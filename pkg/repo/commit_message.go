package repo

import (
	"fmt"
	"strings"
	"theasda/fync/pkg/config"
)

func compileCommitMessage(addedFiles []string, modifiedFiles []string, deletedFiles []string) string {
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

func parseStatus(status string) (addedFiles []string, modifiedFiles []string, deletedFiles []string) {
	lines := strings.Split(status, "\n")
	for _, l := range lines {
		if strings.Index(l, "??") == 0 {
			file := strings.ReplaceAll(l, "?? ", "")
			if strings.Index(file, config.StorageFileName) == 0 || len(file) == 0 {
				continue
			}
			addedFiles = append(addedFiles, file)
			continue
		}
		if strings.Index(l, " D") == 0 {
			file := strings.ReplaceAll(l, " D ", "")
			if strings.Index(file, config.StorageFileName) == 0 || len(file) == 0 {
				continue
			}
			deletedFiles = append(deletedFiles, file)
			continue
		}
		if strings.Index(l, " M") == 0 {
			file := strings.ReplaceAll(l, " M ", "")
			if strings.Index(file, config.StorageFileName) == 0 || len(file) == 0 {
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
