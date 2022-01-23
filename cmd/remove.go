package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleRemove(context *cli.Context) error {
	fileOrId := context.Args().Get(0)
	if len(fileOrId) == 0 {
		return errors.New("file is not provided")
	}
	var filesDb *lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return err
	}
	filePath, err := filepath.Abs(fileOrId)
	if err != nil {
		return err
	}
	file, err := filesDb.RemoveByPath(filePath)
	if err != nil {
		return err
	}
	if err = container.Call(func(filesProcessor lib.FilesProcessor) error {
		return filesProcessor.Remove(file)
	}); err != nil {
		return err
	}
	return container.Call(func(repo *lib.Repo) error {
		return repo.CommitFiles()
	})
}
