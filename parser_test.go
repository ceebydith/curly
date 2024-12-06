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

// func (p *TokenParser) Modify(value string, index int) any {
// 	reg := regexp.MustCompile(`[^0-9]+`)
// 	value = reg.ReplaceAllString(value, "")
// 	reg = regexp.MustCompile(`([0-9]{4})`)
// 	value = reg.ReplaceAllString(value, "-$1")
// 	value = strings.Trim(value, "-")
// 	return value
// }

// func TestParse(t *testing.T) {
// 	text := `TRX 2189561, PLN Prepaid 20000 (507) ke 133312627156 Harga 20075.50 ke 133312627156 (YENI PUTRI               ) status:SUCCESSFUL TOKEN:1583-4469-3407-5632-1298 tarif:R1 / 2200 VA kwh:1260 KWM ref:6B10530281CD4A7083792C775BF55BEC rp:Rp18.181 ppj:Rp1.819 orderid:1729153742308900206 SaldoAwal 2925146, SaldoAkhir 2905071`

// 	data, err := curly.Parse(
// 		text,
// 		[]string{
// 			`{ket}, dari`,
// 			`TRX {id},`,
// 			`dari {msisdn},`,
// 			`PLN Prepaid {denom} `,
// 			`ke {dest} Harga`,
// 			`Harga {charge} ke`,
// 			` ke {num} ({pln.nama}) status`,
// 			`status:SUCCESSFUL`,
// 			`TOKEN:{token} tarif`,
// 			`tarif:{pln.tarif|delete( )} kwh`,
// 			`kwh:{pln.kwh}{num:2} KWM ref`,
// 			`kwh:{pln.kwhdes|sub(-2)} KWM ref`,
// 			`ref:{pln.ref} `,
// 			`SaldoAkhir {saldo};`,
// 			`SN={token}/`,
// 			`SN={any}/{dua}/`,
// 			`SN={any}/{any}/{tiga}/`,
// 			`SN={any}/{any}/{any}/{empat}/{num}.{any}`,
// 			`SN={any}/{any}/{any}/{any}/{any}/{lima}`,
// 		},
// 		//`TRX {id},|dari {msisdn},|PLN Prepaid {denom+100} |ke {dest} Harga|Harga {charge} ke| ke {num} ({pln.nama}) status|status:SUCCESSFUL|TOKEN:{token} tarif|tarif:{pln.tarif\|delete( )} kwh|kwh:{pln.kwh}{num:2} KWM ref|kwh:{pln.kwhdes\|sub(-2)} KWM ref|ref:{pln.ref} |SaldoAkhir {saldo}`,
// 		NewTokenParser("token"),
// 		curly.NewMsisdnParser(62, -1, "msisdn"),
// 		curly.NewNumberParser("denom", "charge", "saldo"),
// 		//curly.NewStringParser(false),
// 	)
// 	require.NoError(t, err)
// 	t.Logf("%#v %T", data, data["denom"])

// 	//modifier := `post(\|{pln.nama})|post(\|{dest})|post(\|{pln.kwh})|post(.{pln.kwhdes})`
// 	modifier := []string{
// 		`post(/{pln.nama})`,
// 		`post(/{dest})`,
// 		`post(/{pln.tarif})`,
// 		`post(/{pln.kwh})`,
// 		`post(.{pln.kwhdes})`,
// 	}

// 	sn, err := curly.StringModify(data["token"].(string), modifier, curly.NewMapFormatter(data), curly.NewDatetimeFormatter())
// 	require.NoError(t, err)
// 	t.Logf("%#v", sn)
// }
