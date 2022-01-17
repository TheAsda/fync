package lib

import (
	"errors"
	"io/fs"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "./fync_config.yaml"

type Config struct {
	Files []string
}

func (config *Config) AddFile(file string) error {
	for _, f := range config.Files {
		if f == file {
			return errors.New("File is already added")
		}
	}
	config.Files = append(config.Files, file)

	return nil
}

func (config *Config) RemoveFile(file string) error {
	for i, f := range config.Files {
		if f == file {
			newFiles := append(config.Files[:i], config.Files[i+1:]...)
			config.Files = newFiles
			return nil
		}
	}
	return errors.New("File was not added")
}

func GetConfig() (config *Config, err error) {
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
