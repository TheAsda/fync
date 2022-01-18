package cmd

import (
	// "theasda/fync/lib"

	"theasda/fync/lib"

	"github.com/urfave/cli/v2"
)

func HandleSync(context *cli.Context) error {
	config, err := lib.GetConfig()
	if err != nil {
		return err
	}
	files, err := lib.GetFiles(config.GetFilesPath())
	if err != nil {
		return err
	}
	list, err := files.GetFiles()
	if err != nil {
		return err
	}
	return lib.SyncFiles(list, *config)
}
