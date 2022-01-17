package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/urfave/cli/v2"
)

func HandleRemove(context *cli.Context) error {
	file := context.Args().Get(0)
	if len(file) == 0 {
		return errors.New("File is not provided")
	}
	config, err := lib.GetConfig()
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	err = config.RemoveFile(absPath)
	if err != nil {
		return err
	}
	return lib.SaveConfig(config)
}
