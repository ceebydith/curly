package curly

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	mu               sync.RWMutex
	defaultModifiers []Modifier
)

// DefaultModifier manages the default list of modifiers.
func DefaultModifier(modifiers ...Modifier) []Modifier {
	if len(modifiers) != 0 {
		mu.Lock()
		defaultModifiers = modifiers
		mu.Unlock()
	}
	mu.RLock()
	defer mu.RUnlock()
	return defaultModifiers
}

// Modifier is an interface that defines methods for validating and modifying values with specific modifiers.
type Modifier interface {
	Valid(modifier string) bool
	Modify(value string, modifier string, onfailed ...func(value any, modifier string) (any, error)) (any, error)
}

// NewNumberModifier creates a new instance of NumberModifier.
func NewNumberModifier() *NumberModifier {
	return &NumberModifier{}
}

// NewStringModifier creates a new instance of StringModifier.
func NewStringModifier() *StringModifier {
	return &StringModifier{}
}

// NewFormatModifier creates a new instance of FormatModifier.
func NewFormatModifier() *FormatModifier {
	return &FormatModifier{}
}

// NumberModifier implements Modifier for numerical expressions.
type NumberModifier struct{}

// Valid checks if the modifier is a valid numerical expression.
func (m *NumberModifier) Valid(modifier string) bool {
	reg := regexp.MustCompile(`^\s*[\*/\+\-][\s\.\*/\+\-\(\)0-9]*?[0-9\)]\s*$`)
	return reg.MatchString(modifier) && charCount('(', modifier) == charCount(')', modifier)
}

// Modify applies the numerical modifier to the given value.
func (m *NumberModifier) Modify(value string, modifier string, onfailed ...func(value any, modifier string) (any, error)) (any, error) {
	modifier = value + modifier
	syntax := "(" + modifier + ")"
	if charCount('(', syntax) != charCount(')', syntax) {
		return nil, fmt.Errorf("invalid expression: \"%s\"", modifier)
	}
	regParenthesis := regexp.MustCompile(`\(\s*((\-\s*)?[0-9]+(\.[0-9]+)?)?\s*\)`)
	regOperator := []*regexp.Regexp{
		regexp.MustCompile(`\s*([\*/\+\-\(])\s*((\-\s*)?[0-9]+(\.[0-9]+)?)\s*([\*/])\s*((\-\s*)?[0-9]+(\.[0-9]+)?)`),
		regexp.MustCompile(`\s*([\*/\+\-\(])\s*((\-\s*)?[0-9]+(\.[0-9]+)?)\s*([\+\-])\s*((\-\s*)?[0-9]+(\.[0-9]+)?)`),
	}

	for {
		found := false
		matches := regParenthesis.FindAllStringSubmatch(syntax, -1)
		for _, match := range matches {
			if match[1] == "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", match[0])
			}
			syntax = strings.Replace(syntax, match[0], match[1], 1)
			found = true
		}
		for _, reg := range regOperator {
			for {
				matches := reg.FindAllStringSubmatch(syntax, -1)
				if len(matches) == 0 {
					break
				}
				for _, match := range matches {
					val, err := calculate(match[2], match[5], match[6])
					if err != nil {
						return nil, err
					}
					if math.IsInf(val, 0) {
						return nil, fmt.Errorf("division by zero: \"%s\"", modifier)
					}
					if strings.Contains(match[2]+match[6], ".") {
						syntax = strings.Replace(syntax, match[0], fmt.Sprintf("%s%f", match[1], val), 1)
					} else {
						syntax = strings.Replace(syntax, match[0], fmt.Sprintf("%s%v", match[1], val), 1)
					}
					found = true
				}
			}
		}
		if !strings.Contains(syntax, "(") || !strings.Contains(syntax, ")") || !found {
			break
		}
	}
	if charCount('(', syntax)+charCount(')', syntax) > 0 {
		return nil, fmt.Errorf("invalid expression: \"%s\"", modifier)
	}
	if strings.Contains(syntax, ".") {
		if val, err := strconv.ParseFloat(syntax, 64); err == nil {
			return val, nil
		}
		return nil, fmt.Errorf("invalid expression: \"%s\"", modifier)
	}
	if val, err := strconv.ParseInt(syntax, 10, 64); err == nil {
		return val, nil
	}
	return nil, fmt.Errorf("invalid expression: \"%s\"", modifier)
}

// StringModifier implements Modifier for string transformations.
type StringModifier struct{}

// Valid checks if the modifier is a valid string transformation expression.
func (m *StringModifier) Valid(modifier string) bool {
	modifier = strings.Trim(modifier, " |") + "|"
	reg := regexp.MustCompile(`(?i)^(\s*(pre|post|sub|cut|flip|remove|delete)\((.*?)\)\s*\|)+$`)
	return reg.MatchString(modifier)
}

// Modify applies the string transformation modifier to the given value.
func (m *StringModifier) Modify(value string, modifier string, onfailed ...func(value any, modifier string) (any, error)) (any, error) {
	var result any = value
	syntax := strings.Trim(modifier, " |")
	reg := regexp.MustCompile(`(?i)^\s*(pre|post|sub|cut|flip|remove|delete)\((.*?)\)\s*$`)
	for _, expression := range stringSplit(syntax) {
		value := fmt.Sprintf("%v", result)
		match := reg.FindStringSubmatch(expression)
		if len(match) == 0 {
			if len(onfailed) == 0 || onfailed[0] == nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			val, err := onfailed[0](value, expression)
			if err != nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			result = val
			continue
		}
		switch match[1] {
		case "pre":
			if match[2] == "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			result = match[2] + value
		case "post":
			if match[2] == "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			result = value + match[2]
		case "sub":
			n, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			if n < 0 {
				n = int64(len(value)) + n
				if n < 0 {
					n = 0
				}
				result = value[n:]
			} else {
				if l := int64(len(value)); n > l {
					n = l
				}
				result = value[:n]
			}
		case "cut":
			n, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil || n == 0 {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			if n < 0 {
				n = int64(len(value)) + n
				if n < 0 {
					result = ""
				} else {
					result = value[:n]
				}
			} else {
				if l := int64(len(value)); n > l {
					result = ""
				} else {
					result = value[n:]
				}
			}
		case "flip":
			if match[2] != "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			flip := ""
			for _, char := range value {
				flip = string(char) + flip
			}
			result = flip
		case "remove":
			if match[2] == "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			result = strings.ReplaceAll(value, match[2], "")
		case "delete":
			if match[2] == "" {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			reg := regexp.MustCompile(`(?i)([` + regexp.QuoteMeta(match[2]) + `]+)`)
			result = reg.ReplaceAllString(value, "")
		default:
			return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
		}
	}
	return result, nil
}

// FormatModifier implements Modifier for formatting transformations.
type FormatModifier struct{}

// Valid checks if the modifier is a valid formatting expression.
func (m *FormatModifier) Valid(modifier string) bool {
	modifier = strings.Trim(modifier, " |") + "|"
	reg := regexp.MustCompile(`(?i)^(\s*(money|left|center|right)\((.*?)\)\s*\|)+$`)
	return reg.MatchString(modifier)
}

// Modify applies the formatting modifier to the given value.
func (m *FormatModifier) Modify(value string, modifier string, onfailed ...func(value any, modifier string) (any, error)) (any, error) {
	var result any = value
	syntax := strings.Trim(modifier, " |")
	reg := regexp.MustCompile(`(?i)^\s*(money|left|center|right)\((.*?)\)\s*$`)
	for _, expression := range stringSplit(syntax) {
		value := fmt.Sprintf("%v", result)
		match := reg.FindStringSubmatch(expression)
		if len(match) == 0 {
			if len(onfailed) == 0 || onfailed[0] == nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			val, err := onfailed[0](value, expression)
			if err != nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			result = val
			continue
		}
		switch match[1] {
		case "money":
			dec := 0
			sym := "."
			if match[2] != "" {
				reg := regexp.MustCompile(`^\s*([\.,]([0-9]*))\s*$`)
				m := reg.FindStringSubmatch(match[2])
				if len(m) <= 0 {
					return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
				}
				sym = m[1][:1]
				if m[2] != "" {
					dec, _ = strconv.Atoi(m[2])
				}
			}
			n, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			val := fmt.Sprintf("%."+strconv.Itoa(dec)+"f", n)
			if sym == "," {
				val = strings.ReplaceAll(val, ".", sym)
				sym = "."
			} else {
				sym = ","
			}
			reg := regexp.MustCompile(`^([1-9][0-9]*?)([0-9]{3})($|[\.,])`)
			for reg.MatchString(val) {
				val = reg.ReplaceAllString(val, "$1"+sym+"$2$3")
			}
			result = val
		case "left":
			n, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			if int(n) > len(value) {
				result = value + strings.Repeat(" ", int(n)-len(value))
			} else {
				result = value[:n]
			}
		case "center":
			n, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			if int(n) > len(value) {
				half := (int(n) - len(value)) / 2
				value = strings.Repeat(" ", half) + value
				value = value + strings.Repeat(" ", int(n)-len(value))
				result = value[:n]
			} else {
				result = value[:n]
			}
		case "right":
			n, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
			}
			if int(n) > len(value) {
				result = strings.Repeat(" ", int(n)-len(value)) + value
			} else {
				result = value[len(value)-int(n):]
			}
		default:
			return nil, fmt.Errorf("invalid expression: \"%s\"", expression)
		}
	}
	return result, nil
}

// Initialize default modifiers.
func init() {
	DefaultModifier(NewFormatModifier(), NewNumberModifier(), NewStringModifier())
}
