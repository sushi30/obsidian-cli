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

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name           string
		frontmatter    map[string]interface{}
		filters        map[string]string
		expectedResult bool
	}{
		{
			name: "String equality match",
			frontmatter: map[string]interface{}{
				"status": "active",
			},
			filters: map[string]string{
				"status": "active",
			},
			expectedResult: true,
		},
		{
			name: "String equality no match",
			frontmatter: map[string]interface{}{
				"status": "inactive",
			},
			filters: map[string]string{
				"status": "active",
			},
			expectedResult: false,
		},
		{
			name: "Boolean equality match",
			frontmatter: map[string]interface{}{
				"published": true,
			},
			filters: map[string]string{
				"published": "true",
			},
			expectedResult: true,
		},
		{
			name: "Boolean equality no match",
			frontmatter: map[string]interface{}{
				"published": false,
			},
			filters: map[string]string{
				"published": "true",
			},
			expectedResult: false,
		},
		{
			name: "List contains match",
			frontmatter: map[string]interface{}{
				"tags": []interface{}{"go", "testing", "obsidian"},
			},
			filters: map[string]string{
				"tags": "testing",
			},
			expectedResult: true,
		},
		{
			name: "List contains no match",
			frontmatter: map[string]interface{}{
				"tags": []interface{}{"go", "testing", "obsidian"},
			},
			filters: map[string]string{
				"tags": "python",
			},
			expectedResult: false,
		},
		{
			name: "String list contains match",
			frontmatter: map[string]interface{}{
				"tags": []string{"go", "testing", "obsidian"},
			},
			filters: map[string]string{
				"tags": "go",
			},
			expectedResult: true,
		},
		{
			name: "Multiple filters all match",
			frontmatter: map[string]interface{}{
				"status": "active",
				"type":   "project",
			},
			filters: map[string]string{
				"status": "active",
				"type":   "project",
			},
			expectedResult: true,
		},
		{
			name: "Multiple filters one doesn't match",
			frontmatter: map[string]interface{}{
				"status": "active",
				"type":   "note",
			},
			filters: map[string]string{
				"status": "active",
				"type":   "project",
			},
			expectedResult: false,
		},
		{
			name: "Filter key doesn't exist",
			frontmatter: map[string]interface{}{
				"status": "active",
			},
			filters: map[string]string{
				"type": "project",
			},
			expectedResult: false,
		},
		{
			name: "Integer match",
			frontmatter: map[string]interface{}{
				"priority": 1,
			},
			filters: map[string]string{
				"priority": "1",
			},
			expectedResult: true,
		},
		{
			name: "Float match",
			frontmatter: map[string]interface{}{
				"version": 1.5,
			},
			filters: map[string]string{
				"version": "1.5",
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := frontmatter.MatchesFilter(tt.frontmatter, tt.filters)
			if result != tt.expectedResult {
				t.Errorf("MatchesFilter() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestParseFilters(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expected    map[string]string
		expectError bool
	}{
		{
			name:  "Single filter",
			input: []string{"status=active"},
			expected: map[string]string{
				"status": "active",
			},
			expectError: false,
		},
		{
			name:  "Multiple filters",
			input: []string{"status=active", "type=project"},
			expected: map[string]string{
				"status": "active",
				"type":   "project",
			},
			expectError: false,
		},
		{
			name:  "Filter with spaces",
			input: []string{"status = active"},
			expected: map[string]string{
				"status": "active",
			},
			expectError: false,
		},
		{
			name:        "Invalid filter format (no =)",
			input:       []string{"status"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid filter format (empty key)",
			input:       []string{"=value"},
			expected:    nil,
			expectError: true,
		},
		{
			name:  "Filter with = in value",
			input: []string{"formula=x=y+z"},
			expected: map[string]string{
				"formula": "x=y+z",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := frontmatter.ParseFilters(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseFilters() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ParseFilters() unexpected error: %v", err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ParseFilters() length = %v, want %v", len(result), len(tt.expected))
				}
				for key, value := range tt.expected {
					if result[key] != value {
						t.Errorf("ParseFilters()[%s] = %v, want %v", key, result[key], value)
					}
				}
			}
		})
	}
}

