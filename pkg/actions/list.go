package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type ListParams struct {
	Path            string
	FullPath        bool
	MetadataFilters map[string]string
}

func ListEntries(vault obsidian.VaultManager, params ListParams) ([]string, error) {
	_, err := vault.DefaultName()
	if err != nil {
		return nil, err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return nil, err
	}

	var entries []string
	if obsidian.ContainsGlob(params.Path) {
		entries, err = obsidian.GlobEntries(vaultPath, params.Path)
	} else {
		entries, err = obsidian.ListEntries(vaultPath, params.Path)
	}
	if err != nil {
		return nil, err
	}

	// Apply metadata filtering if filters are provided
	if len(params.MetadataFilters) > 0 {
		entries, err = filterEntriesByMetadata(vaultPath, entries, params.MetadataFilters)
		if err != nil {
			return nil, err
		}
	}

	if params.FullPath {
		basePath := vaultPath
		if params.Path != "" && !obsidian.ContainsGlob(params.Path) {
			basePath = vaultPath + "/" + params.Path
		}
		for i, entry := range entries {
			entries[i] = basePath + "/" + entry
		}
	}

	return entries, nil
}

func filterEntriesByMetadata(vaultPath string, entries []string, filters map[string]string) ([]string, error) {
	var filtered []string

	for _, entry := range entries {
		fullPath := filepath.Join(vaultPath, entry)

		// Only process markdown files
		if !strings.HasSuffix(entry, ".md") {
			continue
		}

		// Check if it's a file (not a directory)
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			continue
		}

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
			filtered = append(filtered, entry)
		}
	}

	return filtered, nil
}
