package lib

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "./fync_config.yaml"

const (
	SymlinkMode = "symlink"
	CopyMode    = "copy"
)

type Config struct {
	Repository   string
	Path         string
	SyncOnAction bool `yaml:"syncOnAction"`
	Mode         string
}

func GetConfig() (config *Config, err error) {
	if !fileExists(ConfigFile) {
		return nil, errors.New("Config does not exist")
	}
	bytes, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

func SaveConfig(config *Config) error {
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigFile, bytes, fs.ModeCharDevice)
}

func (config *Config) GetFilesPath() string {
	return path.Join(config.Path, "files")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
