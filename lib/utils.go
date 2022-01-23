package lib

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func CopyFile(from string, to string) error {
	logrus.Debugf("Copying %s to %s", from, to)
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, data, 0644)
}

func CompareFiles(first string, secord string) (bool, error) {
	firstData, err := ioutil.ReadFile(first)
	if err != nil {
		return false, err
	}
	secondData, err := ioutil.ReadFile(secord)
	if err != nil {
		return false, err
	}
	return bytes.Equal(firstData, secondData), nil
}
