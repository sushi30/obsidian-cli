package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Fuzzy searches and opens note in vault",
	Long: `Fuzzy searches and opens note in vault.

Examples:
  # Search all notes
  obsidian-cli search

  # Search notes with status=active in frontmatter
  obsidian-cli search --meta status=active

  # Search notes with multiple filters
  obsidian-cli search --meta status=active --meta type=project`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		note := obsidian.Note{}
		uri := obsidian.Uri{}
		fuzzyFinder := obsidian.FuzzyFinder{}
		useEditor, err := cmd.Flags().GetBool("editor")
		if err != nil {
			log.Fatalf("failed to retrieve 'editor' flag: %v", err)
		}
		metadataFlags, _ := cmd.Flags().GetStringSlice("meta")

		var metadataFilters map[string]string
		if len(metadataFlags) > 0 {
			metadataFilters, err = frontmatter.ParseFilters(metadataFlags)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = actions.SearchNotes(&vault, &note, &uri, &fuzzyFinder, useEditor, metadataFilters)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	searchCmd.Flags().BoolP("editor", "e", false, "open in editor instead of Obsidian")
	searchCmd.Flags().StringSliceP("meta", "m", []string{}, "filter by frontmatter metadata (key=value)")
	rootCmd.AddCommand(searchCmd)
}
