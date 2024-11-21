package curly_test

import (
	"testing"
	"time"

	"github.com/ceebydith/curly"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		input    string
		maps     []map[string]any
		expected string
		hasError bool
	}{
		{
			// Test default mappings without custom maps
			input: "Year: {yyyy}, Month: {mm}, Day: {dd}, Hour: {hh}, Minute: {nn}, Second: {ss}, App Directory: {appdir}, Current Directory: {curdir}",
			maps:  nil,
			expected: "Year: " + time.Now().Format("2006") + ", Month: " + time.Now().Format("01") +
				", Day: " + time.Now().Format("02") + ", Hour: " + time.Now().Format("15") +
				", Minute: " + time.Now().Format("04") + ", Second: " + time.Now().Format("05") +
				", App Directory: " + curly.DefaultMaps["appdir"].(string) + ", Current Directory: " + curly.DefaultMaps["curdir"].(string),
			hasError: false,
		},
		{
			// Test custom mappings
			input:    "Year: {yyyy}, Custom: {custom}",
			maps:     []map[string]any{{"custom": "Custom Value"}},
			expected: "Year: " + time.Now().Format("2006") + ", Custom: Custom Value",
			hasError: false,
		},
		{
			// Test unresolved placeholders
			input:    "Unresolved: {unresolved}",
			maps:     nil,
			expected: "",
			hasError: true,
		},
		{
			// Test escaped placeholders
			input:    "Escaped: \\{escaped\\}",
			maps:     nil,
			expected: "Escaped: {escaped}",
			hasError: false,
		},
		{
			// Test replacement of escaped placeholders
			input:    "Escaped and resolved: \\{escaped\\} and {yyyy}",
			maps:     nil,
			expected: "Escaped and resolved: {escaped} and " + time.Now().Format("2006"),
			hasError: false,
		},
	}

	for _, tt := range tests {
		result, err := curly.Format(tt.input, tt.maps...)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}
