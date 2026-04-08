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
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		note := obsidian.Note{}
		uri := obsidian.Uri{}
		fuzzyFinder := obsidian.FuzzyFinder{}
		useEditor, err := cmd.Flags().GetBool("editor")
		if err != nil {
			log.Fatalf("failed to retrieve 'editor' flag: %v", err)
		}
		whereStr, err := cmd.Flags().GetString("where")
		if err != nil {
			log.Fatalf("failed to retrieve 'where' flag: %v", err)
		}
		where, err := frontmatter.ParseWhere(whereStr)
		if err != nil {
			log.Fatal(err)
		}
		err = actions.SearchNotes(&vault, &note, &uri, &fuzzyFinder, useEditor, where)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	searchCmd.Flags().BoolP("editor", "e", false, "open in editor instead of Obsidian")
	searchCmd.Flags().StringP("where", "w", "", "filter by frontmatter key=value pairs (comma-separated)")
	rootCmd.AddCommand(searchCmd)
}
