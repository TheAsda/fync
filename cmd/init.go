package cmd

import (
	"errors"
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	r "theasda/fync/pkg/repo"

	"github.com/golobby/container/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func HandleInit(context *cli.Context) error {
	var initialized bool
	if e := container.NamedResolve(&initialized, "initialized"); e != nil {
		panic(e)
	}

	var config c.Config
	var repo *r.Repo
	if initialized {
		if e := container.Resolve(&repo); e != nil {
			panic(e)
		}
	} else {
		config, err := c.PromptConfig()
		if err != nil {
			return err
		}
		err = c.SaveConfig(config)
		if err != nil {
			return err
		}
		repo = r.NewRepo(config)
	}

	if repo.Exists() {
		return errors.New("repository already initialized")
	}

	if err := repo.Clone(); err != nil {
		return err
	}

	var err error
	container.Call(func(filesProcessor files_processor.FilesProcessor) {
		var files []string
		files, err = filesProcessor.Files()
		if err != nil {
			return
		}
		if len(files) != 0 {
			ignoredFiles, filesMapping, err := c.PromptFiles(files)
			if err != nil {
				return
			}
			config.IgnoredFiles = ignoredFiles
			config.FilesMapping = filesMapping
			err = c.SaveConfig(config)
			if err != nil {
				return
			}
		}
	})
	if err != nil {
		return err
	}

	logrus.Info("Initialization completed")
	return nil
}
