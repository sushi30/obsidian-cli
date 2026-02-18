package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

type vaultStub struct {
	path       string
	defaultErr error
	pathErr    error
}

func (v *vaultStub) DefaultName() (string, error) {
	if v.defaultErr != nil {
		return "", v.defaultErr
	}
	return "example-vault", nil
}

func (v *vaultStub) SetDefaultName(_ string) error {
	return nil
}

func (v *vaultStub) Path() (string, error) {
	if v.pathErr != nil {
		return "", v.pathErr
	}
	return v.path, nil
}

func (v *vaultStub) DailyNotePattern() (string, error) {
	return "", nil
}

func (v *vaultStub) ResolveDailyNote() (string, error) {
	return "", nil
}

func TestListEntries(t *testing.T) {
	t.Run("List vault root", func(t *testing.T) {
		vaultDir := t.TempDir()

		err := os.Mkdir(filepath.Join(vaultDir, "Project Alpha"), 0755)
		assert.NoError(t, err)

		err = os.WriteFile(filepath.Join(vaultDir, "Ideas.md"), []byte(""), 0644)
		assert.NoError(t, err)

		err = os.WriteFile(filepath.Join(vaultDir, "Meeting Notes.md"), []byte(""), 0644)
		assert.NoError(t, err)

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Project Alpha/", "Ideas.md", "Meeting Notes.md"}, entries)
	})

	t.Run("List subdirectory", func(t *testing.T) {
		vaultDir := t.TempDir()
		subDir := filepath.Join(vaultDir, "001 Notes")
		assert.NoError(t, os.Mkdir(subDir, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(subDir, "Daily.md"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "001 Notes"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Daily.md"}, entries)
	})

	t.Run("Rejects path traversal", func(t *testing.T) {
		vault := &vaultStub{path: t.TempDir()}
		_, err := actions.ListEntries(vault, actions.ListParams{Path: "../"})
		assert.ErrorIs(t, err, obsidian.ErrPathTraversal)
	})

	t.Run("Returns error when path is not directory", func(t *testing.T) {
		vaultDir := t.TempDir()
		fileName := "Ideas.md"
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, fileName), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		_, err := actions.ListEntries(vault, actions.ListParams{Path: fileName})
		assert.EqualError(t, err, obsidian.VaultAccessError)
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		vault := &vaultStub{path: t.TempDir(), defaultErr: errors.New("default error")}
		_, err := actions.ListEntries(vault, actions.ListParams{})
		assert.Equal(t, vault.defaultErr, err)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		vault := &vaultStub{pathErr: errors.New("path error")}
		_, err := actions.ListEntries(vault, actions.ListParams{})
		assert.Equal(t, vault.pathErr, err)
	})

	t.Run("List with full path at vault root", func(t *testing.T) {
		vaultDir := t.TempDir()

		err := os.WriteFile(filepath.Join(vaultDir, "Notes.md"), []byte(""), 0644)
		assert.NoError(t, err)

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{FullPath: true})
		assert.NoError(t, err)
		assert.Equal(t, []string{vaultDir + "/Notes.md"}, entries)
	})

	t.Run("List with full path in subdirectory", func(t *testing.T) {
		vaultDir := t.TempDir()
		subDir := filepath.Join(vaultDir, "Projects")
		assert.NoError(t, os.Mkdir(subDir, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(subDir, "Todo.md"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "Projects", FullPath: true})
		assert.NoError(t, err)
		assert.Equal(t, []string{vaultDir + "/Projects/Todo.md"}, entries)
	})

	t.Run("Glob pattern with * matches files", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "notes.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "ideas.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "readme.txt"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "*.md"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ideas.md", "notes.md"}, entries)
	})

	t.Run("Glob pattern with ** matches recursively", func(t *testing.T) {
		vaultDir := t.TempDir()
		subDir := filepath.Join(vaultDir, "Projects")
		assert.NoError(t, os.Mkdir(subDir, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "root.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(subDir, "project.md"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "**/*.md"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Projects/project.md", "root.md"}, entries)
	})

	t.Run("Glob with full path", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "note.md"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "*.md", FullPath: true})
		assert.NoError(t, err)
		assert.Equal(t, []string{vaultDir + "/note.md"}, entries)
	})

	t.Run("Glob rejects path traversal", func(t *testing.T) {
		vault := &vaultStub{path: t.TempDir()}
		_, err := actions.ListEntries(vault, actions.ListParams{Path: "../*.md"})
		assert.ErrorIs(t, err, obsidian.ErrPathTraversal)
	})
}
