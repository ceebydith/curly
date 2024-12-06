package curly

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ValidIdentifier(identifier string, identifiers []string) bool {
	identifier = strings.ToLower(identifier)
	for _, i := range identifiers {
		if strings.ToLower(i) == identifier {
			return true
		}
	}
	return false
}

func charCount(char rune, str string) int {
	count := 0
	for _, c := range str {
		if c == char {
			count++
		}
	}
	return count
}

func numberOf(num string) (float64, error) {
	n := strings.ReplaceAll(num, " ", "")
	v, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: \"%s\"", num)
	}
	return v, nil
}

func calculate(a string, opr string, b string) (float64, error) {
	x, err := numberOf(a)
	if err != nil {
		return 0, err
	}
	y, err := numberOf(b)
	if err != nil {
		return 0, err
	}
	switch strings.TrimSpace(opr) {
	case "+":
		return x + y, nil
	case "-":
		return x - y, nil
	case "*":
		return x * y, nil
	case "/":
		return x / y, nil
	}
	return 0, fmt.Errorf("invalid operator: \"%s\"", opr)
}

func stringSplit[T string | []string](str T) []string {
	switch val := any(str).(type) {
	case []string:
		return val
	case string:
		val = strings.Trim(val, "|")
		val = strings.ReplaceAll(val, "%", "%25")
		val = strings.ReplaceAll(val, "\\|", "%7C")
		result := strings.Split(val, "|")
		for i := range result {
			result[i] = strings.ReplaceAll(result[i], "%7C", "|")
			result[i] = strings.ReplaceAll(result[i], "%25", "%")
		}
		return result
	}
	return nil
}

func stringJoin[T string | []string](str T) string {
	switch val := any(str).(type) {
	case string:
		return val
	case []string:
		result := []string{}
		reg := regexp.MustCompile(`([^\\])\|`)
		for _, s := range val {
			result = append(result, reg.ReplaceAllString(s, "$1\\|"))
		}
		return strings.Join(result, "|")
	}
	return ""
}

func execModifier(value any, modifier string) (any, error) {
	if modifier == "" {
		return value, nil
	}

	var modif Modifier
	for _, m := range DefaultModifier() {
		if m.Valid(modifier) {
			modif = m
			break
		}
	}
	if modif == nil {
		return nil, fmt.Errorf("invalid modifier: \"%s\"", modifier)
	}
	return modif.Modify(fmt.Sprintf("%v", value), modifier, execModifier)
}
