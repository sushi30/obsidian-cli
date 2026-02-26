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
	Short:   "Fuzzy searches and prints note path",
	Long: `Fuzzy searches and prints note path.

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
		fuzzyFinder := obsidian.FuzzyFinder{}
		metadataFlags, _ := cmd.Flags().GetStringSlice("meta")

		var metadataFilters map[string]string
		var err error
		if len(metadataFlags) > 0 {
			metadataFilters, err = frontmatter.ParseFilters(metadataFlags)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = actions.SearchNotes(&vault, &note, &fuzzyFinder, metadataFilters)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	searchCmd.Flags().StringSliceP("meta", "m", []string{}, "filter by frontmatter metadata (key=value)")
	rootCmd.AddCommand(searchCmd)
}
