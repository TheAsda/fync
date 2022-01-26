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
	id, err := getId(fullPath, context.String("name"))
	if err != nil {
		return err
	}
	utils.CheckInitialization()
	file := storage.File{ID: id, Path: fullPath}
	if e := container.Call(func(
		storage *storage.Storage,
		filesProcessor files_processor.FilesProcessor,
		repo *repo.Repo,
	) {
		err = storage.Add(file)
		if err != nil {
			return
		}
		err = filesProcessor.Add(file)
		if err != nil {
			return
		}
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}

func getId(path string, name string) (string, error) {
	var storage *storage.Storage
	if err := container.Resolve(&storage); err != nil {
		panic(err)
	}
	var id string
	if len(name) != 0 {
		id = name
	} else if name := filepath.Base(path); !storage.Exists(name) {
		id = filepath.Base(name)
	} else {
		return "", errors.New("file name already taken, please specify custom name")
	}
	return id, nil
}
