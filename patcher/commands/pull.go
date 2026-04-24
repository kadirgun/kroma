package commands

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"patcher/models"
	"patcher/store"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

var Pull = cli.Command{
	Name:  "pull",
	Usage: "Pull patches from a source",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "reset",
			Usage: "Delete current patches and re-pull everything",
		},
	},
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "target",
			UsageText: "Target path to pull patches from",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		if c.Bool("reset") {
			log.Info("Resetting patches...")
			patchesPath := filepath.Join(store.Config.BasePath, "patches")
			err := os.RemoveAll(patchesPath)
			if err != nil {
				return err
			}
		}

		targetPath := strings.TrimPrefix(c.StringArg("target"), store.Config.Repo)
		scanDirs := append(store.Config.Submodules, ".")
		if targetPath != "" {
			log.Infof("Pulling patches from %s", targetPath)
			scanDirs = []string{targetPath}
		}

		patches, err := getPatchesFromPaths(scanDirs)
		log.Debugf("Found %d patches", len(patches))
		if err != nil {
			return err
		}

		for _, patch := range patches {
			if strings.HasSuffix(patch.TargetPath, "rej") {
				return fmt.Errorf("Found reject file %s, review and resolve it before pulling", patch.TargetPath)
			} else if patch.HasRejects() {
				return fmt.Errorf("Found patch %s with rejects %s, review and resolve them before pulling", filepath.Base(patch.SourcePath), patch.RejectsPath)
			}
		}

		if targetPath != "" && len(patches) == 0 {
			log.Warnf("No patches found in %s, make sure the path is correct and has changes", targetPath)
		}

		savedCount, err := saveDiffs(patches)
		if err != nil {
			log.Error("Error saving diffs")
		}

		log.Infof("Saved %d patches", savedCount)

		return nil
	},
}

func getPatchesFromPath(path string) ([]*models.Patch, error) {
	log.Debugf("Getting patches from path: %s", path)

	// Check if path is a file or directory
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var dir string
	if info.IsDir() {
		dir = path
	} else {
		dir = filepath.Dir(path)
	}

	args := []string{"status", "--short", "--untracked-files=all", "."}
	cmd := exec.Command("git", args...)
	log.Debugf("Running command: git %s", strings.Join(args, " "))

	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var patches []*models.Patch
	for scanner.Scan() {
		line := scanner.Text()

		log.Debugf("Git status line: %s", line)

		fileName := strings.TrimSpace(line[3:])

		if shouldIgnore(fileName) {
			log.Debugf("Ignoring file %s due to match pattern", fileName)
			continue
		}

		targetPath := filepath.Join(cmd.Dir, fileName)
		patch, err := models.NewPatchFromTarget(targetPath)
		if err != nil {
			log.Errorf("Error creating patch for %s: %v", fileName, err)
			continue
		}

		patch.GitOperation = strings.TrimSpace(line[:2])

		patches = append(patches, patch)
	}

	return patches, nil
}

func getPatchesFromPaths(paths []string) ([]*models.Patch, error) {
	var patches []*models.Patch

	for _, path := range paths {
		submoduleDir := filepath.Join(store.Config.BasePath, store.Config.Repo, path)

		// If path is a file, get the directory containing it
		info, err := os.Stat(submoduleDir)
		if err != nil {
			return nil, err
		}

		dir := submoduleDir
		if !info.IsDir() {
			dir = filepath.Dir(submoduleDir)
		}

		submodulePatches, err := getPatchesFromPath(dir)
		if err != nil {
			return nil, err
		}

		patches = append(patches, submodulePatches...)
	}

	return patches, nil
}

func saveDiffs(patches []*models.Patch) (int, error) {
	savedCount := 0

	for _, patch := range patches {
		err := patch.SaveDiff()
		if err != nil {
			return savedCount, err
		}

		savedCount++
	}

	return savedCount, nil
}

func shouldIgnore(fileName string) bool {
	for _, pattern := range store.Config.IgnorePatterns {
		matched, err := filepath.Match(pattern, fileName)
		if err == nil && matched {
			return true
		}
	}
	return false
}
