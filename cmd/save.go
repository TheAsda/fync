package cmd

import (
	"theasda/fync/lib"

	"github.com/urfave/cli/v2"
)

func HandleSave(context *cli.Context) error {
	config, err := lib.GetConfig()
	if err != nil {
		return err
	}

	return nil
}
