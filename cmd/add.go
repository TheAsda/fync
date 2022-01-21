package cmd

import (
	"crypto/sha1"
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleAdd(context *cli.Context) error {
	file := context.Args().Get(0)
	if len(file) == 0 {
		return errors.New("file is not provided")
	}
	fullPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	id, err := getId(fullPath, context.String("name"))
	if err != nil {
		return err
	}
	var filesDb lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return err
	}
	if err = filesDb.Add(lib.File{ID: id, Path: fullPath}); err != nil {
		return err
	}
	return container.Call(func(repo lib.Repo) error {
		return repo.CommitFiles()
	})
}

func getId(path string, name string) (string, error) {
	var filesDb lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return "", err
	}
	var id string
	if len(name) != 0 {
		id = name
	} else if name := filepath.Base(path); !filesDb.Exists(name) {
		id = filepath.Base(name)
	} else {
		hash := sha1.Sum([]byte(path))
		id = string(hash[:])
	}
	return id, nil
}
