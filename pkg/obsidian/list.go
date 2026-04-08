package obsidian

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
)

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

func ListEntriesWithFilter(vaultPath, relativePath string, filters map[string]string) ([]string, error) {
	targetPath := vaultPath
	if strings.TrimSpace(relativePath) != "" {
		validatedPath, err := ValidatePath(vaultPath, relativePath)
		if err != nil {
			return nil, err
		}
		targetPath = validatedPath
	}

	var results []string
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != targetPath {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Size() > 10*1024*1024 {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(content)
		if !frontmatter.HasFrontmatter(text) {
			return nil
		}
		fm, _, err := frontmatter.Parse(text)
		if err != nil {
			return nil
		}
		if !frontmatter.MatchesFilters(fm, filters) {
			return nil
		}

		relPath, err := filepath.Rel(targetPath, path)
		if err != nil {
			return err
		}
		results = append(results, relPath)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(results)
	return results, nil
}
