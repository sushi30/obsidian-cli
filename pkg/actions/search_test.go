package actions_test

import (
	"errors"
	"testing"

	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestSearchNotes(t *testing.T) {
	t.Run("Successful search note", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		err := actions.SearchNotes(&vault, &note, &fuzzyFinder, nil)
		assert.NoError(t, err, "Expected no error")
	})

	t.Run("fuzzy find returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{
			FindErr: errors.New("Fuzzy find error"),
		}
		err := actions.SearchNotes(&vault, &note, &fuzzyFinder, nil)
		assert.Equal(t, err, fuzzyFinder.FindErr)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			PathError: errors.New("Failed to get vault path"),
		}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		err := actions.SearchNotes(&vault, &note, &fuzzyFinder, nil)
		assert.Equal(t, err, vault.PathError)
	})
}
