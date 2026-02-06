package cmd

import (
	"fmt"
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:     "append <note> <content>",
	Aliases: []string{"a"},
	Short:   "Append content to an existing note (use @daily for daily note)",
	Long: `Append text content to the end of an existing note.

Supports escape sequences like \n for newlines and \t for tabs.

Examples:
  obsidian-cli append "My Note" "New paragraph at the end"
  obsidian-cli append "Daily Note" "\n## New Section\nContent here"
  obsidian-cli append "Todo" -v work "\n- [ ] New task"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		noteName, err := ResolveNoteName(&vault, args[0])
		if err != nil {
			log.Fatal(err)
		}
		content := args[1]
		note := obsidian.Note{}

		params := actions.AppendParams{
			NoteName: noteName,
			Content:  content,
		}

		output, err := actions.AppendToNote(&vault, &note, params)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(output)
	},
}

func init() {
	appendCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	rootCmd.AddCommand(appendCmd)
}
