package cmd

import (
	"io"
	"log"
	"os"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

// resolveContent returns content from the flag if set, otherwise reads from
// stdin when it's piped (not a terminal). isTerminal reports whether stdin
// is an interactive terminal.
func resolveContent(flagContent string, r io.Reader, isTerminal func() bool) string {
	if flagContent != "" {
		return flagContent
	}
	if !isTerminal() {
		b, err := io.ReadAll(r)
		if err != nil {
			return ""
		}
		return string(b)
	}
	return ""
}

func stdinIsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return true
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

var shouldAppend bool
var shouldOverwrite bool
var content string
var createNoteCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Creates note in vault",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		uri := obsidian.Uri{}
		noteName := args[0]
		useEditor, err := cmd.Flags().GetBool("editor")
		if err != nil {
			log.Fatalf("Failed to parse --editor flag: %v", err)
		}
		resolved := resolveContent(content, os.Stdin, stdinIsTerminal)
		params := actions.CreateParams{
			NoteName:        noteName,
			Content:         resolved,
			ShouldAppend:    shouldAppend,
			ShouldOverwrite: shouldOverwrite,
			ShouldOpen:      shouldOpen,
			UseEditor:       useEditor,
		}
		err = actions.CreateNote(&vault, &uri, params)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createNoteCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	createNoteCmd.Flags().BoolVarP(&shouldOpen, "open", "", false, "open created note")
	createNoteCmd.Flags().StringVarP(&content, "content", "c", "", "text to add to note")
	createNoteCmd.Flags().BoolVarP(&shouldAppend, "append", "a", false, "append to note")
	createNoteCmd.Flags().BoolVarP(&shouldOverwrite, "overwrite", "o", false, "overwrite note")
	createNoteCmd.Flags().BoolP("editor", "e", false, "open in editor instead of Obsidian (requires --open flag)")
	createNoteCmd.MarkFlagsMutuallyExclusive("append", "overwrite")
	rootCmd.AddCommand(createNoteCmd)
}
