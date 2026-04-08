package actions

import (
	"fmt"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

func SearchNotesContent(vault obsidian.VaultManager, note obsidian.NoteManager, fuzzyFinder obsidian.FuzzyFinderManager, searchTerm string, metadataFilters map[string]string) error {
	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	matches, err := note.SearchNotesWithSnippets(vaultPath, searchTerm)
	if err != nil {
		return err
	}

	if len(metadataFilters) > 0 {
		uniquePaths := make(map[string]bool)
		for _, m := range matches {
			uniquePaths[m.FilePath] = true
		}
		pathList := make([]string, 0, len(uniquePaths))
		for p := range uniquePaths {
			pathList = append(pathList, p)
		}
		filtered, err := filterNotesByMetadata(vaultPath, pathList, metadataFilters)
		if err != nil {
			return err
		}
		allowedPaths := make(map[string]bool)
		for _, p := range filtered {
			allowedPaths[p] = true
		}
		var filteredMatches []obsidian.NoteMatch
		for _, m := range matches {
			if allowedPaths[m.FilePath] {
				filteredMatches = append(filteredMatches, m)
			}
		}
		matches = filteredMatches
	}

	if len(matches) == 0 {
		fmt.Printf("No notes found containing '%s'\n", searchTerm)
		return nil
	}

	if len(matches) == 1 {
		fmt.Println(matches[0].FilePath)
		return nil
	}

	displayItems := formatMatchesForDisplay(matches)

	index, err := fuzzyFinder.Find(displayItems, func(i int) string {
		return displayItems[i]
	})
	if err != nil {
		return err
	}

	fmt.Println(matches[index].FilePath)
	return nil
}

func formatMatchesForDisplay(matches []obsidian.NoteMatch) []string {
	maxPathLength := calculateMaxPathLength(matches)

	var displayItems []string
	for _, match := range matches {
		displayStr := formatSingleMatch(match, maxPathLength)
		displayItems = append(displayItems, displayStr)
	}

	return displayItems
}

func calculateMaxPathLength(matches []obsidian.NoteMatch) int {
	maxLength := 0
	for _, match := range matches {
		pathWithLine := formatPathWithLine(match)
		if len(pathWithLine) > maxLength {
			maxLength = len(pathWithLine)
		}
	}
	return maxLength
}

func formatPathWithLine(match obsidian.NoteMatch) string {
	if match.LineNumber > 0 {
		return fmt.Sprintf("%s:%d", match.FilePath, match.LineNumber)
	}
	return match.FilePath
}

func formatSingleMatch(match obsidian.NoteMatch, maxPathLength int) string {
	pathWithLine := formatPathWithLine(match)
	if match.LineNumber == 0 {
		// Filename match - show path and indicate it's a filename match
		return fmt.Sprintf("%-*s | %s", maxPathLength, pathWithLine, match.MatchLine)
	}
	// Content match - show path:line | snippet
	return fmt.Sprintf("%-*s | %s", maxPathLength, pathWithLine, match.MatchLine)
}
