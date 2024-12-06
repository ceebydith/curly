package curly_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ceebydith/curly"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	testMap := formatterTester{
		formatter: curly.NewMapFormatter(map[string]any{
			"id":     "20241205183927",
			"name":   "John Doe",
			"age":    30,
			"active": true,
			"amount": 53.51,
		}),
		scenarios: []formatterScenarioTest{
			{"id", true, "20241205183927", nil},
			{"name", true, "John Doe", nil},
			{"age", true, 30, nil},
			{"active", true, true, nil},
			{"amount", true, 53.51, nil},
			{"not_exist", false, nil, fmt.Errorf("invalid identifier: \"not_exist\"")},
		},
	}

	testDatetime := formatterTester{
		formatter: curly.NewDatetimeFormatter(),
		scenarios: []formatterScenarioTest{
			{"yyyy", true, time.Now().Format("2006"), nil},
			{"yy", true, time.Now().Format("06"), nil},
			{"mm", true, time.Now().Format("01"), nil},
			{"dd", true, time.Now().Format("02"), nil},
			{"hh", true, time.Now().Format("15"), nil},
			{"nn", true, time.Now().Format("04"), nil},
			{"ss", true, time.Now().Format("05"), nil},
			{"xx", false, nil, fmt.Errorf("invalid identifier: \"xx\"")},
		},
	}

	testDirectory := formatterTester{
		formatter: curly.NewDirectoryFormatter(),
		scenarios: []formatterScenarioTest{
			{"appdir", true, func() string { exe, _ := os.Executable(); return filepath.Dir(exe) }(), nil},
			{"curdir", true, func() string { dir, _ := os.Getwd(); return dir }(), nil},
			{"workdir", false, nil, fmt.Errorf("invalid identifier: \"workdir\"")},
		},
	}
	t.Run("TestMapFormatter", testMap.Test)
	t.Run("TestDatetimeFormatter", testDatetime.Test)
	t.Run("TestDirectoryFormatter", testDirectory.Test)
}

type formatterScenarioTest struct {
	identifier  string
	expectValid bool
	expectValue any
	expectError error
}

type formatterTester struct {
	formatter curly.Formatter
	scenarios []formatterScenarioTest
}

func (tester *formatterTester) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d %s", i, scenario.identifier)

		valid := tester.formatter.Valid(scenario.identifier)
		require.Equal(t, scenario.expectValid, valid, "Valid"+msg)

		value, err := tester.formatter.Value(scenario.identifier)
		require.Equal(t, scenario.expectValue, value, "Value "+msg)
		require.Equal(t, scenario.expectError, err, "Value "+msg)
	}
}
