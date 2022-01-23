package cmd

import (
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleSync(context *cli.Context) error {
	if err := container.Call(func(filesDb *lib.FilesDB, filesProcessor lib.FilesProcessor) error {
		return filesProcessor.Update(filesDb.GetAll())
	}); err != nil {
		return err
	}
	return container.Call(func(repo *lib.Repo) error {
		return repo.CommitFiles()
	})
}
