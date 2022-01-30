package cmd

import (
	"path"
	"strings"
	c "theasda/fync/pkg/config"
	"theasda/fync/pkg/files_processor"
	r "theasda/fync/pkg/repo"
	"theasda/fync/pkg/utils"

	"github.com/golobby/container/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func HandleCheck(context *cli.Context) error {
	utils.CheckInitialization()
	var err error
	if e := container.Call(func(
		config c.Config,
		filesProcessor files_processor.FilesProcessor,
		repo *r.Repo,
	) {
		var files []string
		files, err = filesProcessor.Files()
		if err != nil {
			return
		}
		var changedFiles []string
		for _, file := range files {
			filePath := config.FilesMapping[file]

			areEqual := utils.CompareFiles(path.Join(config.Path, file), filePath)
			if !areEqual {
				changedFiles = append(changedFiles, file)
			}
		}
		if len(changedFiles) == 0 {
			logrus.Info("No changes")
			return
		}
		logrus.Info("Changed files: ", strings.Join(changedFiles, ", "))
	}); e != nil {
		panic(e)
	}
	return err
}
