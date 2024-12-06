package curly

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Formatter is an interface for formatting values based on identifiers.
type Formatter interface {
	Valid(identifier string) bool
	Value(identifier string) (any, error)
}

// NewMapFormatter creates a new MapFormatter with the provided map.
func NewMapFormatter(maps map[string]any) *MapFormatter {
	return &MapFormatter{
		maps: maps,
	}
}

// NewDatetimeFormatter creates a new DatetimeFormatter.
func NewDatetimeFormatter() *DatetimeFormatter {
	return &DatetimeFormatter{}
}

// NewDirectoryFormatter creates a new DirectoryFormatter.
func NewDirectoryFormatter() *DirectoryFormatter {
	return &DirectoryFormatter{}
}

// MapFormatter formats values based on a map of identifiers.
type MapFormatter struct {
	maps map[string]any
}

// Valid checks if the identifier is valid in the map.
func (f *MapFormatter) Valid(identifier string) bool {
	identifier = strings.ToLower(identifier)
	for key := range f.maps {
		if strings.ToLower(key) == identifier {
			return true
		}
	}
	return false
}

// Value returns the value associated with the identifier in the map.
func (f *MapFormatter) Value(identifier string) (any, error) {
	identifier = strings.ToLower(identifier)
	for key, val := range f.maps {
		if strings.ToLower(key) == identifier {
			return val, nil
		}
	}
	return nil, fmt.Errorf("invalid identifier: \"%s\"", identifier)
}

// DatetimeFormatter formats values based on date and time.
type DatetimeFormatter struct{}

// Valid checks if the identifier is a valid date or time identifier.
func (f *DatetimeFormatter) Valid(identifier string) bool {
	identifier = strings.ToLower(identifier)
	identifiers := []string{"yyyy", "yy", "mm", "dd", "hh", "nn", "ss"}
	for _, iden := range identifiers {
		if identifier == iden {
			return true
		}
	}
	return false
}

// Value returns the current date or time based on the identifier.
func (f *DatetimeFormatter) Value(identifier string) (any, error) {
	identifier = strings.ToLower(identifier)
	switch identifier {
	case "yyyy":
		return time.Now().Format("2006"), nil
	case "yy":
		return time.Now().Format("06"), nil
	case "mm":
		return time.Now().Format("01"), nil
	case "dd":
		return time.Now().Format("02"), nil
	case "hh":
		return time.Now().Format("15"), nil
	case "nn":
		return time.Now().Format("04"), nil
	case "ss":
		return time.Now().Format("05"), nil
	}
	return nil, fmt.Errorf("invalid identifier: \"%s\"", identifier)
}

// DirectoryFormatter formats values based on the directory paths.
type DirectoryFormatter struct {
	once   sync.Once
	appdir string
	curdir string
	err    error
}

// Valid checks if the identifier is a valid directory identifier.
func (f *DirectoryFormatter) Valid(identifier string) bool {
	return identifier == "appdir" || identifier == "curdir"
}

// Value returns the directory path based on the identifier.
func (f *DirectoryFormatter) Value(identifier string) (any, error) {
	f.once.Do(func() {
		exe, err := os.Executable()
		if err != nil {
			f.err = err
			return
		}
		f.appdir = filepath.Dir(exe)
		f.curdir, f.err = os.Getwd()
	})

	if f.err != nil {
		return nil, f.err
	} else if identifier == "appdir" {
		return f.appdir, nil
	} else if identifier == "curdir" {
		return f.curdir, nil
	}
	return nil, fmt.Errorf("invalid identifier: \"%s\"", identifier)
}
