package cmd

import (
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	r "theasda/fync/pkg/repo"
	"theasda/fync/pkg/utils"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleSync(context *cli.Context) error {
	utils.CheckInitialization()
	var err error
	if e := container.Call(func(
		config c.Config,
		filesProcessor files_processor.FilesProcessor,
		repo *r.Repo,
	) {
		err = filesProcessor.Update(config.FilesMapping)
		if err != nil {
			return
		}
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}
