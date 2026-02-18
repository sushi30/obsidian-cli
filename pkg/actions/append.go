package actions

import (
	"fmt"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type AppendParams struct {
	NoteName string
	Content  string
}

func AppendToNote(vault obsidian.VaultManager, note obsidian.NoteManager, params AppendParams) (string, error) {
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

	normalizedContent := NormalizeContent(params.Content)
	updatedContent := contents + normalizedContent

	err = note.SetContents(vaultPath, params.NoteName, updatedContent)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Appended content to %s", params.NoteName), nil
}
