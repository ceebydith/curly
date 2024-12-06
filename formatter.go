package curly

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Formatter interface {
	Valid(identifier string) bool
	Value(identifier string) (any, error)
}

func NewMapFormatter(maps map[string]any) *MapFormatter {
	return &MapFormatter{
		maps: maps,
	}
}

func NewDatetimeFormatter() *DatetimeFormatter {
	return &DatetimeFormatter{}
}

func NewDirectoryFormatter() *DirectoryFormatter {
	return &DirectoryFormatter{}
}

type MapFormatter struct {
	maps map[string]any
}

func (f *MapFormatter) Valid(identifier string) bool {
	identifier = strings.ToLower(identifier)
	for key := range f.maps {
		if strings.ToLower(key) == identifier {
			return true
		}
	}
	return false
}

func (f *MapFormatter) Value(identifier string) (any, error) {
	identifier = strings.ToLower(identifier)
	for key, val := range f.maps {
		if strings.ToLower(key) == identifier {
			return val, nil
		}
	}
	return nil, fmt.Errorf("invalid identifier: \"%s\"", identifier)
}

type DatetimeFormatter struct{}

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

type DirectoryFormatter struct {
	once   sync.Once
	appdir string
	curdir string
	err    error
}

func (f *DirectoryFormatter) Valid(identifier string) bool {
	return identifier == "appdir" || identifier == "curdir"
}

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
