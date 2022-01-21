package main

import (
	"errors"
	"log"
	"os"
	"theasda/fync/cmd"
	"theasda/fync/lib"
	"time"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "fync",
		Usage:    "Sync specified files with provided git repository",
		Compiled: time.Now(),
		Authors:  []*cli.Author{{Name: "Andrey Kiselev", Email: "omega-faworit@yandex.ru"}},
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
				Flags:     []cli.Flag{&cli.StringFlag{Name: "name", Usage: "Name which will be used as ID of file"}},
				Action:    cmd.HandleAdd,
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

	container.Singleton(func() *lib.Config {
		config, err := lib.GetConfig()
		if err != nil {
			return nil
		}
		return &config
	})
	container.Singleton(func(config lib.Config) *lib.FilesDB {
		files, err := lib.NewFilesDb(config.GetFilesPath())
		if err != nil {
			panic(err)
		}
		return files
	})
	container.Singleton(func(config lib.Config) *lib.Repo {
		return lib.NewRepo(config)
	})
	container.Singleton(func(config lib.Config) lib.FilesProcessor {
		if config.Mode == lib.SymlinkMode {
			return lib.NewSymlinkProcessor(config)
		}
		if config.Mode == lib.CopyMode {
			return lib.NewCopyProcessor(config)
		}
		panic(errors.New("unknown mode"))
	})

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ok")
	}
}
