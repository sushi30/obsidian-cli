package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <note> <old-string> <new-string>",
	Short: "Replace text in a note (use @daily for daily note)",
	Long: `Replace exact string matches in a note. By default, only replaces the first occurrence.
Use --all flag to replace all occurrences.

Examples:
  obsidian-cli edit "My Note" "old text" "new text"
  obsidian-cli edit "@daily" "TODO" "DONE" --all
  obsidian-cli edit "Project" "phase 1" "phase 2" -v work`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		note := obsidian.Note{}
		originalNoteName := args[0]
		noteName, err := ResolveNoteName(&vault, originalNoteName)
		if err != nil {
			log.Fatal(err)
		}

		replaceAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatalf("Failed to parse --all flag: %v", err)
		}

		params := actions.EditParams{
			NoteName:   noteName,
			OldString:  args[1],
			NewString:  args[2],
			ReplaceAll: replaceAll,
		}

		output, err := actions.EditNote(&vault, &note, params)
		if err != nil {
			log.Fatal(WrapDailyNoteError(originalNoteName, err))
		}

		log.Println(output)
	},
}

func init() {
	editCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	editCmd.Flags().BoolP("all", "a", false, "replace all occurrences of old string")
	rootCmd.AddCommand(editCmd)
}
