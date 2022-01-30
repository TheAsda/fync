package files_processor

import (
	"os"
	"theasda/fync/pkg/config"
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

func (sp *CopyProcessor) Add(file string, path string) error {
	logrus.Debug("Copying file")
	return utils.CopyFile(path, sp.FilesProcessorBase.getFilePath(file))
}

func (sp *CopyProcessor) Remove(file string) error {
	logrus.Debug("Removing file")
	return os.Remove(sp.FilesProcessorBase.getFilePath(file))
}

func (sp *CopyProcessor) Update(files map[string]string) error {
	logrus.Info("Updating files")
	for file, path := range files {
		filePath := sp.FilesProcessorBase.getFilePath(file)
		logrus.Debugf("Checking %s", file)
		areEqual, err := utils.CompareFiles(filePath, path)
		if err != nil {
			return err
		}
		if areEqual {
			logrus.Debug("File did not change")
			continue
		}
		logrus.Debug("File changed")
		if err := utils.CopyFile(path, filePath); err != nil {
			return err
		}
	}
	return nil
}

func (sp *CopyProcessor) Exists(file string) bool {
	files, err := sp.FilesProcessorBase.readDir()
	if err != nil {
		return false
	}
	for _, f := range files {
		if f == file {
			return true
		}
	}
	return false
}

func (sp *CopyProcessor) Files() ([]string, error) {
	return sp.FilesProcessorBase.readDir()
}
