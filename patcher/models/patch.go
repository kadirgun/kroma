package models

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"patcher/store"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type Patch struct {
	SourcePath   string
	TargetPath   string
	RejectsPath  string
	GitOperation string
	Submodule    string
}

func NewPatch(sourcePatg string) (*Patch, error) {
	patch := &Patch{
		SourcePath: sourcePatg,
	}

	var err error
	patch.TargetPath, err = filepath.Rel(store.Config.GetPatchesDir(), patch.SourcePath)
	if err != nil {
		return nil, err
	}

	patch.TargetPath = filepath.Join(store.Config.BasePath, store.Config.Repo, patch.TargetPath)
	patch.TargetPath = strings.TrimSuffix(patch.TargetPath, ".diff")

	patch.RejectsPath = fmt.Sprintf("%s.rej", patch.TargetPath)

	patch.Submodule, err = patch.getSubmodule()
	if err != nil {
		return nil, err
	}

	return patch, nil
}

func NewPatchFromTarget(targetPath string) (*Patch, error) {
	relPath, err := filepath.Rel(filepath.Join(store.Config.BasePath, store.Config.Repo), targetPath)
	if err != nil {
		return nil, err
	}

	sourcePath := filepath.Join(store.Config.GetPatchesDir(), relPath) + ".diff"

	return NewPatch(sourcePath)
}

func (p *Patch) Apply() error {
	cmd := exec.Command("git", "apply", "--whitespace=fix", "--reject", p.SourcePath)
	cmd.Dir = filepath.Dir(p.TargetPath)
	output, err := cmd.CombinedOutput()
	outputString := string(output)
	if err != nil {
		if strings.Contains(outputString, "Rejected hunk") {
			return fmt.Errorf("Patch %s applied with rejects %s, review and resolve them before pushing", filepath.Base(p.SourcePath), p.RejectsPath)
		}

		return fmt.Errorf("Error applying patch %s: %v, output: %s", p.SourcePath, err, string(output))
	}

	log.Debugf("Successfully applied patch %s", p.SourcePath)

	return nil
}

func (p *Patch) Pull() error {
	return nil
}

func (p *Patch) HasRejects() bool {
	_, err := os.Stat(p.RejectsPath)
	return err == nil
}

func (p *Patch) IsDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain", p.TargetPath)
	cmd.Dir = filepath.Dir(p.TargetPath)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("Error checking if patch is dirty for %s: %v", p.TargetPath, err)
		return false
	}

	return len(output) > 0
}

func (p *Patch) IsApplied() bool {
	cmd := exec.Command("git", "apply", "--check", "--reverse", "--ignore-whitespace", p.SourcePath)
	cmd.Dir = filepath.Dir(p.TargetPath)
	err := cmd.Run()
	if err == nil {
		return true
	}

	return false
}

func (p *Patch) getSubmodule() (string, error) {
	for _, submodule := range append(store.Config.Submodules, ".") {
		if strings.HasPrefix(filepath.Dir(p.TargetPath), filepath.Join(store.Config.BasePath, store.Config.Repo, submodule)) {
			return submodule, nil
		}
	}

	return "", fmt.Errorf("no matching submodule found for patch %s", p.TargetPath)
}

func (p *Patch) SaveDiff() error {
	file := filepath.Base(p.TargetPath)
	dir := filepath.Dir(p.TargetPath)

	args := []string{"diff"}

	if p.GitOperation == "??" {
		args = append(args, "--no-index", "/dev/null")
	} else {
		args = append(args, "--full-index")
	}

	args = append(args, file)

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	output, err := cmd.Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 1 {
				log.Errorf("Error generating diff for %s: %v", p.TargetPath, err)
				return err
			}
		} else {
			log.Errorf("Error generating diff for %s: %v", p.TargetPath, err)
			return err
		}
	}

	diff := string(output)

	// if has -- /dev/null line, it means it's a new file, we need to fix the diff to be applicable
	if strings.Contains(diff, "-- /dev/null") {
		fileName := filepath.Base(p.TargetPath)
		repoRoot := filepath.Join(store.Config.BasePath, store.Config.Repo)
		relativePath, err := filepath.Rel(repoRoot, p.TargetPath)
		if err != nil {
			log.Fatal(err)
		}

		relativePath = filepath.ToSlash(relativePath)

		diff = strings.ReplaceAll(diff, fmt.Sprintf("+++ b/%s", fileName), fmt.Sprintf("+++ b/%s", relativePath))
		diff = strings.ReplaceAll(
			diff,
			fmt.Sprintf("diff --git a/%s b/%s", fileName, fileName),
			fmt.Sprintf("diff --git a/%s b/%s", relativePath, relativePath),
		)
	}

	diffBytes := []byte(diff)

	newHash := fmt.Sprintf("%x", sha1.Sum(diffBytes))

	if existingContent, err := os.ReadFile(p.SourcePath); err == nil {
		existingHash := fmt.Sprintf("%x", sha1.Sum(existingContent))

		if newHash == existingHash {
			log.Debugf("Diff for %s unchanged, skipping", filepath.Base(p.TargetPath))
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(p.SourcePath), 0755); err != nil {
		log.Errorf("Error creating directories for %s: %v", p.SourcePath, err)
		return err
	}

	err = os.WriteFile(p.SourcePath, diffBytes, 0755)
	if err != nil {
		log.Errorf("Error writing diff for %s: %v", p, err)
		return err
	}

	return nil
}
