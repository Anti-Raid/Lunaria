// helpers.go
package lunaria

import (
	"regexp"
	"strconv"
	"strings"
)

// GetAttr retrieves an attribute value by name from a Node
func GetAttr(node Node, name string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// HasAttr checks if a node has a specific attribute
func HasAttr(node Node, name string) bool {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return true
		}
	}
	return false
}

// GetAttrWithDefault retrieves an attribute value with a default fallback
func GetAttrWithDefault(node Node, name, defaultValue string) string {
	if value := GetAttr(node, name); value != "" {
		return value
	}
	return defaultValue
}

// GetBoolAttr retrieves a boolean attribute value
func GetBoolAttr(node Node, name string) bool {
	value := GetAttr(node, name)
	return value == "true" || value == "1" || value == "yes"
}

// Interpolate replaces {{var}} patterns with Luau string concatenation
func Interpolate(text string) string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		varName := strings.TrimSpace(match[2 : len(match)-2])
		return `" .. tostring(` + varName + `) .. "`
	})
}

// ParseNumber safely converts a string to a number, defaulting to 0
func ParseNumber(s string) int {
	if num, err := strconv.Atoi(s); err == nil {
		return num
	}
	return 0
}

// ParseFloat safely converts a string to a float, defaulting to 0.0
func ParseFloat(s string) float64 {
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num
	}
	return 0.0
}

// EscapeString properly escapes a string for Luau
func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

// IsValidIdentifier checks if a string is a valid Luau identifier
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	// Must start with letter or underscore
	first := s[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}

	// Rest can be letters, digits, or underscores
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	// Check against Luau keywords
	keywords := []string{
		"and", "break", "do", "else", "elseif", "end", "false", "for",
		"function", "if", "in", "local", "nil", "not", "or", "repeat",
		"return", "then", "true", "until", "while",
	}

	for _, keyword := range keywords {
		if s == keyword {
			return false
		}
	}

	return true
}

// SplitParameters splits a parameter string into individual parameters
func SplitParameters(params string) []string {
	if params == "" {
		return []string{}
	}

	var result []string
	var current strings.Builder
	var inString bool
	var stringChar rune
	var depth int

	for _, r := range params {
		switch r {
		case '"', '\'':
			if !inString {
				inString = true
				stringChar = r
			} else if r == stringChar {
				inString = false
			}
			current.WriteRune(r)
		case '(', '[', '{':
			if !inString {
				depth++
			}
			current.WriteRune(r)
		case ')', ']', '}':
			if !inString {
				depth--
			}
			current.WriteRune(r)
		case ',':
			if !inString && depth == 0 {
				param := strings.TrimSpace(current.String())
				if param != "" {
					result = append(result, param)
				}
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add the last parameter
	param := strings.TrimSpace(current.String())
	if param != "" {
		result = append(result, param)
	}

	return result
}

// IndentLines adds indentation to each line of a multi-line string
func IndentLines(text string, indent string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	result := make([]string, len(lines))

	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			result[i] = indent + line
		} else {
			result[i] = ""
		}
	}

	return strings.Join(result, "\n")
}

// FormatComment formats a string as a Luau comment
func FormatComment(text string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(text), "\n")
	result := make([]string, len(lines))

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result[i] = "-- " + line
		} else {
			result[i] = "--"
		}
	}

	return strings.Join(result, "\n")
}

// GenerateVariableName generates a unique variable name with a prefix
func GenerateVariableName(prefix string, counter int) string {
	if prefix == "" {
		prefix = "var"
	}
	return prefix + strconv.Itoa(counter)
}

// IsStringLiteral checks if a string is a Luau string literal
func IsStringLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return false
	}

	return (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') ||
		(strings.HasPrefix(s, "[[") && strings.HasSuffix(s, "]]"))
}

// IsNumberLiteral checks if a string is a valid Luau number
func IsNumberLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	// Try to parse as int
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}

	return false
}

// WrapInQuotes wraps a string in quotes if it's not already a string literal
func WrapInQuotes(s string) string {
	if IsStringLiteral(s) || IsNumberLiteral(s) {
		return s
	}

	// Check if it looks like a variable or expression
	if IsValidIdentifier(s) || strings.Contains(s, "(") || strings.Contains(s, ".") {
		return s
	}

	return `"` + EscapeString(s) + `"`
}

// JoinWithCommas joins strings with commas, filtering out empty strings
func JoinWithCommas(strs []string) string {
	var filtered []string
	for _, s := range strs {
		if s = strings.TrimSpace(s); s != "" {
			filtered = append(filtered, s)
		}
	}
	return strings.Join(filtered, ", ")
}

// TrimAndClean trims whitespace and removes extra spaces
func TrimAndClean(s string) string {
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
