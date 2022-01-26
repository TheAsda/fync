package cmd

import (
	"errors"
	с "theasda/fync/pkg/config"
	r "theasda/fync/pkg/repo"

	"github.com/golobby/container/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func HandleInit(context *cli.Context) error {
	var initialized bool
	if e := container.NamedResolve(&initialized, "initialized"); e != nil {
		panic(e)
	}

	var repo *r.Repo
	if initialized {
		if e := container.Resolve(&repo); e != nil {
			panic(e)
		}
	} else {
		config, err := с.PromptConfig()
		if err != nil {
			return err
		}
		err = с.SaveConfig(config)
		if err != nil {
			return err
		}
		repo = r.NewRepo(config)
	}
	if repo.Exists() {
		return errors.New("repository already initialized")
	}
	if err := repo.Clone(); err != nil {
		return err
	}
	logrus.Info("Initialization completed")
	return nil
}
