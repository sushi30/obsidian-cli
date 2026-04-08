package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"

	"github.com/spf13/cobra"
)

var searchContentCmd = &cobra.Command{
	Use:     "search-content [search term]",
	Short:   "Search note content and print matching note path",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"sc"},
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		note := obsidian.Note{}
		fuzzyFinder := obsidian.FuzzyFinder{}

		searchTerm := args[0]
		metadataFlags, _ := cmd.Flags().GetStringSlice("meta")

		var metadataFilters map[string]string
		var err error
		if len(metadataFlags) > 0 {
			metadataFilters, err = frontmatter.ParseFilters(metadataFlags)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = actions.SearchNotesContent(&vault, &note, &fuzzyFinder, searchTerm, metadataFilters)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchContentCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	searchContentCmd.Flags().StringSliceP("meta", "m", []string{}, "filter by frontmatter metadata (key=value)")
	rootCmd.AddCommand(searchContentCmd)
}
