package files_processor

import (
	"theasda/fync/pkg/storage"
)

type FilesProcessor interface {
	Add(file storage.File) error
	Remove(file storage.File) error
	Update(files []storage.File) error
}
