package cmd

import (
	"crypto/sha1"
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/urfave/cli/v2"
)

func HandleAdd(context *cli.Context) error {
	file := context.Args().Get(0)
	if len(file) == 0 {
		return errors.New("File is not provided")
	}
	filePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	config, err := lib.GetConfig()
	if err != nil {
		return err
	}
	files, err := lib.GetFiles(config.GetFilesPath())
	if err != nil {
		return err
	}
	var id string
	if name := context.String("name"); len(name) != 0 {
		id = name
	} else if name := filepath.Base(filePath); !files.Exists(name) {
		id = filepath.Base(name)
	} else {
		hash := sha1.Sum([]byte(filePath))
		id = string(hash[:])
	}
	err = files.AddFile(id, filePath)
	if err != nil {
		return err
	}
	return lib.AddFile(id, filePath, *config)
}
