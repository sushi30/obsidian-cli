package actions_test

import (
	"errors"
	"testing"

	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestAppendToNote(t *testing.T) {
	t.Run("Successfully append content to note", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "Existing content",
		}

		output, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "test-note",
			Content:  "\nNew content",
		})

		assert.NoError(t, err)
		assert.Contains(t, output, "Appended content to test-note")
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			DefaultNameErr: errors.New("vault error"),
		}
		note := mocks.MockNoteManager{}

		_, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "test-note",
			Content:  "content",
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

		_, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "test-note",
			Content:  "content",
		})

		assert.Error(t, err)
		assert.Equal(t, "path error", err.Error())
	})

	t.Run("Note not found error propagates", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			GetContentsError: errors.New("note not found"),
		}

		_, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "nonexistent",
			Content:  "content",
		})

		assert.Error(t, err)
		assert.Equal(t, "note not found", err.Error())
	})

	t.Run("Write error propagates", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents:         "existing",
			SetContentsError: errors.New("write error"),
		}

		_, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "test-note",
			Content:  "content",
		})

		assert.Error(t, err)
		assert.Equal(t, "write error", err.Error())
	})

	t.Run("Empty content appends nothing but succeeds", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			Contents: "Existing content",
		}

		output, err := actions.AppendToNote(&vault, &note, actions.AppendParams{
			NoteName: "test-note",
			Content:  "",
		})

		assert.NoError(t, err)
		assert.Contains(t, output, "Appended content")
	})
}
