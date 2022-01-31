package cmd

import (
	"errors"
	"path/filepath"
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	"theasda/fync/pkg/repo"
	"theasda/fync/pkg/utils"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleAdd(context *cli.Context) error {
	path := context.Args().Get(0)
	if len(path) == 0 {
		return errors.New("file is not provided")
	}
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file, err := getFile(fullPath, context.String("name"))
	if err != nil {
		return err
	}
	utils.CheckInitialization()
	if e := container.Call(func(
		config c.Config,
		filesProcessor files_processor.FilesProcessor,
		repo *repo.Repo,
	) {
		config.FilesMapping[file] = fullPath
		err = c.SaveConfig(config)
		if err != nil {
			return
		}
		err = filesProcessor.Add(file, fullPath)
		if err != nil {
			return
		}
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}

func getFile(path string, name string) (string, error) {
	var filesProcessor files_processor.FilesProcessor
	if err := container.Resolve(&filesProcessor); err != nil {
		panic(err)
	}
	var file string
	if len(name) != 0 {
		file = name
	} else if name := filepath.Base(path); !filesProcessor.Exists(name) {
		file = filepath.Base(name)
	} else {
		return "", errors.New("file name already taken, please specify custom name")
	}
	return file, nil
}
