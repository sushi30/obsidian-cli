package obsidian

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func ContainsGlob(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

func GlobEntries(vaultPath, pattern string) ([]string, error) {
	if strings.Contains(pattern, "..") {
		return nil, ErrPathTraversal
	}

	fullPattern := filepath.Join(vaultPath, pattern)

	matches, err := doublestar.FilepathGlob(fullPattern)
	if err != nil {
		return nil, errors.New(VaultReadError)
	}

	results := make([]string, 0, len(matches))
	for _, match := range matches {
		relPath, err := filepath.Rel(vaultPath, match)
		if err != nil {
			continue
		}
		if strings.HasPrefix(relPath, ".") || strings.Contains(relPath, "/.") {
			continue
		}
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.IsDir() {
			results = append(results, relPath+"/")
		} else {
			results = append(results, relPath)
		}
	}

	sort.Strings(results)
	return results, nil
}

func ListEntries(vaultPath, relativePath string) ([]string, error) {
	targetPath := vaultPath
	if strings.TrimSpace(relativePath) != "" {
		validatedPath, err := ValidatePath(vaultPath, relativePath)
		if err != nil {
			return nil, err
		}
		targetPath = validatedPath
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, errors.New(VaultAccessError)
	}
	if !info.IsDir() {
		return nil, errors.New(VaultAccessError)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, errors.New(VaultReadError)
	}

	dirs := make([]string, 0, len(entries))
	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, name+"/")
			continue
		}
		files = append(files, name)
	}

	sort.Strings(dirs)
	sort.Strings(files)

	return append(dirs, files...), nil
}
