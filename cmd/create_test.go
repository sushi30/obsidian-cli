package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveContent(t *testing.T) {
	t.Run("flag takes precedence over stdin", func(t *testing.T) {
		r := strings.NewReader("stdin content")
		result := resolveContent("flag content", r, func() bool { return false })
		assert.Equal(t, "flag content", result)
	})

	t.Run("reads stdin when piped", func(t *testing.T) {
		r := strings.NewReader("piped content\nwith newlines\n")
		result := resolveContent("", r, func() bool { return false })
		assert.Equal(t, "piped content\nwith newlines\n", result)
	})

	t.Run("ignores stdin when terminal", func(t *testing.T) {
		r := strings.NewReader("should be ignored")
		result := resolveContent("", r, func() bool { return true })
		assert.Equal(t, "", result)
	})

	t.Run("empty when no flag and terminal", func(t *testing.T) {
		r := strings.NewReader("")
		result := resolveContent("", r, func() bool { return true })
		assert.Equal(t, "", result)
	})
}
