package actions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	t.Run("Successful create note", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:  "note",
			Content:   "test content",
			UseEditor: false,
		})
		// Assert
		assert.NoError(t, err, "Expected no error")
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{
			PathError: errors.New("Failed to get vault path"),
		}
		// Act
		err := actions.CreateNote(&vault, &mocks.MockUriManager{}, actions.CreateParams{
			NoteName:  "note-name",
			Content:   "test",
			UseEditor: false,
		})
		// Assert
		assert.Equal(t, vault.PathError, err)
	})

	t.Run("uri.Execute returns error when opening note", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		uri := mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note-name",
			Content:    "test",
			ShouldOpen: true,
			UseEditor:  false,
		})
		// Assert
		assert.Equal(t, uri.ExecuteErr, err)
	})

	t.Run("Successful create note with editor flag and open", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Set EDITOR to a command that will succeed
		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			Content:    "test",
			ShouldOpen: true,
			UseEditor:  true,
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Create note with editor flag fails when editor fails", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Set EDITOR to a command that will fail
		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "false")

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			Content:    "test",
			ShouldOpen: true,
			UseEditor:  true,
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("Create note with editor flag without open does not use editor", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Act - UseEditor is true but ShouldOpen is false
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			Content:    "test",
			ShouldOpen: false,
			UseEditor:  true,
		})

		// Assert - should succeed and write the file
		assert.NoError(t, err)
	})
}

func TestCreateNote_FileExists(t *testing.T) {
	t.Run("returns error when file exists without overwrite flag", func(t *testing.T) {
		// Arrange - create temp dir with existing file
		tmpDir := t.TempDir()
		existingFile := tmpDir + "/existing-note.md"
		err := os.WriteFile(existingFile, []byte("original content"), 0644)
		assert.NoError(t, err)

		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Act - try to create note without overwrite flag
		err = actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:        "existing-note",
			Content:         "new content",
			ShouldOverwrite: false,
		})

		// Assert - should return an error, not silently fail
		if assert.Error(t, err, "Expected error when file exists without overwrite flag") {
			assert.Contains(t, err.Error(), "already exists")
		}

		// Verify original content unchanged
		content, _ := os.ReadFile(existingFile)
		assert.Equal(t, "original content", string(content))
	})

	t.Run("overwrites file when overwrite flag is set", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		existingFile := tmpDir + "/existing-note.md"
		err := os.WriteFile(existingFile, []byte("original content"), 0644)
		assert.NoError(t, err)

		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Act - create note WITH overwrite flag
		err = actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:        "existing-note",
			Content:         "new content",
			ShouldOverwrite: true,
		})

		// Assert
		assert.NoError(t, err)
		content, _ := os.ReadFile(existingFile)
		assert.Equal(t, "new content", string(content))
	})

	t.Run("creates new file when it does not exist", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", VaultPath: tmpDir}
		uri := mocks.MockUriManager{}

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "new-note",
			Content:  "some content",
		})

		// Assert
		assert.NoError(t, err)
		content, err := os.ReadFile(tmpDir + "/new-note.md")
		assert.NoError(t, err)
		assert.Equal(t, "some content", string(content))
	})
}

func TestNormalizeContent(t *testing.T) {
	t.Run("Replaces escape sequences with actual characters", func(t *testing.T) {
		// Arrange
		input := "Hello\\nWorld\\tTabbed\\rReturn\\\"Quote\\'SingleQuote\\\\Backslash"
		expected := "Hello\nWorld\tTabbed\rReturn\"Quote'SingleQuote\\Backslash"

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "The content should have the escape sequences replaced correctly")
	})

	t.Run("Handles empty input", func(t *testing.T) {
		// Arrange
		input := ""
		expected := ""

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "Empty input should return empty output")
	})

	t.Run("No escape sequences in input", func(t *testing.T) {
		// Arrange
		input := "Plain text with no escapes"
		expected := "Plain text with no escapes"

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "Content without escape sequences should remain unchanged")
	})
}
