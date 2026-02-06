package cmd

import (
	"fmt"
	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"log"

	"github.com/spf13/cobra"
)

var shouldRenderMarkdown bool
var includeMentions bool

var printCmd = &cobra.Command{
	Use:     "print",
	Aliases: []string{"p"},
	Short:   "Print contents of note",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		originalNoteName := args[0]
		noteName, err := ResolveNoteName(&vault, originalNoteName)
		if err != nil {
			log.Fatal(err)
		}
		note := obsidian.Note{}
		params := actions.PrintParams{
			NoteName:        noteName,
			IncludeMentions: includeMentions,
		}
		contents, err := actions.PrintNote(&vault, &note, params)
		if err != nil {
			log.Fatal(WrapDailyNoteError(originalNoteName, err))
		}
		fmt.Println(contents)
	},
}

func init() {
	printCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	printCmd.Flags().BoolVarP(&includeMentions, "mentions", "m", false, "include linked mentions at the end")
	rootCmd.AddCommand(printCmd)
}
