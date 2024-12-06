package curly

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parser is an interface for parsing values based on identifiers.
type Parser interface {
	Valid(identifier string) bool
	Expressions() []string
	Modify(value string, index int) any
}

// NewMsisdnParser creates a new MsisdnParser with the provided country code, format, and identifiers.
func NewMsisdnParser(country uint, format int, identifiers ...string) *MsisdnParser {
	return &MsisdnParser{
		identifiers: identifiers,
		country:     country,
		format:      format,
	}
}

// NewNumberParser creates a new NumberParser with the provided identifiers.
func NewNumberParser(identifiers ...string) *NumberParser {
	return &NumberParser{
		identifiers: identifiers,
	}
}

// NewStringParser creates a new StringParser. If trimspace is true, the parser will trim spaces.
func NewStringParser(trimspace ...bool) *StringParser {
	return &StringParser{
		trimspace: len(trimspace) == 0 || trimspace[0],
	}
}

// MsisdnParser parses mobile subscriber numbers (MSISDN) based on country code and format.
type MsisdnParser struct {
	identifiers []string
	country     uint
	format      int
}

// Valid checks if the identifier is valid for MSISDN parsing.
func (p *MsisdnParser) Valid(identifier string) bool {
	return ValidIdentifier(identifier, p.identifiers)
}

// Expressions returns the regex expressions for parsing MSISDN based on the country code.
func (p *MsisdnParser) Expressions() []string {
	if p.country == 0 {
		return []string{`0[1-9][0-9]+`}
	}
	return []string{fmt.Sprintf(`(0|\+?%d)[1-9][0-9]+`, p.country)}
}

// Modify formats the MSISDN based on the provided format.
func (p *MsisdnParser) Modify(value string, index int) any {
	if p.country == 0 || index != 0 {
		return value
	}
	var replace string
	if p.format < 0 {
		replace = "0"
	} else if p.format > 0 {
		replace = fmt.Sprintf("+%d", p.country)
	} else {
		replace = fmt.Sprintf("%d", p.country)
	}
	reg := regexp.MustCompile(fmt.Sprintf(`^(0|\+?%d)`, p.country))
	return reg.ReplaceAllString(value, replace)
}

// NumberParser parses numerical values based on specific formats.
type NumberParser struct {
	identifiers []string
}

// Valid checks if the identifier is valid for number parsing.
func (p *NumberParser) Valid(identifier string) bool {
	return ValidIdentifier(identifier, p.identifiers)
}

// Expressions returns the regex expressions for parsing numbers in various formats.
func (p *NumberParser) Expressions() []string {
	return []string{
		`\-?([0-9]{1,3}(\.[0-9]{3})*|([0-9]+))(\,[0-9]{1,2})?`,
		`\-?([0-9]{1,3}(\,[0-9]{3})*|([0-9]+))(\.[0-9]{1,2})?`,
	}
}

// Modify formats the parsed number based on the provided format.
func (p *NumberParser) Modify(value string, index int) any {
	format := [][2]string{
		{".", ","},
		{",", "."},
	}
	value = strings.TrimSpace(value)
	if 0 <= index && index < len(format) {
		value = strings.ReplaceAll(value, format[index][0], "")
		value = strings.ReplaceAll(value, format[index][1], ".")
	}
	if strings.Contains(value, ".") {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v
		}
	} else {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return value
}

// StringParser parses string values, optionally trimming spaces.
type StringParser struct {
	trimspace bool
}

// Valid always returns true for string parsing, as all strings are valid.
func (p *StringParser) Valid(identifier string) bool {
	return true
}

// Expressions returns the regex expressions for parsing strings.
func (p *StringParser) Expressions() []string {
	return []string{
		`.+`,
	}
}

// Modify formats the parsed string value, optionally trimming spaces.
func (p *StringParser) Modify(value string, index int) any {
	if p.trimspace {
		return strings.TrimSpace(value)
	}
	return value
}
