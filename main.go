package main

import (
	"errors"
	"os"
	"theasda/fync/cmd"
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	"theasda/fync/pkg/repo"

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
	config, err := c.GetConfig()

	if e := container.NamedSingleton("initialized", func() bool {
		return err == nil
	}); e != nil {
		panic(e)
	}

	if err != nil {
		return
	}

	if e := container.Singleton(func() c.Config {
		return config
	}); e != nil {
		panic(e)
	}

	if e := container.Singleton(func(config c.Config) *repo.Repo {
		return repo.NewRepo(config)
	}); e != nil {
		panic(e)
	}

	if err := container.Singleton(func(config c.Config) files_processor.FilesProcessor {
		logrus.Debug("Initializing file processor")
		if config.Mode == c.SymlinkMode {
			logrus.Debug("Using Symlink Processor")
			return files_processor.NewSymlinkProcessor(config)
		}
		if config.Mode == c.CopyMode {
			logrus.Debug("Using Copy Processor")
			return files_processor.NewCopyProcessor(config)
		}
		panic(errors.New("unknown mode"))
	}); err != nil {
		panic(err)
	}
}
