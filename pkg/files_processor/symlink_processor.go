package files_processor

import (
	"os"
	"theasda/fync/pkg/config"

	"github.com/sirupsen/logrus"
)

type SymlinkProcessor struct {
	FilesProcessorBase
}

func NewSymlinkProcessor(config config.Config) *SymlinkProcessor {
	return &SymlinkProcessor{
		FilesProcessorBase{
			config: config,
		},
	}
}

func (sp *SymlinkProcessor) Add(file string, path string) error {
	logrus.Debug("Creating symlink")
	return os.Symlink(path, sp.FilesProcessorBase.getFilePath(file))
}

func (sp *SymlinkProcessor) Remove(file string) error {
	logrus.Debug("Removing symlink")
	return os.Remove(sp.FilesProcessorBase.getFilePath(file))
}

func (sp *SymlinkProcessor) Update(files map[string]string) error {
	return nil
}

func (sp *SymlinkProcessor) Exists(file string) bool {
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

func (sp *SymlinkProcessor) Files() ([]string, error) {
	return sp.FilesProcessorBase.readDir()
}
