package actions

import (
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type ListParams struct {
	Path  string
	Where map[string]string
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

	if len(params.Where) > 0 {
		return obsidian.ListEntriesWithFilter(vaultPath, params.Path, params.Where)
	}

	return obsidian.ListEntries(vaultPath, params.Path)
}
