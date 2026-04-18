package main

import (
	"context"
	"os"
	"patcher/commands"
	"patcher/store"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func main() {
	app := cli.Command{
		Name:  "patcher",
		Usage: "A tool for applying patches to files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "Enable verbose logging",
				Destination: &store.Flags.Verbose,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if store.Flags.Verbose {
				log.SetLevel(log.DebugLevel)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return ctx, err
			}

			store.Config.BasePath, err = findRoot(cwd)
			if err != nil {
				return ctx, err
			}

			configPath := filepath.Join(store.Config.BasePath, "patcher.json")
			if err := store.InitConfig(configPath); err != nil {
				return ctx, err
			}

			return ctx, nil
		},
		Commands: []*cli.Command{
			&commands.Pull,
			&commands.Push,
			&commands.Submodules,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func findRoot(dir string) (string, error) {
	_, err := os.Stat(filepath.Join(dir, "patcher.json"))
	if err == nil {
		return dir, nil
	}

	parent := filepath.Dir(dir)
	if parent == dir {
		return "", os.ErrNotExist
	}

	return findRoot(parent)
}
