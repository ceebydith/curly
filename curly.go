// Package curly provides a simple way to format strings using dynamic values.
// It allows users to define placeholders within a string and replace them with values
// from predefined maps or custom maps. This package is useful for generating dynamic
// content such as configuration files, messages, and paths.
package curly

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ErrResolve is returned when there are unresolved placeholders in the string
var ErrResolve = errors.New("resolve error")

// DefaultMaps contains predefined keys with static and dynamic values/functions
var DefaultMaps = map[string]any{
	"yyyy":   formatTime("2006"), // Full year (e.g., 2024)
	"yy":     formatTime("06"),   // Last two digits of the year (e.g., 24)
	"mm":     formatTime("01"),   // Month (01-12)
	"dd":     formatTime("02"),   // Day of the month (01-31)
	"hh":     formatTime("15"),   // Hour (00-23)
	"nn":     formatTime("04"),   // Minute (00-59)
	"ss":     formatTime("05"),   // Second (00-59)
	"appdir": getAppDir(),        // Directory of the executable
	"curdir": getCurDir(),        // Current working directory
}

// formatTime creates a function to format current time based on the provided layout
func formatTime(layout string) func(string) string {
	return func(key string) string { return time.Now().Format(layout) }
}

// getAppDir returns the directory of the executable file
func getAppDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "{appdir}" // Fallback value if there's an error
	}
	return filepath.Dir(exe)
}

// getCurDir returns the current working directory
func getCurDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "{curdir}" // Fallback value if there's an error
	}
	return dir
}

// resolved checks if all placeholders in the string have been resolved.
// Returns an error if there are unresolved placeholders.
func resolved(str string) error {
	// Remove escaped placeholders from consideration
	str = strings.ReplaceAll(str, "\\{", "")
	str = strings.ReplaceAll(str, "\\}", "")
	if strings.Contains(str, "{") || strings.Contains(str, "}") {
		return ErrResolve
	}
	return nil
}

// Format replaces placeholders in the input string with their corresponding values from the provided maps.
// If multiple maps are provided, the values from the later maps will override those from the earlier ones.
// Returns the formatted string and an error if unresolved placeholders remain.
func Format(str string, maps ...map[string]any) (string, error) {
	// Append the DefaultMaps to the provided maps
	maps = append(maps, DefaultMaps)
	for _, m := range maps {
		for key, val := range m {
			switch v := val.(type) {
			case string:
				str = strings.ReplaceAll(str, "{"+key+"}", v)
			case func(string) string:
				str = strings.ReplaceAll(str, "{"+key+"}", v(key))
			}
		}
	}
	// Check for unresolved placeholders
	if err := resolved(str); err != nil {
		return "", err
	}
	// Restore escaped placeholders
	str = strings.ReplaceAll(str, "\\{", "{")
	str = strings.ReplaceAll(str, "\\}", "}")
	return str, nil
}
