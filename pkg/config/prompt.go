package config

import (
	"fmt"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

func PromptConfig() (Config, error) {
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
		Items:        []string{SymlinkMode, CopyMode},
		HideSelected: true,
	}

	repository, err := repositoryPrompt.Run()
	if err != nil {
		return Config{}, err
	}
	path, err := pathPrompt.Run()
	if err != nil {
		return Config{}, err
	}
	_, syncOnAction, err := syncOnActionPrompt.Run()
	if err != nil {
		return Config{}, err
	}
	_, mode, err := modePrompt.Run()
	if err != nil {
		return Config{}, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Repository:   repository,
		Path:         absPath,
		SyncOnAction: syncOnAction == "yes",
		Mode:         mode,
		IgnoredFiles: []string{},
	}, nil
}

func PromptFiles(files []string) (ignoredFiles []string, filesMapping map[string]string, err error) {
	filesMapping = make(map[string]string)
	for _, file := range files {
		ignoreFilePrompt := promptui.SelectWithAdd{
			Label:    fmt.Sprintf("Would you like to ignore %s", file),
			Items:    []string{"yes", "no"},
			AddLabel: "Path",
		}
		_, answer, err := ignoreFilePrompt.Run()
		if err != nil {
			return []string{}, nil, err
		}
		if answer == "yes" {
			ignoredFiles = append(ignoredFiles, file)
			continue
		}
		absPath, err := filepath.Abs(answer)
		if err != nil {
			return []string{}, nil, err
		}
		filesMapping[file] = absPath
	}
	return ignoredFiles, filesMapping, nil
}
