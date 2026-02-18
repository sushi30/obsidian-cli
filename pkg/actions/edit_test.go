package actions_test

import (
	"errors"
	"testing"

	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestEditNote(t *testing.T) {

	t.Run("Successfully replace single occurrence", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "This is a test note.\nIt has multiple lines.",
		}

		output, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "test note",
			NewString:  "sample note",
			ReplaceAll: false,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Replaced 1 occurrence in test", output)
	})

	t.Run("Successfully replace all occurrences", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "This is a test note.\nIt has multiple test lines.\nThis is another test case.",
		}

		output, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "test",
			NewString:  "example",
			ReplaceAll: true,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Replaced 3 occurrences in test", output)
	})

	t.Run("Fail when strings are identical", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{Contents: "content"}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "same",
			NewString:  "same",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "old string and new string must be different")
	})

	t.Run("Fail when old string is empty", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{Contents: "content"}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "old string cannot be empty")
	})

	t.Run("Fail when string not found", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "This is a test note.",
		}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "nonexistent",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "string not found in note")
	})

	t.Run("Fail when multiple occurrences without replaceAll", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "This is a test note.\nAnother test here.",
		}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "test",
			NewString:  "example",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "string appears 2 times in note")
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			DefaultNameErr: errors.New("vault error"),
		}
		note := mocks.MockNoteManager{}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "old",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Equal(t, "vault error", err.Error())
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			Name:      "myVault",
			PathError: errors.New("path error"),
		}
		note := mocks.MockNoteManager{}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "old",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Equal(t, "path error", err.Error())
	})

	t.Run("Note not found error propagates", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			GetContentsError: errors.New("note not found"),
		}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "nonexistent",
			OldString:  "old",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Equal(t, "note not found", err.Error())
	})

	t.Run("Write error propagates", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents:         "old text",
			SetContentsError: errors.New("write error"),
		}

		_, err := actions.EditNote(&vault, &note, actions.EditParams{
			NoteName:   "test",
			OldString:  "old",
			NewString:  "new",
			ReplaceAll: false,
		})

		assert.Error(t, err)
		assert.Equal(t, "write error", err.Error())
	})
}
