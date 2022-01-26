package repo

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"theasda/fync/pkg/config"
	"theasda/fync/pkg/utils"

	"github.com/sirupsen/logrus"
)

type Repo struct {
	config config.Config
}

func NewRepo(config config.Config) *Repo {
	logrus.Debug("Initializing repo")
	return &Repo{
		config: config,
	}
}

func (repo Repo) Exists() bool {
	if !utils.FileExists(repo.config.Path) {
		logrus.Debug("Repo path does not exist")
		return false
	}
	statusBuffer, err := repo.runCommand("git", "status")
	exists := err == nil && !strings.Contains(statusBuffer.String(), "not a git repository")
	if exists {
		logrus.Debug("Repo exists")
	} else {
		logrus.Debug("Repo does not exist")
	}
	return exists
}

func (repo Repo) Clone() error {
	logrus.Debug("Cloning repository")
	if err := os.MkdirAll(repo.config.Path, 0644); err != nil {
		return err
	}
	_, err := repo.runCommand("git", "clone", repo.config.Repository, ".")
	return err
}

func (repo Repo) StageFiles() error {
	logrus.Debug("Staging files")
	_, err := repo.runCommand("git", "add", "-A")
	return err
}

func (repo Repo) CommitFiles() error {
	commitMessage, err := repo.getCommitMessage()
	if err != nil {
		return err
	}

	if err = repo.StageFiles(); err != nil {
		return err
	}

	logrus.Info("Committing files")
	_, err = repo.runCommand("git", "commit", "-m", commitMessage)
	return err
}

func (repo Repo) Push() error {
	logrus.Info("Pushing")
	_, err := repo.runCommand("git", "push")
	return err
}

func (repo *Repo) UpdateRepo() error {
	if err := repo.CommitFiles(); err != nil {
		return err
	}
	return repo.Push()
}

func (repo Repo) getCommitMessage() (string, error) {
	statusBuffer, err := repo.runCommand("git", "status", "-sb")
	if err != nil {
		return "", err
	}
	status := statusBuffer.String()
	addedFiles, modifiedFiles, deletedFiles := parseStatus(status)
	if len(addedFiles)+len(modifiedFiles)+len(deletedFiles) == 0 {
		return "", errors.New("no changes in files")
	}
	return compileCommitMessage(addedFiles, modifiedFiles, deletedFiles), nil
}

func (repo Repo) runCommand(name string, arg ...string) (stdout bytes.Buffer, err error) {
	cmd := exec.Command("git", "status", "-sb")
	cmd.Path = repo.config.Path
	cmd.Stdout = &stdout
	err = cmd.Run()
	return stdout, err
}
