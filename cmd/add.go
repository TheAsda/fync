package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/pkg/files_processor"
	"theasda/fync/pkg/repo"
	"theasda/fync/pkg/storage"
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
		storage *storage.Storage,
		filesProcessor files_processor.FilesProcessor,
		repo *repo.Repo,
	) {
		err = storage.Add(file)
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
	var storage *storage.Storage
	if err := container.Resolve(&storage); err != nil {
		panic(err)
	}
	var file string
	if len(name) != 0 {
		file = name
	} else if name := filepath.Base(path); !storage.Exists(name) {
		file = filepath.Base(name)
	} else {
		return "", errors.New("file name already taken, please specify custom name")
	}
	return file, nil
}
