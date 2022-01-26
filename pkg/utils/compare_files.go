package utils

import (
	"bytes"
	"io/ioutil"
)

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
