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
	var err error = nil
	if e := container.Call(func(repo *lib.Repo) {
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}
