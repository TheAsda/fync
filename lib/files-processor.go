package lib

import (
	"os"
	"path"
)

type FilesProcessor interface {
	Add(file File) error
	Remove(file File) error
	Update(files []File) error
}

type FilesProcessorBase struct {
	config Config
}

func (fpb FilesProcessorBase) getIdPath(id string) string {
	return path.Join(fpb.config.Path, id)
}

type SymlinkProcessor struct {
	base FilesProcessorBase
}

type CopyProcessor struct {
	base FilesProcessorBase
}

func NewSymlinkProcessor(config Config) *SymlinkProcessor {
	return &SymlinkProcessor{
		base: FilesProcessorBase{
			config: config,
		},
	}
}

func NewCopyProcessor(config Config) *CopyProcessor {
	return &CopyProcessor{
		base: FilesProcessorBase{
			config: config,
		},
	}
}

func (sp *SymlinkProcessor) Add(file File) error {
	return os.Symlink(file.Path, sp.base.getIdPath(file.ID))
}

func (sp *CopyProcessor) Add(file File) error {
	return CopyFile(file.Path, sp.base.getIdPath(file.ID))
}

func (sp *SymlinkProcessor) Remove(file File) error {
	return os.Remove(sp.base.getIdPath(file.ID))
}

func (sp *CopyProcessor) Remove(file File) error {
	return os.Remove(sp.base.getIdPath(file.ID))
}

func (sp *SymlinkProcessor) Update(files []File) error {
	return nil
}

func (sp *CopyProcessor) Update(files []File) error {
	for _, file := range files {
		idPath := sp.base.getIdPath(file.ID)
		areEqual, err := CompareFiles(file.Path, idPath)
		if err != nil {
			return err
		}
		if areEqual {
			continue
		}
		if err := CopyFile(file.Path, idPath); err != nil {
			return err
		}
	}
	return nil
}
