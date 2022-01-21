package cmd

import (
	// "theasda/fync/lib"

	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleSync(context *cli.Context) error {
	var filesDb lib.FilesDB
	if err := container.Resolve(&filesDb); err != nil {
		return err
	}
	var filesProcessor lib.FilesProcessor
	if err := container.Resolve(&filesProcessor); err != nil {
		return err
	}
	if err := filesProcessor.Update(filesDb.GetAll()); err != nil {
		return err
	}
	var repo *lib.Repo
	if err := container.Resolve(&repo); err != nil {
		return err
	}
	return repo.CommitFiles()
}
