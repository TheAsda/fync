package cmd

import (
	"errors"
	"fmt"
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
		var err error
		config, err = c.PromptConfig()
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
	var filesProcessor files_processor.FilesProcessor
	switch config.Mode {
	case c.SymlinkMode:
		logrus.Debug("Using Symlink Processor")
		filesProcessor = files_processor.NewSymlinkProcessor(config)
	case c.CopyMode:
		logrus.Debug("Using Copy Processor")
		filesProcessor = files_processor.NewCopyProcessor(config)
	default:
		panic(errors.New(fmt.Sprintf("Unknown mode \"%s\"", config.Mode)))
	}

	var files []string
	files, err = filesProcessor.Files()
	if err != nil {
		return err
	}

	if len(files) != 0 {
		logrus.Debugf("%d files in repository")
		var ignoredFiles []string
		var filesMapping map[string]string
		ignoredFiles, filesMapping, err = c.PromptFiles(files)
		if err != nil {
			return err
		}
		config.IgnoredFiles = ignoredFiles
		config.FilesMapping = filesMapping
		err = c.SaveConfig(config)
		if err != nil {
			return err
		}
		err = filesProcessor.Load(config.FilesMapping)
	} else {
		logrus.Debug("No files in repository")
	}

	if err != nil {
		return err
	}

	logrus.Info("Initialization completed")
	return nil
}
