package config

import (
	"errors"
	"io/ioutil"
	"path"
	"theasda/fync/pkg/utils"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "./fync_config.yaml"

const (
	SymlinkMode = "symlink"
	CopyMode    = "copy"
)

type Config struct {
	Repository   string            `yaml:"repository"`
	Path         string            `yaml:"path"`
	SyncOnAction bool              `yaml:"syncOnAction"`
	Mode         string            `yaml:"mode"`
	IgnoredFiles []string          `yaml:"ignoredFiles"`
	FilesMapping map[string]string `yaml:"files"`
}

const StorageFileName = "_storage.json"

func (config Config) GetStoragePath() string {
	return path.Join(config.Path, StorageFileName)
}

func (config Config) FindFile(path string) (string, error) {
	for file, p := range config.FilesMapping {
		if p == path {
			return file, nil
		}
	}
	return "", errors.New("cannot find file")
}

func GetConfig() (config Config, err error) {
	logrus.Debug("Getting config")
	if !utils.FileExists(ConfigFile) {
		return Config{}, errors.New("Config does not exist")
	}
	bytes, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

func SaveConfig(config Config) error {
	logrus.Debug("Saving config")
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigFile, bytes, 0777)
}
