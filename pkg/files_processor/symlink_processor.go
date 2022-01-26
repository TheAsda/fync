package files_processor

import (
	"os"
	"theasda/fync/pkg/config"
	"theasda/fync/pkg/storage"

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

func (sp *SymlinkProcessor) Add(file storage.File) error {
	logrus.Debug("Creating symlink")
	return os.Symlink(file.Path, sp.FilesProcessorBase.getIdPath(file.ID))
}

func (sp *SymlinkProcessor) Remove(file storage.File) error {
	logrus.Debug("Removing symlink")
	return os.Remove(sp.FilesProcessorBase.getIdPath(file.ID))
}

func (sp *SymlinkProcessor) Update(files []storage.File) error {
	return nil
}
