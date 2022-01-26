package utils

import (
	"errors"

	"github.com/golobby/container/v3"
)

func CheckInitialization() {
	var initialized bool
	if e := container.NamedResolve(&initialized, "initialized"); e != nil {
		panic(e)
	}
	if !initialized {
		panic(errors.New("Not initialized"))
	}
}
