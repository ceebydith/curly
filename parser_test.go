package curly_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ceebydith/curly"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	testNumber := parserTester{
		parser: curly.NewNumberParser("number", "amount", "total"),
		scenarios: []parserScenarioTest{
			{"number", "100", true, true, int64(100)},
			{"number", "100.00", true, true, float64(100)},
			{"amount", "100,00", true, true, float64(100)},
			{"amount", "100.000,00", true, true, float64(100000)},
			{"total", "100.000.00", true, true, nil},
			{"total", "-100.000", true, true, int64(-100000)},
			{"fee", "200", false, true, int64(200)},
		},
	}

	testString := parserTester{
		parser: curly.NewStringParser(),
		scenarios: []parserScenarioTest{
			{"name", "John Doe", true, true, "John Doe"},
			{"user", "root", true, true, "root"},
			{"email", "you@example.com", true, true, "you@example.com"},
			{"password", "p455w0rd", true, true, "p455w0rd"},
			{"address", "Mountain St. ", true, true, "Mountain St."},
			{"status", "1", true, true, "1"},
		},
	}

	testStringNoTrim := parserTester{
		parser: curly.NewStringParser(false),
		scenarios: []parserScenarioTest{
			{"name", "John Doe ", true, true, "John Doe "},
			{"user", " root", true, true, " root"},
			{"address", " Mountain St. ", true, true, " Mountain St. "},
		},
	}

	testMsisdn := parserTester{
		parser: curly.NewMsisdnParser(62, 0, "msisdn", "phone"),
		scenarios: []parserScenarioTest{
			{"msisdn", "081234567890", true, true, "6281234567890"},
			{"phone", "081234567890", true, true, "6281234567890"},
			{"phone", "+6281234567890", true, true, "6281234567890"},
			{"phone", "6281234567890", true, true, "6281234567890"},
			{"msisdn", "81234567890", true, true, nil},
			{"whatsapp", "081234567890", false, true, "6281234567890"},
		},
	}

	testMsisdnPlus := parserTester{
		parser: curly.NewMsisdnParser(62, 1, "msisdn", "phone"),
		scenarios: []parserScenarioTest{
			{"msisdn", "081234567890", true, true, "+6281234567890"},
			{"phone", "081234567890", true, true, "+6281234567890"},
			{"phone", "+6281234567890", true, true, "+6281234567890"},
			{"phone", "6281234567890", true, true, "+6281234567890"},
			{"msisdn", "81234567890", true, true, nil},
			{"whatsapp", "081234567890", false, true, "+6281234567890"},
		},
	}

	testMsisdnZero := parserTester{
		parser: curly.NewMsisdnParser(62, -1, "msisdn", "phone"),
		scenarios: []parserScenarioTest{
			{"msisdn", "081234567890", true, true, "081234567890"},
			{"phone", "081234567890", true, true, "081234567890"},
			{"phone", "+6281234567890", true, true, "081234567890"},
			{"phone", "6281234567890", true, true, "081234567890"},
			{"msisdn", "81234567890", true, true, nil},
			{"whatsapp", "081234567890", false, true, "081234567890"},
		},
	}

	testMsisdnNoCode := parserTester{
		parser: curly.NewMsisdnParser(0, 0, "msisdn", "phone"),
		scenarios: []parserScenarioTest{
			{"msisdn", "081234567890", true, true, "081234567890"},
			{"phone", "081234567890", true, true, "081234567890"},
			{"phone", "+6281234567890", true, true, nil},
			{"phone", "6281234567890", true, true, nil},
			{"msisdn", "81234567890", true, true, nil},
			{"whatsapp", "081234567890", false, true, "081234567890"},
		},
	}

	t.Run("TestNumberParser", testNumber.Test)
	t.Run("TestStringParser", testString.Test)
	t.Run("TestStringNoTrimParser", testStringNoTrim.Test)
	t.Run("TestMsisdnParser", testMsisdn.Test)
	t.Run("TestMsisdnPlusParser", testMsisdnPlus.Test)
	t.Run("TestMsisdnZeroParser", testMsisdnZero.Test)
	t.Run("TestMsisdnNoCodeParser", testMsisdnNoCode.Test)
}

type parserScenarioTest struct {
	identifier       string
	value            string
	expectValid      bool
	expectExpression bool
	expectModify     any
}

type parserTester struct {
	parser    curly.Parser
	scenarios []parserScenarioTest
}

func (tester *parserTester) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d: %s", i, scenario.identifier)

		valid := tester.parser.Valid(scenario.identifier)
		require.Equal(t, scenario.expectValid, valid, "Valid "+msg)

		expressions := tester.parser.Expressions()
		if !scenario.expectExpression {
			require.Empty(t, expressions, "Expressions "+msg)
		} else {
			var value any
			for i, expression := range expressions {
				reg := regexp.MustCompile(`^` + expression + `$`)
				if reg.MatchString(scenario.value) {
					value = tester.parser.Modify(scenario.value, i)
					break
				}
			}
			require.Equal(t, scenario.expectModify, value, "Modify "+msg)
		}
	}
}
