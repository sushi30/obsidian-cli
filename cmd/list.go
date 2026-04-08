package cmd

import (
	"fmt"
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List files and folders in vault",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		}

		whereStr, err := cmd.Flags().GetString("where")
		if err != nil {
			log.Fatalf("failed to retrieve 'where' flag: %v", err)
		}
		where, err := frontmatter.ParseWhere(whereStr)
		if err != nil {
			log.Fatal(err)
		}

		vault := obsidian.Vault{Name: vaultName}
		entries, err := actions.ListEntries(&vault, actions.ListParams{
			Path:  targetPath,
			Where: where,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			fmt.Printf("• %s\n", entry)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	listCmd.Flags().StringP("where", "w", "", "filter by frontmatter key=value pairs (comma-separated)")
	rootCmd.AddCommand(listCmd)
}
