package cmd

import (
	"fmt"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

// ResolveNoteName checks if noteName is @daily and resolves it, otherwise returns as-is
func ResolveNoteName(vault *obsidian.Vault, noteName string) (string, error) {
	if !obsidian.IsDailyReference(noteName) {
		return noteName, nil
	}
	return vault.ResolveDailyNote()
}

// WrapDailyNoteError wraps an error with a helpful message if the original
// noteName was @daily and the note doesn't exist
func WrapDailyNoteError(originalNoteName string, err error) error {
	if err == nil {
		return nil
	}
	if obsidian.IsDailyReference(originalNoteName) && err.Error() == obsidian.NoteDoesNotExistError {
		return fmt.Errorf("%s\nYou can create today's daily note with: obsidian create \"@daily\"", err.Error())
	}
	return err
}
