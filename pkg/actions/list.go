package actions

import (
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type ListParams struct {
	Path     string
	FullPath bool
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
