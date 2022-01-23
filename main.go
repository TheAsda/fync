package main

import (
	"errors"
	"os"
	"theasda/fync/cmd"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	app := &cli.App{
		Name:  "fync",
		Usage: "Sync specified files with provided git repository",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Show debug logs",
			},
		},
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
		Before: func(c *cli.Context) error {
			enableDebug := c.Bool("debug")
			if enableDebug {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}
			initializeContainer()
			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func initializeContainer() {
	if err := container.Singleton(func() *lib.Config {
		logrus.Debug("Initializing Config")
		config, err := lib.GetConfig()
		if err != nil {
			return nil
		}
		return &config
	}); err != nil {
		panic(err)
	}

	if err := container.Singleton(func(config *lib.Config) *lib.FilesDB {
		logrus.Debug("Initializing Files DB")
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
		logrus.Debug("Initializing Repo")
		if config == nil {
			return nil
		}
		return lib.NewRepo(*config)
	}); err != nil {
		panic(err)
	}

	if err := container.Singleton(func(config *lib.Config) lib.FilesProcessor {
		logrus.Debug("Initializing File Processor")
		if config == nil {
			return &lib.CopyProcessor{}
		}
		if config.Mode == lib.SymlinkMode {
			logrus.Debug("Using Symlink Processor")
			return lib.NewSymlinkProcessor(*config)
		}
		if config.Mode == lib.CopyMode {
			logrus.Debug("Using Copy Processor")
			return lib.NewCopyProcessor(*config)
		}
		panic(errors.New("unknown mode"))
	}); err != nil {
		panic(err)
	}
}
