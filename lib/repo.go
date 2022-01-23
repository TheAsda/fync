package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Repo struct {
	config Config
}

func NewRepo(config Config) *Repo {
	return &Repo{
		config: config,
	}
}

func (repo Repo) Exists() bool {
	if !FileExists(repo.config.Path) {
		return false
	}
	cmd := exec.Command("git", "status")
	cmd.Dir = repo.config.Path
	err := cmd.Run()
	return err != nil
}

func (repo Repo) Clone() error {
	if err := os.MkdirAll(repo.config.Path, 0644); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", repo.config.Repository, ".")
	cmd.Dir = repo.config.Path
	return cmd.Run()
}

func (repo Repo) StageFiles() error {
	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = repo.config.Path
	return addCmd.Run()
}

func (repo Repo) CommitFiles() error {
	commitMessage, err := repo.getCommitMessage()
	if err != nil {
		return err
	}

	if err = repo.StageFiles(); err != nil {
		return err
	}

	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	commitCmd.Dir = repo.config.Path
	return commitCmd.Run()
}

func (repo Repo) Push() error {
	pushCmd := exec.Command("git", "push")
	pushCmd.Dir = repo.config.Path

	return pushCmd.Run()
}

func (repo *Repo) UpdateRepo() error {
	if err := repo.CommitFiles(); err != nil {
		return err
	}
	return repo.Push()
}

func (repo Repo) getCommitMessage() (string, error) {
	statusCmd := exec.Command("git", "status", "-sb")
	statusCmd.Dir = repo.config.Path
	var stdBuffer bytes.Buffer
	w := io.MultiWriter(&stdBuffer)
	statusCmd.Stdout = w
	err := statusCmd.Run()
	if err != nil {
		return "", err
	}
	status := stdBuffer.String()
	addedFiles, modifiedFiles, deletedFiles := parseStatus(status)
	if len(addedFiles)+len(modifiedFiles)+len(deletedFiles) == 0 {
		return "", errors.New("no changes in files")
	}
	return getCommitMessage(addedFiles, modifiedFiles, deletedFiles), nil
}

func parseStatus(status string) (addedFiles []string, modifiedFiles []string, deletedFiles []string) {
	lines := strings.Split(status, "\n")
	for _, l := range lines {
		if strings.Index(l, "??") == 0 {
			file := strings.ReplaceAll(l, "?? ", "")
			if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
				continue
			}
			addedFiles = append(addedFiles, file)
			continue
		}
		if strings.Index(l, " D") == 0 {
			file := strings.ReplaceAll(l, " D ", "")
			if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
				continue
			}
			deletedFiles = append(deletedFiles, file)
			continue
		}
		if strings.Index(l, " M") == 0 {
			file := strings.ReplaceAll(l, " M ", "")
			if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
				continue
			}
			modifiedFiles = append(modifiedFiles, file)
			continue
		}
	}
	return addedFiles, modifiedFiles, deletedFiles
}

func toMDList(list []string) (result []string) {
	for _, item := range list {
		result = append(result, "* "+item)
	}
	return result
}

func getCommitMessage(addedFiles []string, modifiedFiles []string, deletedFiles []string) string {
	changesCount := len(addedFiles) + len(modifiedFiles) + len(deletedFiles)
	commitLines := []string{fmt.Sprintf("Update %d files", changesCount), ""}
	if len(addedFiles) > 0 {
		commitLines = append(commitLines, "Added files:", strings.Join(toMDList(addedFiles), "\n"))
	}
	if len(modifiedFiles) > 0 {
		commitLines = append(commitLines, "Modified files:", strings.Join(toMDList(modifiedFiles), "\n"))
	}
	if len(deletedFiles) > 0 {
		commitLines = append(commitLines, "Deleted files:", strings.Join(toMDList(deletedFiles), "\n"))
	}

	return strings.Join(commitLines, "\n")
}
