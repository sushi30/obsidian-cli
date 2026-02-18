package obsidian

import (
	"strings"
	"time"
)

const DailyReference = "@daily"

// ExpandDatePattern replaces date format tokens with actual date values.
// Supported tokens: YYYY, YY, MM, DD, MMM (short month), MMMM (full month)
func ExpandDatePattern(pattern string, t time.Time) string {
	replacer := strings.NewReplacer(
		"YYYY", t.Format("2006"),
		"YY", t.Format("06"),
		"MMMM", t.Format("January"),
		"MMM", t.Format("Jan"),
		"MM", t.Format("01"),
		"DD", t.Format("02"),
	)
	return replacer.Replace(pattern)
}

// IsDailyReference checks if the note name is the @daily special reference
func IsDailyReference(noteName string) bool {
	return noteName == DailyReference
}
