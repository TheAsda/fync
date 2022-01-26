package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"theasda/fync/pkg/config"
	"theasda/fync/pkg/utils"

	"github.com/sirupsen/logrus"
)

type Storage struct {
	config config.Config
	Files  []File `json:"files"`
}

func NewStorage(config config.Config) (*Storage, error) {
	logrus.Debug("Initializing storage")
	if !utils.FileExists(config.GetStoragePath()) {
		logrus.Debug("Storage path does not exist using empty storage")
		return &Storage{
			Files: []File{},
		}, nil
	}
	b, err := ioutil.ReadFile(config.GetStoragePath())
	if err != nil {
		return nil, err
	}
	var storage Storage
	err = json.Unmarshal(b, &storage)
	return &storage, err
}

func (s *Storage) save() error {
	logrus.Debug("Saving storage")
	b, err := json.Marshal(s.Files)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.config.GetStoragePath(), b, 0644)
}

func (s *Storage) remove(index int) {
	s.Files = append(s.Files[:index], s.Files[index+1:]...)
}

func (s *Storage) Exists(id string) bool {
	for _, file := range s.Files {
		if file.ID == id {
			return true
		}
	}
	return false
}

func (s *Storage) Add(file File) error {
	logrus.Debug("Adding file")
	if s.Exists(file.ID) {
		return errors.New("ID already exists")
	}
	s.Files = append(s.Files, file)
	return s.save()
}

func (s *Storage) Remove(id string) error {
	logrus.Debug("Removing file")
	for i, file := range s.Files {
		if file.ID == id {
			s.remove(i)
			return s.save()
		}
	}
	return errors.New("cannot find file")
}

func (s *Storage) RemoveByPath(path string) (File, error) {
	logrus.Debug("Removing file by path")
	for i, file := range s.Files {
		if file.Path == path {
			s.remove(i)
			return file, s.save()
		}
	}
	return File{}, errors.New("cannot find file")
}
