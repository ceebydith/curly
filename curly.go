package curly

import (
	"fmt"
	"regexp"
	"strings"
)

// Format applies a series of formatters to the given text and returns the formatted string.
func Format(text string, formatters ...Formatter) (string, error) {
	if len(formatters) == 0 {
		formatters = append(formatters, NewDatetimeFormatter(), NewDirectoryFormatter())
	}

	// Escape special characters
	result := strings.ReplaceAll(text, "%", "%25")
	result = strings.ReplaceAll(result, "\\\\", "%5C")
	result = strings.ReplaceAll(result, "\\{", "%7B")
	result = strings.ReplaceAll(result, "\\}", "%7D")

	// Regex to find placeholders in the text
	reg := regexp.MustCompile(`(?i)\{\s*([a-z]+([\._]?[a-z0-9]+)*)\s*([\*/\+\-:\|][^\}]+)?\}`)
	matches := reg.FindAllStringSubmatch(result, -1)

	for _, match := range matches {
		identifier := match[1]
		modifier := strings.TrimSpace(match[3])

		// Find the appropriate formatter
		var formatter Formatter
		for _, f := range formatters {
			if f.Valid(identifier) {
				formatter = f
				break
			}
		}
		if formatter == nil {
			return "", fmt.Errorf("invalid expression: \"%s\"", match[1])
		}

		// Get the value from the formatter
		value, err := formatter.Value(identifier)
		if err != nil {
			return "", err
		}

		// Apply the modifier if present
		if val, err := execModifier(value, modifier); err != nil {
			return "", fmt.Errorf("invalid expression: \"%s\"", match[0])
		} else {
			value = val
		}

		// Replace the placeholder with the formatted value
		result = strings.Replace(result, match[0], fmt.Sprintf("%v", value), 1)
	}

	// Make sure there are no remaining placeholders
	if strings.Contains(result, "{") || strings.Contains(result, "}") {
		return "", fmt.Errorf("invalid expression: \"%s\"", text)
	}

	// Restore escaped characters
	result = strings.ReplaceAll(result, "%7D", "}")
	result = strings.ReplaceAll(result, "%7B", "{")
	result = strings.ReplaceAll(result, "%5C", "\\")
	result = strings.ReplaceAll(result, "%25", "%")

	return result, nil
}

// Parse extracts data from text based on the provided expression and parsers.
func Parse[T string | []string](text string, expression T, parsers ...Parser) (map[string]any, error) {
	parsers = append(parsers, NewStringParser())
	result := map[string]any{}
	expressions := stringSplit(expression)
	if expressions == nil {
		return nil, fmt.Errorf("invalid expression type: %T", expression)
	}

	// Regex patterns
	regexAlphanum := regexp.MustCompile(`(?i)\{\s*(alphanum|alpha|num|any)(\s*:\s*([1-9][0-9]*))?\s*\}`)
	regexAlpha := regexp.MustCompile(`([@]+)`)
	regexNum := regexp.MustCompile(`([#]+)`)
	regexSpace := regexp.MustCompile(`\s+`)
	regexIdentifier := regexp.MustCompile(`(?i)\{\s*([a-z]+([\._]?[a-z0-9]+)*)\s*([\*/\+\-:\|][^\}]+)?\}`)
	var regexFinal *regexp.Regexp

	for _, expression := range expressions {
		replaces := [][2]string{}
		expression = strings.ReplaceAll(expression, "%", "%25")
		expression = strings.ReplaceAll(expression, "\\{", "%7B")
		expression = strings.ReplaceAll(expression, "\\}", "%7D")
		expression = strings.ReplaceAll(expression, "\\#", "%23")
		expression = strings.ReplaceAll(expression, "\\@", "%40")
		exp := expression

		// Process alphanum, alpha, num, and any patterns
		for _, match := range regexAlphanum.FindAllStringSubmatch(exp, -1) {
			var count = ""
			if match[3] != "" {
				count = fmt.Sprintf("{%s}", match[3])
			} else if match[1] == "any" {
				count = "*?"
			} else {
				count = "+?"
			}
			switch match[1] {
			case "num":
				replaces = append(replaces, [2]string{match[0], "[0-9]" + count})
			case "alpha":
				replaces = append(replaces, [2]string{match[0], "[a-z\\s]" + count})
			case "alphanum":
				replaces = append(replaces, [2]string{match[0], "[a-z0-9\\s]" + count})
			case "any":
				replaces = append(replaces, [2]string{match[0], "." + count})
			}
			exp = strings.Replace(exp, match[0], "", 1)
		}

		// Process alpha patterns
		for _, alpha := range regexAlpha.FindAllStringSubmatch(exp, -1) {
			var n string
			if len(alpha[1]) > 1 {
				n = fmt.Sprintf("{%d}", len(alpha[1]))
			}
			replaces = append(replaces, [2]string{alpha[0], fmt.Sprintf("[a-z]%s", n)})
			exp = strings.Replace(exp, alpha[0], "", 1)
		}

		// Process num patterns
		for _, num := range regexNum.FindAllStringSubmatch(exp, -1) {
			var n string
			if len(num[1]) > 1 {
				n = fmt.Sprintf("{%d}", len(num[1]))
			}
			replaces = append(replaces, [2]string{num[0], fmt.Sprintf("[0-9]%s", n)})
			exp = strings.Replace(exp, num[0], "", 1)
		}

		// Process identifier patterns
		match := regexIdentifier.FindAllStringSubmatch(exp, -1)
		if len(match) > 1 {
			return nil, fmt.Errorf("multiple identifier: \"%s\"", expression)
		}

		var parser Parser
		var target [3]string
		if len(match) == 1 {
			target = [3]string{match[0][0], match[0][1], strings.TrimSpace(match[0][3])}
			for _, p := range parsers {
				if p.Valid(target[1]) {
					parser = p
					break
				}
			}
			if parser == nil {
				return nil, fmt.Errorf("invalid expression : \"%s\"", expression)
			}
		}

		// Restore escaped characters
		expression = strings.ReplaceAll(expression, "%40", "@")
		expression = strings.ReplaceAll(expression, "%23", "#")
		expression = strings.ReplaceAll(expression, "%7D", "}")
		expression = strings.ReplaceAll(expression, "%7B", "{")
		expression = strings.ReplaceAll(expression, "%25", "%")
		exp = regexp.QuoteMeta(expression)
		for _, replace := range replaces {
			exp = strings.Replace(exp, regexp.QuoteMeta(replace[0]), replace[1], 1)
		}

		// Apply regex to extract data
		if parser == nil {
			exp = `(?i)` + regexSpace.ReplaceAllString(exp, `\s+`)
			regexFinal = regexp.MustCompile(exp)
			if !regexFinal.MatchString(text) {
				return nil, fmt.Errorf("invalid expression : \"%s\"", expression)
			}
		} else {
			parsed := false
			for i, regex := range parser.Expressions() {
				treg := regexp.QuoteMeta(target[0])
				if m, n := len(exp), len(treg); n > m || exp[m-n:] != treg {
					regex = regex + "?"
				} else {
					regex = regex + "$"
				}
				regex = `(?i)` + strings.Replace(exp, treg, "("+regex+")", 1)
				regex = regexSpace.ReplaceAllString(regex, `\s+`)
				regexFinal = regexp.MustCompile(regex)
				if match := regexFinal.FindStringSubmatch(text); match != nil {
					var value any
					value = parser.Modify(match[1], i)

					// Apply the modifier if present
					if val, err := execModifier(value, target[2]); err != nil {
						return nil, fmt.Errorf("invalid expression : \"%s\"", target[0])
					} else {
						value = val
					}
					result[target[1]] = value
					parsed = true
					break
				}
			}
			if !parsed {
				return nil, fmt.Errorf("invalid expression : \"%s\"", expression)
			}
		}
	}
	return result, nil
}

// NumberCalculate evaluates a mathematical expression after formatting it.
func NumberCalculate(expression string, formatters ...Formatter) (any, error) {
	expression, err := Format(expression, formatters...)
	if err != nil {
		return nil, err
	}
	return NewNumberModifier().Modify("", expression)
}

// StringModify applies modifications to the given text based on provided expressions.
func StringModify[T string | []string](text string, expressions T, formatters ...Formatter) (string, error) {
	expression := stringJoin(expressions)
	expression, err := Format(expression, formatters...)
	if err != nil {
		return "", err
	}
	val, err := NewStringModifier().Modify(text, expression)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", val), nil
}
