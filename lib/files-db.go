package lib

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type File struct {
	Path string
	ID   string
}

type FilesDB struct {
	files []File
	path  string
}

const FilesCollectionName = "files"

func NewFilesDb(path string) (*FilesDB, error) {
	if !FileExists(path) {
		return &FilesDB{
			files: []File{},
			path:  path,
		}, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var files []File
	err = json.Unmarshal(b, &files)
	if err != nil {
		return nil, err
	}
	return &FilesDB{
		files: files,
		path:  path,
	}, nil
}

func (f *FilesDB) save() error {
	b, err := json.Marshal(f.files)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.path, b, 0644)
}

func (f *FilesDB) remove(index int) {
	f.files = append(f.files[:index], f.files[index+1:]...)
}

func (f *FilesDB) Exists(id string) bool {
	for _, file := range f.files {
		if file.ID == id {
			return true
		}
	}
	return false
}

func (f *FilesDB) Add(file File) error {
	if f.Exists(file.ID) {
		return errors.New("ID already exists")
	}
	f.files = append(f.files, file)
	return f.save()
}

func (f *FilesDB) Remove(id string) error {
	for i, file := range f.files {
		if file.ID == id {
			f.remove(i)
			return f.save()
		}
	}
	return errors.New("cannot find file")
}

func (f *FilesDB) RemoveByPath(path string) (File, error) {
	for i, file := range f.files {
		if file.Path == path {
			f.remove(i)
			return file, f.save()
		}
	}
	return File{}, errors.New("cannot find file")
}

func (f *FilesDB) Update(newFile File) error {
	for i, file := range f.files {
		if file.ID == newFile.ID {
			f.files[i] = file
			return f.save()
		}
	}
	return errors.New("cannot find file")
}

func (f *FilesDB) GetAll() (files []File) {
	return f.files
}
