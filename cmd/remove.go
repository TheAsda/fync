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

func HandleRemove(context *cli.Context) error {
	utils.CheckInitialization()
	fileOrId := context.Args().Get(0)
	if len(fileOrId) == 0 {
		return errors.New("file is not provided")
	}
	filePath, err := filepath.Abs(fileOrId)
	if err != nil {
		return err
	}
	if e := container.Call(func(
		storage *storage.Storage,
		filesProcessor files_processor.FilesProcessor,
		repo *repo.Repo,
	) {
		file, err := storage.RemoveByPath(filePath)
		if err != nil {
			return
		}
		err = filesProcessor.Remove(file)
		if err != nil {
			return
		}
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}
