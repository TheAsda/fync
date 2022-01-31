package files_processor

import (
	"io/ioutil"
	"path"
	"strings"
	"theasda/fync/pkg/config"
)

type FilesProcessorBase struct {
	config config.Config
}

func (processor FilesProcessorBase) getFilePath(id string) string {
	return path.Join(processor.config.Path, id)
}

func (processor FilesProcessorBase) readDir() ([]string, error) {
	fileInfos, err := ioutil.ReadDir(processor.config.Path)
	if err != nil {
		return []string{}, err
	}
	var files []string
	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), ".git") {
			continue
		}
		files = append(files, fileInfo.Name())
	}
	return files, err
}
