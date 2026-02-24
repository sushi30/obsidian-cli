package frontmatter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

const (
	Delimiter              = "---"
	NoFrontmatterError     = "note does not contain frontmatter"
	InvalidFrontmatterError = "frontmatter contains invalid YAML"
)

// Parse extracts and parses frontmatter from note content.
// Returns the frontmatter as a map, the body content, and any error.
func Parse(content string) (map[string]interface{}, string, error) {
	var fm map[string]interface{}
	rest, err := frontmatter.Parse(strings.NewReader(content), &fm)
	if err != nil {
		return nil, "", errors.New(InvalidFrontmatterError)
	}
	return fm, string(rest), nil
}

// Format converts a frontmatter map to a YAML string.
func Format(fm map[string]interface{}) (string, error) {
	if fm == nil || len(fm) == 0 {
		return "", nil
	}
	data, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HasFrontmatter checks if content starts with frontmatter delimiters.
func HasFrontmatter(content string) bool {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return false
	}
	return strings.TrimSpace(lines[0]) == Delimiter
}

// SetKey updates or adds a key in the frontmatter, returning the full updated content.
// If no frontmatter exists, it creates new frontmatter with the key.
func SetKey(content, key, value string) (string, error) {
	parsedValue := parseValue(value)

	if !HasFrontmatter(content) {
		// Create new frontmatter
		fm := map[string]interface{}{key: parsedValue}
		fmStr, err := yaml.Marshal(fm)
		if err != nil {
			return "", err
		}
		return Delimiter + "\n" + string(fmStr) + Delimiter + "\n" + content, nil
	}

	// Parse existing frontmatter
	fm, body, err := Parse(content)
	if err != nil {
		return "", err
	}

	if fm == nil {
		fm = make(map[string]interface{})
	}

	// Update the key
	fm[key] = parsedValue

	// Reconstruct content
	fmStr, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}

	return Delimiter + "\n" + string(fmStr) + Delimiter + "\n" + body, nil
}

// DeleteKey removes a key from the frontmatter, returning the full updated content.
func DeleteKey(content, key string) (string, error) {
	if !HasFrontmatter(content) {
		return "", errors.New(NoFrontmatterError)
	}

	fm, body, err := Parse(content)
	if err != nil {
		return "", err
	}

	if fm == nil {
		return "", errors.New(NoFrontmatterError)
	}

	// Delete the key
	delete(fm, key)

	// If no keys left, return just the body
	if len(fm) == 0 {
		return strings.TrimPrefix(body, "\n"), nil
	}

	// Reconstruct content
	fmStr, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}

	return Delimiter + "\n" + string(fmStr) + Delimiter + "\n" + body, nil
}

// parseValue attempts to parse the value into appropriate Go types.
// Supports: booleans, arrays (comma-separated in brackets), strings.
func parseValue(value string) interface{} {
	// Try boolean
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Try array (comma-separated in brackets)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := value[1 : len(value)-1]
		if inner == "" {
			return []string{}
		}
		parts := strings.Split(inner, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			result = append(result, strings.TrimSpace(p))
		}
		return result
	}

	// Default to string
	return value
}

// MatchesFilter checks if frontmatter matches the given filter criteria.
// For primitives (string, bool, number), uses equality.
// For lists, checks if the list contains the filter value.
func MatchesFilter(fm map[string]interface{}, filters map[string]string) bool {
	for key, filterValue := range filters {
		fmValue, exists := fm[key]
		if !exists {
			return false
		}

		// Check if the frontmatter value is a slice (list)
		switch v := fmValue.(type) {
		case []interface{}:
			// For lists, check if any element matches
			found := false
			for _, item := range v {
				if matchesPrimitive(item, filterValue) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case []string:
			// Handle string slices specifically
			found := false
			for _, item := range v {
				if matchesPrimitive(item, filterValue) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		default:
			// For primitives, use equality
			if !matchesPrimitive(fmValue, filterValue) {
				return false
			}
		}
	}
	return true
}

// matchesPrimitive checks if a primitive value matches the filter string.
func matchesPrimitive(value interface{}, filterValue string) bool {
	switch v := value.(type) {
	case string:
		return v == filterValue
	case bool:
		filterBool := filterValue == "true"
		return v == filterBool
	case int, int8, int16, int32, int64:
		return strings.TrimSpace(filterValue) == fmt.Sprintf("%d", v)
	case float32, float64:
		return strings.TrimSpace(filterValue) == fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v) == filterValue
	}
}

// ParseFilters converts a slice of "key=value" strings into a map.
func ParseFilters(filterStrings []string) (map[string]string, error) {
	filters := make(map[string]string)
	for _, filterStr := range filterStrings {
		parts := strings.SplitN(filterStr, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid filter format: " + filterStr + " (expected key=value)")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, errors.New("filter key cannot be empty")
		}
		filters[key] = value
	}
	return filters, nil
}
