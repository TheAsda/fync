package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	p "path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func InitRepo(config Config) error {
	r, err := cloneRepo(config.Repository, config.Path)
	if err == transport.ErrEmptyRemoteRepository {
		r, err = createNewRepo(config.Repository, config.Path)
	}
	if err != nil {
		return err
	}
	blobs, err := r.BlobObjects()
	if err != nil {
		return err
	}
	blobs.ForEach(func(b *object.Blob) error {
		println("blob")
		return nil
	})
	return nil
}

func cloneRepo(repository string, path string) (*git.Repository, error) {
	return git.PlainClone(path, false, &git.CloneOptions{
		URL:      repository,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: "TheAsda",
			Password: "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",
		},
	})
}

func createNewRepo(repository string, path string) (*git.Repository, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	r, err := git.PlainInit(path, false)
	if err != nil {
		return nil, err
	}
	_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{repository}})
	return r, err
}

func AddFile(id string, path string, config Config) error {
	idPath := p.Join(config.Path, id)
	switch config.Mode {
	case SymlinkMode:
		if err := os.Symlink(path, idPath); err != nil {
			return err
		}
	case CopyMode:
		if err := copyFile(path, idPath); err != nil {
			return err
		}
	default:
		return errors.New("Unknown mode")
	}
	if !config.SyncOnAction {
		return nil
	}
	return CommitFiles(config)
}

func copyFile(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, data, os.ModeAppend)
}

func CommitFiles(config Config) error {
	r, err := git.PlainOpen(config.Path)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	s, err := w.Status()
	if err != nil {
		return err
	}
	if s.IsClean() {
		return nil
	}
	changedFiles := parseStatus(s.String())
	for i := range changedFiles {
		changedFiles[i] = "* " + changedFiles[i]
	}
	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return err
	}
	commit, err := w.Commit(fmt.Sprintf("Update %v files\n\nUpdated files:\n%s", len(changedFiles), strings.Join(changedFiles, "\n")), &git.CommitOptions{})
	if err != nil {
		return err
	}
	_, err = r.CommitObject(commit)
	return err
}

func parseStatus(status string) []string {
	lines := strings.Split(status, "\n")
	var changedFiles []string
	for _, l := range lines {
		file := strings.ReplaceAll(l, "?? ", "")
		if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
			continue
		}
		changedFiles = append(changedFiles, file)
	}
	return changedFiles
}
