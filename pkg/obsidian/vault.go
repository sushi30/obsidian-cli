package obsidian

type CliConfig struct {
	DefaultVaultName string `json:"default_vault_name"`
	DailyNotePattern string `json:"daily_note_pattern,omitempty"`
}

type ObsidianVaultConfig struct {
	Vaults map[string]struct {
		Path string `json:"path"`
	} `json:"vaults"`
}

type VaultManager interface {
	DefaultName() (string, error)
	SetDefaultName(name string) error
	Path() (string, error)
	DailyNotePattern() (string, error)
	ResolveDailyNote() (string, error)
}

type Vault struct {
	Name string
}
