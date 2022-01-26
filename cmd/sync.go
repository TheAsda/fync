package cmd

import (
	"theasda/fync/pkg/files_processor"
	"theasda/fync/pkg/repo"
	"theasda/fync/pkg/storage"
	"theasda/fync/pkg/utils"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleSync(context *cli.Context) error {
	utils.CheckInitialization()
	var err error
	if e := container.Call(func(
		storage *storage.Storage,
		filesProcessor files_processor.FilesProcessor,
		repo *repo.Repo,
	) {
		err = filesProcessor.Update(storage.Files)
		if err != nil {
			return
		}
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}
