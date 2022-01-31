package utils

import (
	"bytes"
	"io/ioutil"
)

func CompareFiles(first string, secord string) bool {
	firstData, err := ioutil.ReadFile(first)
	if err != nil {
		return false
	}
	secondData, err := ioutil.ReadFile(secord)
	if err != nil {
		return false
	}
	return bytes.Equal(firstData, secondData)
}
