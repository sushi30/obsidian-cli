package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

type vaultStubForSearch struct {
	path string
}

func (v *vaultStubForSearch) DefaultName() (string, error)      { return "test-vault", nil }
func (v *vaultStubForSearch) SetDefaultName(_ string) error     { return nil }
func (v *vaultStubForSearch) Path() (string, error)             { return v.path, nil }
func (v *vaultStubForSearch) DailyNotePattern() (string, error) { return "", nil }
func (v *vaultStubForSearch) ResolveDailyNote() (string, error) { return "", nil }

type CustomMockNoteForSingleMatch struct{}

func (m *CustomMockNoteForSingleMatch) Delete(string) error                        { return nil }
func (m *CustomMockNoteForSingleMatch) Move(string, string) error                  { return nil }
func (m *CustomMockNoteForSingleMatch) UpdateLinks(string, string, string) error   { return nil }
func (m *CustomMockNoteForSingleMatch) GetContents(string, string) (string, error) { return "", nil }
func (m *CustomMockNoteForSingleMatch) SetContents(string, string, string) error   { return nil }
func (m *CustomMockNoteForSingleMatch) GetNotesList(string) ([]string, error)      { return nil, nil }
func (m *CustomMockNoteForSingleMatch) SearchNotesWithSnippets(string, string) ([]obsidian.NoteMatch, error) {
	return []obsidian.NoteMatch{
		{FilePath: "test-note.md", LineNumber: 5, MatchLine: "test content"},
	}, nil
}
func (m *CustomMockNoteForSingleMatch) FindBacklinks(string, string) ([]obsidian.NoteMatch, error) {
	return nil, nil
}

func TestSearchNotesContent(t *testing.T) {
	t.Run("Successful content search with multiple matches", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, "test", nil)
		assert.NoError(t, err)
	})

	t.Run("Successful content search with single match", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := &CustomMockNoteForSingleMatch{}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(&vault, note, &fuzzyFinder, "test", nil)
		assert.NoError(t, err)
	})

	t.Run("No matches found", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{NoMatches: true}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, "nonexistent", nil)
		assert.NoError(t, err)
	})

	t.Run("SearchNotesWithSnippets returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{
			GetContentsError: errors.New("search failed"),
		}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, "test", nil)
		assert.Error(t, err)
	})

	t.Run("vault.Path returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			PathError: errors.New("vault path error"),
		}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, "test", nil)
		assert.Error(t, err)
	})

	t.Run("fuzzy finder returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{
			FindErr: errors.New("fuzzy finder error"),
		}

		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, "test", nil)
		assert.Error(t, err)
	})

	t.Run("Metadata filter excludes non-matching notes", func(t *testing.T) {
		vaultDir := t.TempDir()
		os.WriteFile(filepath.Join(vaultDir, "match.md"),
			[]byte("---\nstatus: done\n---\nhas the keyword"), 0644)
		os.WriteFile(filepath.Join(vaultDir, "nomatch.md"),
			[]byte("---\nstatus: draft\n---\nhas the keyword"), 0644)

		vault := &vaultStubForSearch{path: vaultDir}
		note := &obsidian.Note{}
		fuzzyFinder := mocks.MockFuzzyFinder{}

		err := actions.SearchNotesContent(vault, note, &fuzzyFinder, "keyword",
			map[string]string{"status": "done"})
		assert.NoError(t, err)
	})
}
