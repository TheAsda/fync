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
	return lib.CommitFiles(*config)
}
