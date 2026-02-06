package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var vaultName string
var sectionName string
var createIfNotExist bool
var OpenVaultCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Opens note in vault by note name (use @daily for daily note)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		uri := obsidian.Uri{}
		noteName, err := ResolveNoteName(&vault, args[0])
		if err != nil {
			log.Fatal(err)
		}
		params := actions.OpenParams{NoteName: noteName, Section: sectionName, CreateIfNotExist: createIfNotExist}
		err = actions.OpenNote(&vault, &uri, params)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	OpenVaultCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name (not required if default is set)")
	OpenVaultCmd.Flags().StringVarP(&sectionName, "section", "s", "", "heading text to open within the note (case-sensitive)")
	OpenVaultCmd.Flags().BoolVar(&createIfNotExist, "create-if-not-exist", false, "create an empty note if it does not exist")
	rootCmd.AddCommand(OpenVaultCmd)
}
