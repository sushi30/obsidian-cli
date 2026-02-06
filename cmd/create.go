package cmd

import (
	"io"
	"log"
	"os"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var shouldOverwrite bool
var content string
var createNoteCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Creates note in vault (use @daily for daily note). Reads from stdin if -c not provided.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		uri := obsidian.Uri{}
		noteName, err := ResolveNoteName(&vault, args[0])
		if err != nil {
			log.Fatal(err)
		}
		useEditor, err := cmd.Flags().GetBool("editor")
		if err != nil {
			log.Fatalf("Failed to parse --editor flag: %v", err)
		}

		// Read from stdin if data is being piped and -c is not supplied
		noteContent := content
		if content == "" && !term.IsTerminal(int(os.Stdin.Fd())) {
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Failed to read from stdin: %v", err)
			}
			noteContent = string(stdinBytes)
		}

		if noteContent == "" {
			log.Fatal("No content provided. Use -c flag or pipe from stdin:\n  obs create note -c \"content\"\n  echo \"content\" | obs create note")
		}

		params := actions.CreateParams{
			NoteName:        noteName,
			Content:         noteContent,
			ShouldAppend:    false,
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
	createNoteCmd.Flags().BoolVarP(&shouldOverwrite, "overwrite", "o", false, "overwrite existing note")
	createNoteCmd.Flags().BoolP("editor", "e", false, "open in editor instead of Obsidian (requires --open flag)")
	rootCmd.AddCommand(createNoteCmd)
}
