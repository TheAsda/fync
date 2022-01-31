package files_processor

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func promptOverride(file string, path string) (bool, error) {
	overridePrompt := promptui.Select{
		Label: fmt.Sprintf("Do you want to override %s with %s", path, file),
		Items: []string{"yes", "no"},
	}
	_, answer, err := overridePrompt.Run()
	if err != nil {
		return false, err
	}
	return answer == "yes", nil
}
