package curly_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ceebydith/curly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	testFormat := formatTester{
		format: func(text string) (string, error) {
			mapFormatter := curly.NewMapFormatter(map[string]any{
				"appname": "curly",
				"index":   0,
				"file":    "text.txt",
			})
			return curly.Format(text, mapFormatter, curly.NewDatetimeFormatter(), curly.NewDirectoryFormatter())
		},
		scenarios: []formatScenarioTest{
			{
				"{appdir}/log/{yyyy}/{mm}/{dd}/{appname}_{yyyy}{mm}{dd}.log",
				func() string {
					str := []string{}
					exe, _ := os.Executable()
					str = append(str, filepath.Dir(exe))
					str = append(str, "log")
					str = append(str, time.Now().Format("2006"))
					str = append(str, time.Now().Format("01"))
					str = append(str, time.Now().Format("02"))
					str = append(str, "curly_"+time.Now().Format("20060102")+".log")
					return strings.Join(str, "/")
				}(),
				nil,
			},
			{
				"{curdir}/log/{yyyy}/{mm}/{dd}/{appname}_{yyyy}{mm}{dd}.log",
				func() string {
					str := []string{}
					dir, _ := os.Getwd()
					str = append(str, dir)
					str = append(str, "log")
					str = append(str, time.Now().Format("2006"))
					str = append(str, time.Now().Format("01"))
					str = append(str, time.Now().Format("02"))
					str = append(str, "curly_"+time.Now().Format("20060102")+".log")
					return strings.Join(str, "/")
				}(),
				nil,
			},
			{
				"{workdir}/log/{yyyy}/{mm}/{dd}/{appname}_{yyyy}{mm}{dd}.log",
				"",
				fmt.Errorf("invalid expression: \"workdir\""),
			},
			{
				"{{curdir}}/log/{yyyy}/{mm}/{dd}/{appname}_{yyyy}{mm}{dd}.log",
				"",
				fmt.Errorf("invalid expression: \"{{curdir}}/log/{yyyy}/{mm}/{dd}/{appname}_{yyyy}{mm}{dd}.log\""),
			},
			{
				"Application {appname}, filename {file|remove(.txt)}, posistion {index+1}",
				"Application curly, filename text, posistion 1",
				nil,
			},
		},
	}

	t.Run("Format", testFormat.Test)
}

func TestParse(t *testing.T) {
	testParseString := parseTester[string]{
		parse: func(text string, expression string) (map[string]any, error) {
			return curly.Parse(text, expression, curly.NewNumberParser("index", "age"))
		},
		scenarios: []parseScenarioTest[string]{
			{
				"Message#1: Hello, my name is John Doe, I am 30 years old.",
				"Message\\#{index-1}:|name is {name},|I am {age} years",
				map[string]any{"index": int64(0), "name": "John Doe", "age": int64(30)},
				nil,
			},
		},
	}

	testParseStringList := parseTester[[]string]{
		parse: func(text string, expression []string) (map[string]any, error) {
			return curly.Parse(text, expression, curly.NewNumberParser("index", "age"))
		},
		scenarios: []parseScenarioTest[[]string]{
			{
				"Message#1: Hello, my name is Mr. John Doe, I am 30 years old.",
				[]string{"Message\\#{index-1}:", "name is {name|remove(Mr. )},", "I am {age} years"},
				map[string]any{"index": int64(0), "name": "John Doe", "age": int64(30)},
				nil,
			},
		},
	}

	t.Run("ParseString", testParseString.Test)
	t.Run("ParseStringList", testParseStringList.Test)
}

func TestNumberCalculate(t *testing.T) {
	testNumberCalculate := numberCalculateTester{
		calculate: func(expression string) (any, error) {
			return curly.NumberCalculate(expression)
		},
		scenarios: []numberCalculateScenarioTest{
			{"1+2", int64(3), nil},
		},
	}

	t.Run("NumberCalculate", testNumberCalculate.Test)
}

func TestStringModify(t *testing.T) {
	testStringModifyString := stringModifyTester[string]{
		modify: func(text string, expression string) (string, error) {
			return curly.StringModify(text, expression, curly.NewMapFormatter(map[string]any{
				"name":    "John Doe",
				"amount":  20000,
				"product": "OVO",
				"msisdn":  "081234567890",
			}))
		},
		scenarios: []stringModifyScenarioTest[string]{
			{
				"202412060646200001",
				"post(/{product})|post(/{msisdn})|post(/{name})|post(/{amount})",
				"202412060646200001/OVO/081234567890/John Doe/20000",
				nil,
			},
		},
	}

	testStringModifyList := stringModifyTester[[]string]{
		modify: func(text string, expression []string) (string, error) {
			return curly.StringModify(text, expression, curly.NewMapFormatter(map[string]any{
				"name":    "John Doe",
				"amount":  20000,
				"product": "OVO",
				"msisdn":  "081234567890",
			}))
		},
		scenarios: []stringModifyScenarioTest[[]string]{
			{
				"202412060646200001",
				[]string{"post(/{product})", "post(/{msisdn})", "post(/{name})", "post(/{amount})"},
				"202412060646200001/OVO/081234567890/John Doe/20000",
				nil,
			},
		},
	}

	t.Run("TestStringModifyString", testStringModifyString.Test)
	t.Run("TestStringModifyList", testStringModifyList.Test)
}

func TestCurlyPower(t *testing.T) {
	response := `TRX 2189566, PLN Prepaid 20000 (507) ke 133312626789 Harga 20075 ke 133312626789 (MBOK DARMI               ) status:SUCCESSFUL TOKEN:1582.4499.3217.5678.1234 tarif:R1 / 2200 VA kwh:1260 KWM ref:9C10530281BA4A8783792C775BF55ABC rp:Rp18.181 ppj:Rp1.819 orderid:1729153652308900832 info: 081234567890 SaldoAwal 2925146, SaldoAkhir 2905071`
	data, err := curly.Parse(
		response,
		[]string{
			"TRX {id},",
			"TRX {num}, {product} {num}",
			"TRX {num}, {alpha} {denom} ",
			"ke {dest} Harga",
			"Harga {charge} ke",
			"ke {alphanum} ({pln.nama})",
			"status:SUCCESSFUL",
			"TOKEN:{token} tarif",
			"tarif:{pln.tarifdaya|delete( )} kwh",
			"kwh:{pln.kwh/100} KWM",
			"ref:{pln.ref} rp",
			"rp:Rp{pln.tagihan} ppj",
			"ppj:Rp{pln.ppj} orderid",
			"orderid:{orderid} info",
			"info: {info} SaldoAwal",
			"SaldoAkhir {balance}",
			"### ({code}) @@",
			"TRX {any} ({product.code}) ke {num:11}",
		},
		NewTokenParser("token"),
		curly.NewMsisdnParser(62, 1, "info"),
		curly.NewNumberParser("denom", "charge", "pln.kwh", "pln.tagihan", "pln.ppj", "balance"),
	)
	require.NoError(t, err)

	total, err := curly.NumberCalculate("{pln.tagihan} + {pln.ppj}", curly.NewMapFormatter(data))
	require.NoError(t, err)
	data["pln.total"] = total

	bill := strings.Join([]string{
		"",
		"====================================",
		"         BUKTI PERMBAYARAN",
		"{product|center(36)}",
		"",
		"Nominal        : {denom}",
		"No. Pelanggan  : {dest}",
		"Nama Pelanggan : {pln.nama}",
		"Tarif/Daya     : {pln.tarifdaya}",
		"KWH            : {pln.kwh} KWM",
		"Tagihan        : Rp. {pln.tagihan|money(,)|right(15)}",
		"PPJ            : Rp. {pln.ppj|money(,)|right(15)}",
		"Total          : Rp. {pln.total|money(,)|right(15)|post(YEY)|cut(-3)}",
		"",
		"            TERIMAKASIH",
		"{info|pre(Info: )|center(36)}",
		"====================================",
	}, "\n")

	bill, err = curly.Format(bill, curly.NewMapFormatter(data))
	require.NoError(t, err)

	//t.Logf("%#v", data)
	t.Log(bill)
}

type formatScenarioTest struct {
	text         string
	expectFormat string
	expectError  error
}

type formatTester struct {
	format    func(text string) (string, error)
	scenarios []formatScenarioTest
}

func (tester *formatTester) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d %s", i, scenario.text)
		format, err := tester.format(scenario.text)
		require.Equal(t, scenario.expectFormat, format, "Format "+msg)
		assert.Equal(t, scenario.expectError, err, "Format "+msg)
	}
}

type parseScenarioTest[T string | []string] struct {
	text        string
	expression  T
	expectParse map[string]any
	expectError error
}

type parseTester[T string | []string] struct {
	parse     func(text string, expression T) (map[string]any, error)
	scenarios []parseScenarioTest[T]
}

func (tester *parseTester[T]) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d %s", i, scenario.text)
		parse, err := tester.parse(scenario.text, scenario.expression)
		require.Equal(t, scenario.expectParse, parse, "Parse "+msg)
		assert.Equal(t, scenario.expectError, err, "Parse "+msg)
	}
}

type numberCalculateScenarioTest struct {
	expression      string
	expectCalculate any
	expectError     error
}

type numberCalculateTester struct {
	calculate func(expression string) (any, error)
	scenarios []numberCalculateScenarioTest
}

func (tester *numberCalculateTester) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d %s", i, scenario.expression)
		calculate, err := tester.calculate(scenario.expression)
		require.Equal(t, scenario.expectCalculate, calculate, "Calculate "+msg)
		assert.Equal(t, scenario.expectError, err, "Calculate "+msg)
	}
}

type stringModifyScenarioTest[T string | []string] struct {
	text         string
	expression   T
	expectModify string
	expectError  error
}

type stringModifyTester[T string | []string] struct {
	modify    func(text string, expression T) (string, error)
	scenarios []stringModifyScenarioTest[T]
}

func (tester *stringModifyTester[T]) Test(t *testing.T) {
	for i, scenario := range tester.scenarios {
		msg := fmt.Sprintf("#%d %s", i, scenario.text)
		modify, err := tester.modify(scenario.text, scenario.expression)
		require.Equal(t, scenario.expectModify, modify, "Modify "+msg)
		assert.Equal(t, scenario.expectError, err, "Modify "+msg)
	}
}

func NewTokenParser(identifiers ...string) *TokenParser {
	return &TokenParser{
		identifiers: identifiers,
	}
}

type TokenParser struct {
	identifiers []string
}

func (p *TokenParser) Valid(identifier string) bool {
	return curly.ValidIdentifier(identifier, p.identifiers)
}

func (p *TokenParser) Expressions() []string {
	return []string{
		`[0-9\s\-/\.]+`,
	}
}

func (p *TokenParser) Modify(value string, index int) any {
	reg := regexp.MustCompile(`[^0-9]+`)
	value = reg.ReplaceAllString(value, "")
	reg = regexp.MustCompile(`([0-9]{4})`)
	value = reg.ReplaceAllString(value, "-$1")
	value = strings.Trim(value, "-")
	return value
}
