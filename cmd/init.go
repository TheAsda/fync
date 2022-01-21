package cmd

import (
	"errors"
	"path/filepath"
	"theasda/fync/lib"

	"github.com/golobby/container/v3"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func HandleInit(context *cli.Context) error {
	var config *lib.Config
	var repo *lib.Repo
	err := container.Resolve(&config)
	if err != nil {
		return err
	}
	if config == nil {
		*config, err = PromptConfig()
		if err != nil {
			return err
		}
		if err = lib.SaveConfig(*config); err != nil {
			return err
		}
		repo = lib.NewRepo(*config)
	} else {
		err = container.Resolve(&repo)
		if err != nil {
			return err
		}
	}
	if repo.Exists() {
		return errors.New("repository already initialized")
	}
	return repo.Clone()
}

func PromptConfig() (lib.Config, error) {
	repositoryPrompt := promptui.Prompt{
		Label: "Git repository",
	}
	pathPrompt := promptui.Prompt{
		Label: "Path where local repository will be placed",
	}
	syncOnActionPrompt := promptui.Select{
		Label: "Sync on action",
		Items: []string{"yes", "no"},
	}
	modePrompt := promptui.Select{
		Label: "Mode",
		Items: []string{lib.SymlinkMode, lib.CopyMode},
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
