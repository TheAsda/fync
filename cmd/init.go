package cmd

import (
	"path/filepath"
	"theasda/fync/lib"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func HandleInit(context *cli.Context) error {
	c, err := lib.GetConfig()
	if err == nil {
		return lib.InitRepo(*c)
		// return errors.New("Config already exists")
	}
	config, err := PromptConfig()
	if err != nil {
		return err
	}
	if err = lib.SaveConfig(&config); err != nil {
		return err
	}
	return lib.InitRepo(config)
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
