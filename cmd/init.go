package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func HandleInit(context *cli.Context) error {
	var config *lib.Config
	var repo *lib.Repo

	if e := container.Resolve(&config); e != nil {
		panic(e)
	}
	if config == nil {
		newConfig, err := PromptConfig()
		if err != nil {
			return err
		}
		if err = lib.SaveConfig(newConfig); err != nil {
			return err
		}
		repo = lib.NewRepo(newConfig)
	} else {
		if err := container.Resolve(&repo); err != nil {
			panic(err)
		}
	}
	if repo.Exists() {
		return errors.New("repository already initialized")
	}
	if err := repo.Clone(); err != nil {
		return err
	}
	logrus.Info("Init completed")
	return nil
}

func PromptConfig() (lib.Config, error) {
	repositoryPrompt := promptui.Prompt{
		Label:       "Git repository",
		HideEntered: true,
	}
	pathPrompt := promptui.Prompt{
		Label:       "Path where local repository will be placed",
		HideEntered: true,
	}
	syncOnActionPrompt := promptui.Select{
		Label:        "Sync on action",
		Items:        []string{"yes", "no"},
		HideSelected: true,
	}
	modePrompt := promptui.Select{
		Label:        "Mode",
		Items:        []string{lib.SymlinkMode, lib.CopyMode},
		HideSelected: true,
	}

	repository, err := repositoryPrompt.Run()
	if err != nil {
		return lib.Config{}, err
	}
	path, err := pathPrompt.Run()
	if err != nil {
		return lib.Config{}, err
	}
	_, syncOnAction, err := syncOnActionPrompt.Run()
	if err != nil {
		return lib.Config{}, err
	}
	_, mode, err := modePrompt.Run()
	if err != nil {
		return lib.Config{}, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return lib.Config{}, err
	}

	return lib.Config{
		Repository:   repository,
		Path:         absPath,
		SyncOnAction: syncOnAction == "yes",
		Mode:         mode,
	}, nil
}
