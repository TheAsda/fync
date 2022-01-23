package lib

import (
	"errors"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "./fync_config.yaml"

const (
	SymlinkMode = "symlink"
	CopyMode    = "copy"
)

type Config struct {
	Repository   string `yaml:"repository"`
	Path         string `yaml:"path"`
	SyncOnAction bool   `yaml:"syncOnAction"`
	Mode         string `yaml:"mode"`
}

func (config Config) GetFilesPath() string {
	return path.Join(config.Path, "files")
}

func GetConfig() (config Config, err error) {
	if !FileExists(ConfigFile) {
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
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigFile, bytes, 0644)
}
