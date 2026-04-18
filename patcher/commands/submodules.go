package commands

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"patcher/store"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gammazero/workerpool"
	"github.com/urfave/cli/v3"
)

var Submodules = cli.Command{
	Name:  "submodules",
	Usage: "List dirty submodules",
	Flags: []cli.Flag{},
	Action: func(ctx context.Context, c *cli.Command) error {
		cmd := exec.Command("git", "submodule", "status", "--recursive")
		cmd.Dir = filepath.Join(store.Config.BasePath, store.Config.Repo)

		output, err := cmd.Output()
		if err != nil {
			return err
		}

		submodules := []string{}
		scanner := bufio.NewScanner(bytes.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()
			submodule := strings.Fields(line)[1]
			submodules = append(submodules, submodule)
		}

		log.Infof("Found %d submodules", len(submodules))

		wp := workerpool.New(6)
		for _, submodule := range submodules {
			sub := submodule
			wp.Submit(func() {
				cmd := exec.Command("git", "diff", "--name-only", ".")
				cmd.Dir = filepath.Join(store.Config.BasePath, store.Config.Repo, sub)
				output, err := cmd.Output()
				if err != nil {
					log.Errorf("Error checking submodule %s: %v", sub, err)
				}

				if len(output) <= 0 {
					return
				}

				if slices.Contains(store.Config.Submodules, sub) {
					return
				}

				log.Warnf("Submodule %s is dirty but not in config", sub)
			})
		}
		wp.StopWait()

		return nil
	},
}
