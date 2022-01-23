package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

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
	var filesDb *lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return err
	}
	file := lib.File{ID: id, Path: fullPath}
	if err = filesDb.Add(file); err != nil {
		return err
	}
	container.Call(func(filesProcessor lib.FilesProcessor) error {
		return filesProcessor.Add(file)
	})
	return container.Call(func(repo *lib.Repo) error {
		return repo.CommitFiles()
	})
}

func getId(path string, name string) (string, error) {
	var filesDb *lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return "", err
	}
	var id string
	if len(name) != 0 {
		id = name
	} else if name := filepath.Base(path); !filesDb.Exists(name) {
		id = filepath.Base(name)
	} else {
		return "", errors.New("file name already taken, please specify custom name")
	}
	return id, nil
}
