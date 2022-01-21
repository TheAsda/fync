package lib

import (
	"fmt"
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

func (repo *Repo) Exists() bool {
	if !FileExists(repo.config.Path) {
		return false
	}
	cmd := exec.Command("git", "status")
	cmd.Dir = repo.config.Path
	err := cmd.Run()
	return err != nil
}

func (repo *Repo) Clone() error {
	cmd := exec.Command("git", "clone", repo.config.Repository)
	cmd.Dir = repo.config.Path
	return cmd.Run()
}

func (repo *Repo) CommitFiles() error {
	commitMessage, err := repo.getCommitMessage()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "commit", "-A", "-m", commitMessage)
	cmd.Dir = repo.config.Path
	return cmd.Run()
}

func (repo *Repo) getCommitMessage() (string, error) {
	statusCmd := exec.Command("git", "status")
	statusCmd.Output()
	statusCmd.Dir = repo.config.Path
	err := statusCmd.Run()
	if err != nil {
		return "", err
	}
	status, err := statusCmd.Output()
	if err != nil {
		return "", err
	}
	addedFiles, modifiedFiles, deletedFiles := parseStatus(string(status))
	return getCommitMessage(addedFiles, modifiedFiles, deletedFiles), nil
}

// func AddFile(id string, path string, config Config) error {
// 	log.Printf("Adding %s:%s\n", id, path)
// 	idPath := p.Join(config.Path, id)
// 	switch config.Mode {
// 	case SymlinkMode:
// 		if err := os.Symlink(path, idPath); err != nil {
// 			return err
// 		}
// 	case CopyMode:
// 		if err := copyFile(path, idPath); err != nil {
// 			return err
// 		}
// 	default:
// 		return errors.New("Unknown mode")
// 	}
// 	if !config.SyncOnAction {
// 		return nil
// 	}
// 	return commitFiles(config)
// }

// func DeleteFile(id string, config Config) error {
// 	log.Printf("Deleting %s\n", id)
// 	idPath := p.Join(config.Path, id)
// 	if err := os.Remove(idPath); err != nil {
// 		return err
// 	}
// 	if !config.SyncOnAction {
// 		return nil
// 	}
// 	return commitFiles(config)
// }

// func SyncFiles(files []File, config Config) error {
// 	log.Println("Syncing files")
// 	if config.Mode == SymlinkMode {
// 		return commitFiles(config)
// 	}
// 	for _, file := range files {
// 		idPath := p.Join(config.Path, file.ID)
// 		if err := copyFile(file.Path, idPath); err != nil {
// 			return err
// 		}
// 	}
// 	return commitFiles(config)
// }

// func commitFiles(config Config) error {
// 	log.Println("Committing files")
// 	r, err := git.PlainOpen(config.Path)
// 	if err != nil {
// 		return err
// 	}
// 	w, err := r.Worktree()
// 	if err != nil {
// 		return err
// 	}
// 	s, err := w.Status()
// 	if err != nil {
// 		return err
// 	}
// 	if s.IsClean() {
// 		return nil
// 	}
// 	addedFiles, modifiedFiles, deletedFiles := parseStatus(s.String())
// 	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
// 		return err
// 	}
// 	commitMessage := getCommitMessage(addedFiles, modifiedFiles, deletedFiles)
// 	commit, err := w.Commit(commitMessage, &git.CommitOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	if _, err = r.CommitObject(commit); err != nil {
// 		return err
// 	}

// 	err = r.Push(&git.PushOptions{
// 		Auth: &http.BasicAuth{
// 			Username: "TheAsda",
// 			Password: "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",
// 		}},
// 	)
// 	return err
// }

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
