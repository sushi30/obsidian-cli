package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var DailyCmd = &cobra.Command{
	Use:     "daily",
	Aliases: []string{"d"},
	Short:   "Creates or opens daily note in vault",
	Long: `Opens today's daily note. Equivalent to 'obsidian open @daily'.

If daily note pattern is configured (via set-daily-pattern), uses that path.
Otherwise falls back to Obsidian's native daily note handler.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		uri := obsidian.Uri{}

		noteName, err := vault.ResolveDailyNote()
		if err != nil {
			// Fallback to obsidian://daily when pattern not configured
			err = actions.DailyNote(&vault, &uri)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		params := actions.OpenParams{NoteName: noteName}
		err = actions.OpenNote(&vault, &uri, params)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	DailyCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name (not required if default is set)")
	rootCmd.AddCommand(DailyCmd)
}
