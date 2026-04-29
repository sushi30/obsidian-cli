package obsidian_test

import (
	"errors"
	"net/url"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUriConstruct(t *testing.T) {
	baseUri := "base-uri"
	var tests = []struct {
		testName string
		in       map[string]string
		want     map[string]string
	}{
		{"Empty map", map[string]string{}, nil},
		{"One key", map[string]string{"key": "value"}, map[string]string{"key": "value"}},
		{"Two keys", map[string]string{"key1": "value1", "key2": "value2"}, map[string]string{"key1": "value1", "key2": "value2"}},
		{"Empty value", map[string]string{"key": ""}, nil},
		{"Mix of empty and non-empty values", map[string]string{"key1": "value1", "key2": ""}, map[string]string{"key1": "value1"}},
		{"Value with equals sign", map[string]string{"content": "x = 1"}, map[string]string{"content": "x = 1"}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Act
			uriManager := obsidian.Uri{}
			got := uriManager.Construct(baseUri, test.in)
			// Assert
			if test.want == nil {
				assert.Equal(t, baseUri, got)
			} else {
				parts := strings.SplitN(got, "?", 2)
				assert.Equal(t, baseUri, parts[0])

				if len(parts) > 1 {
					parsedParams, err := url.ParseQuery(parts[1])
					assert.NoError(t, err)
					assert.Equal(t, len(test.want), len(parsedParams), "unexpected number of parameters")
					for key, expectedValue := range test.want {
						assert.Equal(t, expectedValue, parsedParams.Get(key))
					}
				}
			}
		})
	}
}

func TestUriConstructEncodesQueryUnsafeChars(t *testing.T) {
	uriManager := obsidian.Uri{}

	tests := []struct {
		name  string
		value string
		must  string // substring that must NOT appear unescaped in the query string
	}{
		{"equals sign", "x = 1", "x = 1"},
		{"ampersand", "a & b", "a & b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uriManager.Construct("obsidian://new", map[string]string{"content": tt.value})
			query := strings.SplitN(got, "?", 2)[1]
			assert.NotContains(t, query, tt.must, "query-unsafe chars must be percent-encoded")
			// Round-trip: ParseQuery must recover the original value
			parsed, err := url.ParseQuery(query)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, parsed.Get("content"))
		})
	}
}

func TestUriExecute(t *testing.T) {
	// Temporarily override the Run function
	originalOpenerFunc := obsidian.Run
	defer func() { obsidian.Run = originalOpenerFunc }()

	t.Run("Valid URI", func(t *testing.T) {
		obsidian.Run = func(uri string) error {
			return nil
		}
		// Arrange
		uriManager := obsidian.Uri{}
		// Act
		err := uriManager.Execute("https://example.com")
		// Assert
		assert.Equal(t, nil, err)
	})

	t.Run("Invalid URI", func(t *testing.T) {
		obsidian.Run = func(uri string) error {
			return errors.New("mock error")
		}
		// Arrange
		uriManager := obsidian.Uri{}
		// Act
		err := uriManager.Execute("foo")
		// Assert
		assert.Equal(t, obsidian.ExecuteUriError, err.Error())
	})

}
