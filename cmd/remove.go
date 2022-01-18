package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/urfave/cli/v2"
)

func HandleRemove(context *cli.Context) error {
	fileOrId := context.Args().Get(0)
	if len(fileOrId) == 0 {
		return errors.New("File is not provided")
	}
	config, err := lib.GetConfig()
	if err != nil {
		return err
	}
	files, err := lib.GetFiles(config.GetFilesPath())
	if err != nil {
		return err
	}
	if files.Exists(fileOrId) {
		return files.RemoveFile(fileOrId)
	}
	filePath, err := filepath.Abs(fileOrId)
	if err != nil {
		return err
	}
	err = files.RemoveByPath(filePath)
	if err != nil {
		return err
	}
	if !config.SyncOnAction {
		return nil
	}
	println("Syncing")
	return nil
}
