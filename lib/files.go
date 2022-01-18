package lib

import (
	"encoding/json"
	"errors"

	scribble "github.com/nanobox-io/golang-scribble"
)

type File struct {
	Path string
	ID   string
}

type Files struct {
	db *scribble.Driver
}

const FilesCollectionName = "files"

func GetFiles(path string) (*Files, error) {
	db, err := scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return &Files{
		db: db,
	}, nil
}

func (f *Files) AddFile(id string, path string) error {
	err := f.db.Read(FilesCollectionName, id, &File{})
	if err == nil {
		return errors.New("ID already exists")
	}
	return f.db.Write(FilesCollectionName, id, File{
		ID:   id,
		Path: path,
	})
}

func (f *Files) RemoveFile(id string) error {
	return f.db.Delete(FilesCollectionName, id)
}

func (f *Files) RemoveByPath(path string) error {
	strFiles, err := f.db.ReadAll(FilesCollectionName)
	if err != nil {
		return err
	}
	for _, strFile := range strFiles {
		var file File
		if err = json.Unmarshal([]byte(strFile), &file); err != nil {
			return err
		}
		if file.Path == path {
			return f.db.Delete(FilesCollectionName, file.ID)
		}
	}
	return errors.New("Cannot find file with specified path")
}

func (f *Files) UpdateFile(id string, path string) error {
	err := f.RemoveFile(id)
	if err != nil {
		return err
	}
	return f.AddFile(id, path)
}

func (f *Files) Exists(id string) bool {
	var file *File
	return f.db.Read(FilesCollectionName, id, &file) == nil
}

func (f *Files) GetFiles() (files []File, err error) {
	strFiles, err := f.db.ReadAll(FilesCollectionName)
	if err != nil {
		return []File{}, err
	}
	for _, strFile := range strFiles {
		var file File
		if err := json.Unmarshal([]byte(strFile), &file); err != nil {
			return []File{}, err
		}
		files = append(files, file)
	}
	return files, err
}
