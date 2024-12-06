package curly

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser interface {
	Valid(identifier string) bool
	Expressions() []string
	Modify(value string, index int) any
}

func NewMsisdnParser(country uint, format int, identifiers ...string) *MsisdnParser {
	return &MsisdnParser{
		identifiers: identifiers,
		country:     country,
		format:      format,
	}
}

func NewNumberParser(identifiers ...string) *NumberParser {
	return &NumberParser{
		identifiers: identifiers,
	}
}

func NewStringParser(trimspace ...bool) *StringParser {
	return &StringParser{
		trimspace: len(trimspace) == 0 || trimspace[0],
	}
}

type MsisdnParser struct {
	identifiers []string
	country     uint
	format      int
}

func (p *MsisdnParser) Valid(identifier string) bool {
	return ValidIdentifier(identifier, p.identifiers)
}

func (p *MsisdnParser) Expressions() []string {
	if p.country == 0 {
		return []string{`0[1-9][0-9]+`}
	}
	return []string{fmt.Sprintf(`(0|\+?%d)[1-9][0-9]+`, p.country)}
}

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

type NumberParser struct {
	identifiers []string
}

func (p *NumberParser) Valid(identifier string) bool {
	return ValidIdentifier(identifier, p.identifiers)
}

func (p *NumberParser) Expressions() []string {
	return []string{
		`\-?([0-9]{1,3}(\.[0-9]{3})*|([0-9]+))(\,[0-9]{1,2})?`,
		`\-?([0-9]{1,3}(\,[0-9]{3})*|([0-9]+))(\.[0-9]{1,2})?`,
	}
}

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

type StringParser struct {
	trimspace bool
}

func (p *StringParser) Valid(identifier string) bool {
	return true
}

func (p *StringParser) Expressions() []string {
	return []string{
		`.+`,
	}
}

func (p *StringParser) Modify(value string, index int) any {
	if p.trimspace {
		return strings.TrimSpace(value)
	}
	return value
}
