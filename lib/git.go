package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("Initializing repository")
	r, err := cloneRepo(config.Repository, config.Path)
	if err == transport.ErrEmptyRemoteRepository {
		log.Println("Repository is empty")
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
	log.Printf("Cloning repository from %s to %s\n", repository, path)
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
	log.Printf("Initializing new repository in %s\n", path)
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
	log.Printf("Adding %s:%s\n", id, path)
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
	return commitFiles(config)
}

func DeleteFile(id string, config Config) error {
	log.Printf("Deleting %s\n", id)
	idPath := p.Join(config.Path, id)
	if err := os.Remove(idPath); err != nil {
		return err
	}
	if !config.SyncOnAction {
		return nil
	}
	return commitFiles(config)
}

func copyFile(from string, to string) error {
	log.Printf("Copying from %s to %s\n", from, to)
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, data, os.ModeAppend)
}

func SyncFiles(files []File, config Config) error {
	log.Println("Syncing files")
	if config.Mode == SymlinkMode {
		return commitFiles(config)
	}
	for _, file := range files {
		idPath := p.Join(config.Path, file.ID)
		if err := copyFile(file.Path, idPath); err != nil {
			return err
		}
	}
	return commitFiles(config)
}

func commitFiles(config Config) error {
	log.Println("Committing files")
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
	addedFiles, modifiedFiles, deletedFiles := parseStatus(s.String())
	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return err
	}
	commitMessage := getCommitMessage(addedFiles, modifiedFiles, deletedFiles)
	commit, err := w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		return err
	}
	if _, err = r.CommitObject(commit); err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "TheAsda",
			Password: "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",
		}},
	)
	return err
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
		if strings.Index(l, "D") == 0 {
			file := strings.ReplaceAll(l, "D ", "")
			if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
				continue
			}
			deletedFiles = append(deletedFiles, file)
			continue
		}
		if strings.Index(l, "M") == 0 {
			file := strings.ReplaceAll(l, "M ", "")
			if strings.Index(file, FilesCollectionName) == 0 || len(file) == 0 {
				continue
			}
			modifiedFiles = append(modifiedFiles, file)
			continue
		}
	}
	return addedFiles, modifiedFiles, deletedFiles
}

func ToMDList(list []string) (result []string) {
	for _, item := range list {
		result = append(result, "* "+item)
	}
	return result
}

func getCommitMessage(addedFiles []string, modifiedFiles []string, deletedFiles []string) string {
	changesCount := len(addedFiles) + len(modifiedFiles) + len(deletedFiles)
	commitLines := []string{fmt.Sprintf("Update %d files", changesCount), ""}
	if len(addedFiles) > 0 {
		commitLines = append(commitLines, "Added files:", strings.Join(ToMDList(addedFiles), "\n"))
	}
	if len(modifiedFiles) > 0 {
		commitLines = append(commitLines, "Modified files:", strings.Join(ToMDList(modifiedFiles), "\n"))
	}
	if len(deletedFiles) > 0 {
		commitLines = append(commitLines, "Deleted files:", strings.Join(ToMDList(deletedFiles), "\n"))
	}

	return strings.Join(commitLines, "\n")
}
