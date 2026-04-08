package frontmatter_test

import (
	"strings"
	"testing"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("Parse valid frontmatter", func(t *testing.T) {
		content := "---\ntitle: Test\ntags:\n  - a\n  - b\n---\nBody content"
		fm, body, err := frontmatter.Parse(content)
		assert.NoError(t, err)
		assert.Equal(t, "Test", fm["title"])
		assert.Equal(t, "Body content", body)
	})

	t.Run("Parse empty frontmatter", func(t *testing.T) {
		content := "---\n---\nBody content"
		fm, body, err := frontmatter.Parse(content)
		assert.NoError(t, err)
		assert.Empty(t, fm)
		assert.Equal(t, "Body content", body)
	})

	t.Run("No frontmatter returns empty map", func(t *testing.T) {
		content := "Just body content"
		fm, body, err := frontmatter.Parse(content)
		assert.NoError(t, err)
		assert.Empty(t, fm)
		assert.Equal(t, "Just body content", body)
	})

	t.Run("Invalid YAML returns error", func(t *testing.T) {
		content := "---\ninvalid: [unclosed\n---\nBody"
		_, _, err := frontmatter.Parse(content)
		assert.Error(t, err)
	})
}

func TestHasFrontmatter(t *testing.T) {
	t.Run("Has frontmatter", func(t *testing.T) {
		content := "---\ntitle: Test\n---\nBody"
		assert.True(t, frontmatter.HasFrontmatter(content))
	})

	t.Run("No frontmatter", func(t *testing.T) {
		content := "Just body content"
		assert.False(t, frontmatter.HasFrontmatter(content))
	})

	t.Run("Empty content", func(t *testing.T) {
		assert.False(t, frontmatter.HasFrontmatter(""))
	})
}

func TestFormat(t *testing.T) {
	t.Run("Format valid map", func(t *testing.T) {
		fm := map[string]interface{}{
			"title": "Test",
		}
		result, err := frontmatter.Format(fm)
		assert.NoError(t, err)
		assert.Contains(t, result, "title: Test")
	})

	t.Run("Format empty map", func(t *testing.T) {
		result, err := frontmatter.Format(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Format nil map", func(t *testing.T) {
		result, err := frontmatter.Format(nil)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestSetKey(t *testing.T) {
	t.Run("Add key to existing frontmatter", func(t *testing.T) {
		content := "---\ntitle: Test\n---\nBody"
		result, err := frontmatter.SetKey(content, "author", "John")
		assert.NoError(t, err)
		assert.Contains(t, result, "author: John")
		assert.Contains(t, result, "title: Test")
		assert.Contains(t, result, "Body")
	})

	t.Run("Update existing key", func(t *testing.T) {
		content := "---\ntitle: Old\n---\nBody"
		result, err := frontmatter.SetKey(content, "title", "New")
		assert.NoError(t, err)
		assert.Contains(t, result, "title: New")
		assert.NotContains(t, result, "title: Old")
	})

	t.Run("Create frontmatter when none exists", func(t *testing.T) {
		content := "Just body content"
		result, err := frontmatter.SetKey(content, "title", "New")
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(result, "---\n"))
		assert.Contains(t, result, "title: New")
		assert.Contains(t, result, "Just body content")
	})

	t.Run("Parse boolean value true", func(t *testing.T) {
		content := "---\n---\nBody"
		result, err := frontmatter.SetKey(content, "draft", "true")
		assert.NoError(t, err)
		assert.Contains(t, result, "draft: true")
	})

	t.Run("Parse boolean value false", func(t *testing.T) {
		content := "---\n---\nBody"
		result, err := frontmatter.SetKey(content, "published", "false")
		assert.NoError(t, err)
		assert.Contains(t, result, "published: false")
	})

	t.Run("Parse array value", func(t *testing.T) {
		content := "---\n---\nBody"
		result, err := frontmatter.SetKey(content, "tags", "[one, two, three]")
		assert.NoError(t, err)
		assert.Contains(t, result, "tags:")
		assert.Contains(t, result, "- one")
		assert.Contains(t, result, "- two")
		assert.Contains(t, result, "- three")
	})

	t.Run("Parse empty array value", func(t *testing.T) {
		content := "---\n---\nBody"
		result, err := frontmatter.SetKey(content, "tags", "[]")
		assert.NoError(t, err)
		assert.Contains(t, result, "tags: []")
	})
}

func TestParseWhere(t *testing.T) {
	t.Run("Empty string returns empty map", func(t *testing.T) {
		result, err := frontmatter.ParseWhere("")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Single pair", func(t *testing.T) {
		result, err := frontmatter.ParseWhere("status=done")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"status": "done"}, result)
	})

	t.Run("Multiple pairs", func(t *testing.T) {
		result, err := frontmatter.ParseWhere("status=done,tags=work")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"status": "done", "tags": "work"}, result)
	})

	t.Run("Missing equals returns error", func(t *testing.T) {
		_, err := frontmatter.ParseWhere("status")
		assert.Error(t, err)
	})

	t.Run("Missing equals in second pair returns error", func(t *testing.T) {
		_, err := frontmatter.ParseWhere("status=done,bad")
		assert.Error(t, err)
	})
}

func TestMatchesFilters(t *testing.T) {
	t.Run("String exact match", func(t *testing.T) {
		fm := map[string]interface{}{"status": "done"}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"status": "done"}))
	})

	t.Run("String mismatch", func(t *testing.T) {
		fm := map[string]interface{}{"status": "draft"}
		assert.False(t, frontmatter.MatchesFilters(fm, map[string]string{"status": "done"}))
	})

	t.Run("Array contains match", func(t *testing.T) {
		fm := map[string]interface{}{"tags": []interface{}{"work", "personal"}}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"tags": "work"}))
	})

	t.Run("Array miss", func(t *testing.T) {
		fm := map[string]interface{}{"tags": []interface{}{"work", "personal"}}
		assert.False(t, frontmatter.MatchesFilters(fm, map[string]string{"tags": "finance"}))
	})

	t.Run("Bool true", func(t *testing.T) {
		fm := map[string]interface{}{"draft": true}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"draft": "true"}))
	})

	t.Run("Bool false", func(t *testing.T) {
		fm := map[string]interface{}{"draft": false}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"draft": "false"}))
	})

	t.Run("Int match", func(t *testing.T) {
		fm := map[string]interface{}{"priority": 1}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"priority": "1"}))
	})

	t.Run("Float match", func(t *testing.T) {
		fm := map[string]interface{}{"score": 3.5}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"score": "3.5"}))
	})

	t.Run("Missing key returns false", func(t *testing.T) {
		fm := map[string]interface{}{"status": "done"}
		assert.False(t, frontmatter.MatchesFilters(fm, map[string]string{"missing": "value"}))
	})

	t.Run("Multiple filters AND logic", func(t *testing.T) {
		fm := map[string]interface{}{"status": "done", "tags": []interface{}{"work"}}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{"status": "done", "tags": "work"}))
		assert.False(t, frontmatter.MatchesFilters(fm, map[string]string{"status": "done", "tags": "personal"}))
	})

	t.Run("Nil frontmatter returns false", func(t *testing.T) {
		assert.False(t, frontmatter.MatchesFilters(nil, map[string]string{"status": "done"}))
	})

	t.Run("Empty filters returns true", func(t *testing.T) {
		fm := map[string]interface{}{"status": "done"}
		assert.True(t, frontmatter.MatchesFilters(fm, map[string]string{}))
	})

	t.Run("Empty filters with nil frontmatter returns true", func(t *testing.T) {
		assert.True(t, frontmatter.MatchesFilters(nil, map[string]string{}))
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("Delete existing key", func(t *testing.T) {
		content := "---\ntitle: Test\nauthor: John\n---\nBody"
		result, err := frontmatter.DeleteKey(content, "author")
		assert.NoError(t, err)
		assert.Contains(t, result, "title: Test")
		assert.NotContains(t, result, "author")
		assert.Contains(t, result, "Body")
	})

	t.Run("Delete last key removes frontmatter", func(t *testing.T) {
		content := "---\ntitle: Test\n---\nBody"
		result, err := frontmatter.DeleteKey(content, "title")
		assert.NoError(t, err)
		assert.False(t, strings.HasPrefix(result, "---"))
		assert.Contains(t, result, "Body")
	})

	t.Run("Delete from no frontmatter returns error", func(t *testing.T) {
		content := "Just body content"
		_, err := frontmatter.DeleteKey(content, "title")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain frontmatter")
	})

	t.Run("Delete non-existent key succeeds", func(t *testing.T) {
		content := "---\ntitle: Test\n---\nBody"
		result, err := frontmatter.DeleteKey(content, "nonexistent")
		assert.NoError(t, err)
		assert.Contains(t, result, "title: Test")
	})
}
