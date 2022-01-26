package files_processor

import (
	"os"
	"theasda/fync/pkg/config"
	"theasda/fync/pkg/storage"
	"theasda/fync/pkg/utils"

	"github.com/sirupsen/logrus"
)

type CopyProcessor struct {
	FilesProcessorBase
}

func NewCopyProcessor(config config.Config) *CopyProcessor {
	return &CopyProcessor{
		FilesProcessorBase{
			config: config,
		},
	}
}

func (sp *CopyProcessor) Add(file storage.File) error {
	logrus.Debug("Copying file")
	return utils.CopyFile(file.Path, sp.FilesProcessorBase.getIdPath(file.ID))
}

func (sp *CopyProcessor) Remove(file storage.File) error {
	logrus.Debug("Removing file")
	return os.Remove(sp.FilesProcessorBase.getIdPath(file.ID))
}

func (sp *CopyProcessor) Update(files []storage.File) error {
	logrus.Info("Updating files")
	for _, file := range files {
		idPath := sp.FilesProcessorBase.getIdPath(file.ID)
		logrus.Debugf("Checking %s", file.Path)
		areEqual, err := utils.CompareFiles(file.Path, idPath)
		if err != nil {
			return err
		}
		if areEqual {
			logrus.Debug("File did not change")
			continue
		}
		logrus.Debug("File changed")
		if err := utils.CopyFile(file.Path, idPath); err != nil {
			return err
		}
	}
	return nil
}
