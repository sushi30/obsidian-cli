package cmd

import (
	"fmt"
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
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

		fullPath, _ := cmd.Flags().GetBool("full-path")

		vault := obsidian.Vault{Name: vaultName}
		entries, err := actions.ListEntries(&vault, actions.ListParams{Path: targetPath, FullPath: fullPath})
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
	listCmd.Flags().Bool("full-path", false, "display full vault path for each entry")
	rootCmd.AddCommand(listCmd)
}
