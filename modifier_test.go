package curly_test

import (
	"fmt"
	"testing"

	"github.com/ceebydith/curly"
	"github.com/stretchr/testify/require"
)

func TestModifier(t *testing.T) {
	testNumber := modifierTester{
		modifier: curly.NewNumberModifier(),
		scenarios: []modifierScenarioTest{
			{"100", "+100", true, int64(200), nil},
			{"100", "+ 100", true, int64(200), nil},
			{"100", " +100", true, int64(200), nil},
			{"100", "+100 ", true, int64(200), nil},
			{"100", " + 100", true, int64(200), nil},
			{"100", "+ 100 ", true, int64(200), nil},
			{"100", " +100 ", true, int64(200), nil},
			{"100", " + 100 ", true, int64(200), nil},
			{"300", "-100", true, int64(200), nil},
			{"300", "- 100", true, int64(200), nil},
			{"300", " -100", true, int64(200), nil},
			{"300", "-100 ", true, int64(200), nil},
			{"300", " - 100", true, int64(200), nil},
			{"300", "- 100 ", true, int64(200), nil},
			{"300", " -100 ", true, int64(200), nil},
			{"300", " - 100 ", true, int64(200), nil},
			{"100", "*3", true, int64(300), nil},
			{"100", "* 3", true, int64(300), nil},
			{"100", " *3", true, int64(300), nil},
			{"100", "*3 ", true, int64(300), nil},
			{"100", " * 3", true, int64(300), nil},
			{"100", "* 3 ", true, int64(300), nil},
			{"100", " *3 ", true, int64(300), nil},
			{"100", " * 3 ", true, int64(300), nil},
			{"100", "/2", true, int64(50), nil},
			{"100", "/ 2", true, int64(50), nil},
			{"100", " /2", true, int64(50), nil},
			{"100", "/2 ", true, int64(50), nil},
			{"100", " / 2", true, int64(50), nil},
			{"100", "/ 2 ", true, int64(50), nil},
			{"100", " /2 ", true, int64(50), nil},
			{"100", " / 2 ", true, int64(50), nil},
			{"100", "*-3", true, int64(-300), nil},
			{"100", "* -3", true, int64(-300), nil},
			{"100", " *-3", true, int64(-300), nil},
			{"100", "*-3 ", true, int64(-300), nil},
			{"100", " * -3", true, int64(-300), nil},
			{"100", "* -3 ", true, int64(-300), nil},
			{"100", " *-3 ", true, int64(-300), nil},
			{"100", " * -3 ", true, int64(-300), nil},
			{"100", "+100.00", true, float64(200), nil},
			{"100.00", "+100", true, float64(200), nil},
			{"100.00", "+100.00", true, float64(200), nil},
			{"100", "+100,00", false, nil, fmt.Errorf("invalid expression: \"%s%s\"", "100", "+100,00")},
			{"100,00", "+100", true, nil, fmt.Errorf("invalid expression: \"%s%s\"", "100,00", "+100")},
			{"100,00", "+100,00", false, nil, fmt.Errorf("invalid expression: \"%s%s\"", "100,00", "+100,00")},
			{"100.00", "+100.00)", false, nil, fmt.Errorf("invalid expression: \"%s%s\"", "100.00", "+100.00)")},
			{"100.00", "+100.00*/200", true, nil, fmt.Errorf("invalid expression: \"%s%s\"", "100.00", "+100.00*/200")},
			{"100", " *-3 / 0 ", true, nil, fmt.Errorf("division by zero: \"%s%s\"", "100", " *-3 / 0 ")},
		},
	}

	testString := modifierTester{
		modifier: curly.NewStringModifier(),
		scenarios: []modifierScenarioTest{
			{"abcdef", "|pre(XYZ)", true, "XYZabcdef", nil},
			{"abcdef", "|post(XYZ) ", true, "abcdefXYZ", nil},
			{"abcdef", "| sub(3)", true, "abc", nil},
			{"abcdef", "| sub(-3) ", true, "def", nil},
			{"abcdef", " |cut(3)", true, "def", nil},
			{"abcdef", " |cut(-3) ", true, "abc", nil},
			{"abcdef", " | flip()", true, "fedcba", nil},
			{"abcdef", " | remove(cde) ", true, "abf", nil},
			{"aabbccddeeff", " | delete(cde)", true, "aabbff", nil},
			{"abcdefghi", "|replace(ghi,xyz)", true, "abcdefxyz", nil},
			{"abcdef", "pre()", true, nil, fmt.Errorf("invalid expression: \"%s\"", "pre()")},
			{"abcdef", "post()", true, nil, fmt.Errorf("invalid expression: \"%s\"", "post()")},
			{"abcdef", "sub(x)", true, nil, fmt.Errorf("invalid expression: \"%s\"", "sub(x)")},
			{"abcdef", "cut()", true, nil, fmt.Errorf("invalid expression: \"%s\"", "cut()")},
			{"abcdef", "flip(11)", true, nil, fmt.Errorf("invalid expression: \"%s\"", "flip(11)")},
			{"abcdef", "remove()", true, nil, fmt.Errorf("invalid expression: \"%s\"", "remove()")},
			{"abcdef", "delete()", true, nil, fmt.Errorf("invalid expression: \"%s\"", "delete()")},
			{"abcdef", " |cut(30)", true, "", nil},
			{"abcdef", " |cut(-30)", true, "", nil},
			{"abcdef", " |sub(30)", true, "abcdef", nil},
			{"abcdef", " |sub(-30)", true, "abcdef", nil},
			{"abcdef", " |super()", false, nil, fmt.Errorf("invalid expression: \"%s\"", "super()")},
		},
	}

	testFormat := modifierTester{
		modifier: curly.NewFormatModifier(),
		scenarios: []modifierScenarioTest{
			{"100", "|money()", true, "100", nil},
			{"100.00", "|money(,)", true, "100", nil},
			{"1000000.00", "|money(,)", true, "1.000.000", nil},
			{"1000000", "|money(,)", true, "1.000.000", nil},
			{"1000000", "|money(,2)", true, "1.000.000,00", nil},
			{"1000000", "|money(.2)", true, "1,000,000.00", nil},
			{"Elon Musk", "|left(16)", true, "Elon Musk       ", nil},
			{"Elon Musk", "|center(16)", true, "   Elon Musk    ", nil},
			{"Elon Musk", "|right(16)", true, "       Elon Musk", nil},
		},
	}
	t.Run("TestNumberModifier", testNumber.Test)
	t.Run("TestStringModifier", testString.Test)
	t.Run("TestFormatModifier", testFormat.Test)
}

type modifierScenarioTest struct {
	value        string
	modifier     string
	expectValid  bool
	expectModify any
	expectError  error
}

type modifierTester struct {
	modifier  curly.Modifier
	scenarios []modifierScenarioTest
}

func (tester *modifierTester) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d: %s%s", i, scenario.value, scenario.modifier)

		valid := tester.modifier.Valid(scenario.modifier)
		require.Equal(t, scenario.expectValid, valid, "Valid "+msg)

		value, err := tester.modifier.Modify(scenario.value, scenario.modifier)
		require.Equal(t, scenario.expectError, err, "Modify "+msg)
		require.Equal(t, scenario.expectModify, value, "Modify "+msg)
	}
}
