package obsidian

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEntries(t *testing.T) {
	t.Run("Lists directories before files with sorting", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, "Project Alpha"), 0755))
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, "Ideas"), 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Meeting Notes.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Ideas.md"), []byte(""), 0644))

		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Equal(t, []string{"Ideas/", "Project Alpha/", "Ideas.md", "Meeting Notes.md"}, entries)
	})

	t.Run("Filters hidden files and folders", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, ".obsidian"), 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, ".hidden.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Visible.md"), []byte(""), 0644))

		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Equal(t, []string{"Visible.md"}, entries)
	})

	t.Run("Empty directory returns empty list", func(t *testing.T) {
		vaultDir := t.TempDir()
		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("Non-directory path returns vault access error", func(t *testing.T) {
		vaultDir := t.TempDir()
		filePath := filepath.Join(vaultDir, "Note.md")
		assert.NoError(t, os.WriteFile(filePath, []byte(""), 0644))

		_, err := ListEntries(vaultDir, "Note.md")
		assert.EqualError(t, err, VaultAccessError)
	})

	t.Run("Path traversal is rejected", func(t *testing.T) {
		vaultDir := t.TempDir()
		_, err := ListEntries(vaultDir, "../")
		assert.ErrorIs(t, err, ErrPathTraversal)
	})
}

func writeNote(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	assert.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	assert.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func TestListEntriesWithFilter(t *testing.T) {
	t.Run("Matches single filter", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "done.md", "---\nstatus: done\n---\nContent")
		writeNote(t, vaultDir, "draft.md", "---\nstatus: draft\n---\nContent")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"done.md"}, results)
	})

	t.Run("Matches multiple filters AND", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "match.md", "---\nstatus: done\ntags:\n  - work\n---\nContent")
		writeNote(t, vaultDir, "partial.md", "---\nstatus: done\ntags:\n  - personal\n---\nContent")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done", "tags": "work"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"match.md"}, results)
	})

	t.Run("Excludes non-matching notes", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "a.md", "---\nstatus: done\n---\n")
		writeNote(t, vaultDir, "b.md", "---\nstatus: draft\n---\n")
		writeNote(t, vaultDir, "c.md", "---\nstatus: done\n---\n")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		sort.Strings(results)
		assert.Equal(t, []string{"a.md", "c.md"}, results)
	})

	t.Run("Excludes directories", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, "subdir"), 0755))
		writeNote(t, vaultDir, "note.md", "---\nstatus: done\n---\n")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"note.md"}, results)
	})

	t.Run("Excludes hidden files", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, ".hidden.md", "---\nstatus: done\n---\n")
		writeNote(t, vaultDir, "visible.md", "---\nstatus: done\n---\n")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"visible.md"}, results)
	})

	t.Run("Works with subdirectory path scoping", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "root.md", "---\nstatus: done\n---\n")
		writeNote(t, vaultDir, "sub/nested.md", "---\nstatus: done\n---\n")

		results, err := ListEntriesWithFilter(vaultDir, "sub", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"nested.md"}, results)
	})

	t.Run("Returns empty for no matches", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "note.md", "---\nstatus: draft\n---\n")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Handles notes without frontmatter", func(t *testing.T) {
		vaultDir := t.TempDir()
		writeNote(t, vaultDir, "no-fm.md", "Just plain content")
		writeNote(t, vaultDir, "with-fm.md", "---\nstatus: done\n---\nContent")

		results, err := ListEntriesWithFilter(vaultDir, "", map[string]string{"status": "done"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"with-fm.md"}, results)
	})
}
