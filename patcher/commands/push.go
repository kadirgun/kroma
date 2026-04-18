package commands

import (
	"context"
	"fmt"
	"io/fs"
	"patcher/models"
	"patcher/store"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

var Push = cli.Command{
	Name:  "push",
	Usage: "Push patches to a source",
	Flags: []cli.Flag{},
	Action: func(ctx context.Context, c *cli.Command) error {
		return filepath.Walk(store.Config.GetPatchesDir(), func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			log.Debugf("Found patch file: %s", path)

			patch, err := models.NewPatch(path)

			log.Debugf("Target path: %s", patch.TargetPath)

			if patch.HasRejects() {
				return fmt.Errorf("Patch %s has rejects %s, review and resolve them before pushing", filepath.Base(path), patch.RejectsPath)
			}

			if patch.IsApplied() {
				log.Debugf("Patch %s already applied, skipping...", path)
				return nil
			}

			if patch.IsDirty() {
				return fmt.Errorf("%s is dirty, review and pull the latest changes before pushing", patch.TargetPath)
			}

			err = patch.Apply()

			if err != nil {
				return err
			}

			if patch.HasRejects() {
				return fmt.Errorf("Patch %s applied with rejects, review and resolve them before pushing", path)
			}

			return nil
		})
	},
}
