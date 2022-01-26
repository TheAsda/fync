package files_processor

import (
	"path"
	"theasda/fync/pkg/config"
)

type FilesProcessorBase struct {
	config config.Config
}

func (processor FilesProcessorBase) getIdPath(id string) string {
	return path.Join(processor.config.Path, id)
}
