package obsidian

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpandDatePattern(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		pattern  string
		expected string
	}{
		{"simple date", "YYYY-MM-DD", "2024-03-15"},
		{"nested folders", "YYYY/MM/YYYY-MM-DD", "2024/03/2024-03-15"},
		{"with prefix", "daily/YYYY-MM-DD", "daily/2024-03-15"},
		{"short year", "YY-MM-DD", "24-03-15"},
		{"month names short", "YYYY/MMM/DD", "2024/Mar/15"},
		{"month names full", "MMMM DD, YYYY", "March 15, 2024"},
		{"no tokens", "daily/notes", "daily/notes"},
		{"empty pattern", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandDatePattern(tt.pattern, fixedTime)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDailyReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact match", "@daily", true},
		{"without at", "daily", false},
		{"wrong case", "@Daily", false},
		{"regular note", "my-note", false},
		{"empty", "", false},
		{"similar prefix", "@daily-note", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDailyReference(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
