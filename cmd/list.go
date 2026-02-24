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
	Long: `List files and folders in vault.

Examples:
  # List all notes
  obsidian-cli list

  # List notes with status=active in frontmatter
  obsidian-cli list --meta status=active

  # List notes with multiple filters
  obsidian-cli list --meta status=active --meta type=project`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		}

		fullPath, _ := cmd.Flags().GetBool("full-path")
		metadataFlags, _ := cmd.Flags().GetStringSlice("meta")

		var metadataFilters map[string]string
		var err error
		if len(metadataFlags) > 0 {
			metadataFilters, err = frontmatter.ParseFilters(metadataFlags)
			if err != nil {
				log.Fatal(err)
			}
		}

		vault := obsidian.Vault{Name: vaultName}
		entries, err := actions.ListEntries(&vault, actions.ListParams{
			Path:            targetPath,
			FullPath:        fullPath,
			MetadataFilters: metadataFilters,
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
	listCmd.Flags().Bool("full-path", false, "display full vault path for each entry")
	listCmd.Flags().StringSliceP("meta", "m", []string{}, "filter by frontmatter metadata (key=value)")
	rootCmd.AddCommand(listCmd)
}
