package main

import (
	"errors"
	"fmt"
	"os"
	"theasda/fync/cmd"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := container.Singleton(func() *lib.Config {
		config, err := lib.GetConfig()
		if err != nil {
			return nil
		}
		return &config
	}); err != nil {
		panic(err)
	}

	if err := container.Singleton(func(config *lib.Config) *lib.FilesDB {
		if config == nil {
			return nil
		}
		files, err := lib.NewFilesDb(config.GetFilesPath())
		if err != nil {
			panic(err)
		}
		return files
	}); err != nil {
		panic(err)
	}

	if err := container.Singleton(func(config *lib.Config) *lib.Repo {
		if config == nil {
			return nil
		}
		return lib.NewRepo(*config)
	}); err != nil {
		panic(err)
	}

	if err := container.Singleton(func(config *lib.Config) lib.FilesProcessor {
		if config == nil {
			return &lib.CopyProcessor{}
		}
		if config.Mode == lib.SymlinkMode {
			return lib.NewSymlinkProcessor(*config)
		}
		if config.Mode == lib.CopyMode {
			return lib.NewCopyProcessor(*config)
		}
		panic(errors.New("unknown mode"))
	}); err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "fync",
		Usage: "Sync specified files with provided git repository",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize config",
				Action: cmd.HandleInit,
			},
			{
				Name:      "add",
				Usage:     "Add file for syncing",
				ArgsUsage: "[file]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "Name which will be used as ID of file",
					},
				},
				Action: cmd.HandleAdd,
			},
			{
				Name:   "sync",
				Usage:  "Sync files",
				Action: cmd.HandleSync,
			},
			{
				Name:      "remove",
				Usage:     "Remove file from syncing",
				ArgsUsage: "[id or file]",
				Action:    cmd.HandleRemove,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("Success")
	}
}
