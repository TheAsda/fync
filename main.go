package main

import (
	"log"
	"os"
	"theasda/fync/cmd"
	"time"

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
				Name:      "add",
				Usage:     "Add file for syncing",
				ArgsUsage: "[file]",
				Action:    cmd.HandleAdd,
			},
			{
				Name:  "sync",
				Usage: "Sync files",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:   "remove",
				Usage:  "Remove file from syncing",
				Action: cmd.HandleRemove,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ok")
	}
}
