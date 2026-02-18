package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type EditParams struct {
	NoteName   string
	OldString  string
	NewString  string
	ReplaceAll bool
}

func EditNote(vault obsidian.VaultManager, note obsidian.NoteManager, params EditParams) (string, error) {
	if params.OldString == params.NewString {
		return "", errors.New("old string and new string must be different")
	}

	if params.OldString == "" {
		return "", errors.New("old string cannot be empty")
	}

	_, err := vault.DefaultName()
	if err != nil {
		return "", err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return "", err
	}

	contents, err := note.GetContents(vaultPath, params.NoteName)
	if err != nil {
		return "", err
	}

	if !strings.Contains(contents, params.OldString) {
		return "", fmt.Errorf("string not found in note: %q", params.OldString)
	}

	var updatedContent string
	var count int
	if params.ReplaceAll {
		count = strings.Count(contents, params.OldString)
		updatedContent = strings.ReplaceAll(contents, params.OldString, params.NewString)
	} else {
		count = 1
		updatedContent = strings.Replace(contents, params.OldString, params.NewString, 1)

		if strings.Count(contents, params.OldString) > 1 {
			return "", fmt.Errorf("string appears %d times in note. Use --all flag to replace all occurrences or make the old string more specific", strings.Count(contents, params.OldString))
		}
	}

	err = note.SetContents(vaultPath, params.NoteName, updatedContent)
	if err != nil {
		return "", err
	}

	occurrences := "occurrence"
	if count > 1 {
		occurrences = "occurrences"
	}

	return fmt.Sprintf("Replaced %d %s in %s", count, occurrences, params.NoteName), nil
}
