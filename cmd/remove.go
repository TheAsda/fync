package cmd

import (
	"errors"
	"path/filepath"
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	r "theasda/fync/pkg/repo"
	"theasda/fync/pkg/utils"

	"github.com/golobby/container/v3"
	"github.com/urfave/cli/v2"
)

func HandleRemove(context *cli.Context) error {
	utils.CheckInitialization()
	file := context.Args().Get(0)
	if len(file) == 0 {
		return errors.New("file is not provided")
	}
	filePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	if e := container.Call(func(
		config c.Config,
		filesProcessor files_processor.FilesProcessor,
		repo *r.Repo,
	) {
		file, err2 := config.FindFile(filePath)
		if err2 != nil {
			err = err2
			return
		}
		delete(config.FilesMapping, file)
		err = filesProcessor.Remove(file)
		if err != nil {
			return
		}
		c.SaveConfig(config)
		err = repo.UpdateRepo()
	}); e != nil {
		panic(e)
	}
	return err
}
