package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

func SearchNotes(vault obsidian.VaultManager, note obsidian.NoteManager, fuzzyFinder obsidian.FuzzyFinderManager, metadataFilters map[string]string) error {
	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	notes, err := note.GetNotesList(vaultPath)
	if err != nil {
		return err
	}

	// Apply metadata filtering if filters are provided
	if len(metadataFilters) > 0 {
		notes, err = filterNotesByMetadata(vaultPath, notes, metadataFilters)
		if err != nil {
			return err
		}
	}

	if len(notes) == 0 {
		return fmt.Errorf("no notes found matching the criteria")
	}

	index, err := fuzzyFinder.Find(notes, func(i int) string {
		return notes[i]
	})

	if err != nil {
		return err
	}

	fmt.Println(notes[index])
	return nil
}

func filterNotesByMetadata(vaultPath string, notes []string, filters map[string]string) ([]string, error) {
	var filtered []string

	for _, note := range notes {
		fullPath := filepath.Join(vaultPath, note)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		fm, _, err := frontmatter.Parse(string(content))
		if err != nil || fm == nil {
			continue
		}

		if frontmatter.MatchesFilter(fm, filters) {
			filtered = append(filtered, note)
		}
	}

	return filtered, nil
}
