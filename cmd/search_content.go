package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
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
		err := actions.SearchNotesContent(&vault, &note, &fuzzyFinder, searchTerm)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchContentCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	rootCmd.AddCommand(searchContentCmd)
}
