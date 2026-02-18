package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var appendCmd = &cobra.Command{
	Use:     "append <note> <content>",
	Aliases: []string{"a"},
	Short:   "Append content to an existing note (use @daily for daily note). Reads from stdin if content arg not provided.",
	Long: `Append text content to the end of an existing note.

Supports escape sequences like \n for newlines and \t for tabs.
Reads from stdin if content argument is not provided.

Examples:
  obsidian-cli append "My Note" "New paragraph at the end"
  obsidian-cli append "Daily Note" "\n## New Section\nContent here"
  obsidian-cli append "Todo" -v work "\n- [ ] New task"
  echo "piped content" | obsidian-cli append "My Note"`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		originalNoteName := args[0]
		noteName, err := ResolveNoteName(&vault, originalNoteName)
		if err != nil {
			log.Fatal(err)
		}

		var content string
		if len(args) >= 2 {
			content = args[1]
		} else if !term.IsTerminal(int(os.Stdin.Fd())) {
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Failed to read from stdin: %v", err)
			}
			content = string(stdinBytes)
		}

		if content == "" {
			log.Fatal("No content provided. Pass as argument or pipe from stdin:\n  obsidian-cli append \"note\" \"content\"\n  echo \"content\" | obsidian-cli append \"note\"")
		}

		note := obsidian.Note{}

		params := actions.AppendParams{
			NoteName: noteName,
			Content:  content,
		}

		output, err := actions.AppendToNote(&vault, &note, params)
		if err != nil {
			log.Fatal(WrapDailyNoteError(originalNoteName, err))
		}

		fmt.Println(output)
	},
}

func init() {
	appendCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	rootCmd.AddCommand(appendCmd)
}
