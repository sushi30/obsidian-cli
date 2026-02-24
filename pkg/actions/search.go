package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

func SearchNotes(vault obsidian.VaultManager, note obsidian.NoteManager, uri obsidian.UriManager, fuzzyFinder obsidian.FuzzyFinderManager, useEditor bool, metadataFilters map[string]string) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

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

	if useEditor {
		fmt.Printf("Opening note: %s\n", notes[index])
		filePath := filepath.Join(vaultPath, notes[index])
		return obsidian.OpenInEditor(filePath)
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"file":  notes[index],
		"vault": vaultName,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}

	return nil
}

func filterNotesByMetadata(vaultPath string, notes []string, filters map[string]string) ([]string, error) {
	var filtered []string

	for _, note := range notes {
		fullPath := filepath.Join(vaultPath, note)

		// Read file content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		// Parse frontmatter
		fm, _, err := frontmatter.Parse(string(content))
		if err != nil || fm == nil {
			continue
		}

		// Check if frontmatter matches all filters
		if frontmatter.MatchesFilter(fm, filters) {
			filtered = append(filtered, note)
		}
	}

	return filtered, nil
}
