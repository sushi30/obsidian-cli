package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var setDailyPatternCmd = &cobra.Command{
	Use:     "set-daily-pattern <pattern>",
	Aliases: []string{"sdp"},
	Short:   "Sets the daily note path pattern",
	Long: `Sets the pattern for daily note paths. Use date format tokens:
  YYYY - 4-digit year (2024)
  YY   - 2-digit year (24)
  MM   - 2-digit month (01-12)
  MMM  - Short month name (Jan)
  MMMM - Full month name (January)
  DD   - 2-digit day (01-31)

Examples:
  obsidian set-daily-pattern "daily/YYYY-MM-DD"
  obsidian set-daily-pattern "YYYY/MM/YYYY-MM-DD"
  obsidian set-daily-pattern "journal/YYYY/MMM/DD"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]
		v := obsidian.Vault{}
		if err := v.SetDailyNotePattern(pattern); err != nil {
			log.Fatal(err)
		}

		example := obsidian.ExpandDatePattern(pattern, time.Now())
		fmt.Printf("Daily note pattern set to: %s\n", pattern)
		fmt.Printf("Today's daily note would be: %s\n", example)
	},
}

func init() {
	rootCmd.AddCommand(setDailyPatternCmd)
}
