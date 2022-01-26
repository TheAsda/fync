package utils

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func CopyFile(from string, to string) error {
	logrus.Debugf("Copying %s to %s", from, to)
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, data, 0777)
}
