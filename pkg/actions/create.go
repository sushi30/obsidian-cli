package actions

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type CreateParams struct {
	NoteName        string
	ShouldAppend    bool
	ShouldOverwrite bool
	Content         string
	ShouldOpen      bool
	UseEditor       bool
}

func CreateNote(vault obsidian.VaultManager, uri obsidian.UriManager, params CreateParams) error {
	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	normalizedContent := NormalizeContent(params.Content)

	// Build the full file path
	filePath, err := obsidian.ValidatePath(vaultPath, obsidian.AddMdSuffix(params.NoteName))
	if err != nil {
		return err
	}

	// Create parent directories if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	_, err = os.Stat(filePath)
	fileExists := err == nil

	if fileExists && !params.ShouldOverwrite && !params.ShouldAppend {
		// File exists and no overwrite/append flag - let Obsidian handle it via URI
		// (it will show an error or prompt)
		vaultName, err := vault.DefaultName()
		if err != nil {
			return err
		}
		obsidianUri := uri.Construct(ObsCreateUrl, map[string]string{
			"vault":     vaultName,
			"file":      params.NoteName,
			"content":   normalizedContent,
			"silent":    strconv.FormatBool(!params.ShouldOpen),
		})
		return uri.Execute(obsidianUri)
	}

	// Write directly to file
	if params.ShouldAppend && fileExists {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(normalizedContent); err != nil {
			return err
		}
	} else {
		if err := os.WriteFile(filePath, []byte(normalizedContent), 0644); err != nil {
			return err
		}
	}

	// Open in Obsidian or editor if requested
	if params.ShouldOpen {
		if params.UseEditor {
			return obsidian.OpenInEditor(filePath)
		}
		vaultName, err := vault.DefaultName()
		if err != nil {
			return err
		}
		obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
			"vault": vaultName,
			"file":  params.NoteName,
		})
		return uri.Execute(obsidianUri)
	}

	return nil
}

func NormalizeContent(content string) string {
	replacer := strings.NewReplacer(
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
		"\\\\", "\\",
		"\\\"", "\"",
		"\\'", "'",
	)
	return replacer.Replace(content)
}
