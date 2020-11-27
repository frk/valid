package isvalid

import (
	"testing"
)

func Test(t *testing.T) {
	type testcase struct {
		val  string
		want bool
	}

	testtable := []struct {
		Name string
		Func func(string) bool
		pass []string
		fail []string
	}{{
		Name: "Int", Func: Int,
		pass: []string{
			"0",
			"1",
			"-0",
			"123",
			"-987",
			"+717",
			"13",
			"123",
			"+1",
			"01",
			"-01",
			"000",
			"1234567890",
		},
		fail: []string{
			"",
			"0.1",
			".01",
			"123e45",
			"abcdef",
			"      ",
		},
	}, {
		Name: "Uint", Func: Uint,
		pass: []string{
			"0",
			"1",
			"+0",
			"123",
			"+987",
			"13",
			"123",
			"+1",
			"01",
			"+01",
			"000",
			"1234567890",
		},
		fail: []string{
			"",
			"-0",
			".01",
			"0.1",
			"123e45",
			"abcdef",
			"      ",
			"-987654321",
		},
	}, {
		Name: "Float", Func: Float,
		pass: []string{
			"123",
			"123.",
			"123.123",
			"-123.123",
			"-0.123",
			"+0.123",
			"0.123",
			".0",
			"-.123",
			"+.123",
			"01.123",
			"-0.22250738585072011e-307",
		},
		fail: []string{
			"+",
			"-",
			"  ",
			"",
			".",
			"foo",
			"20.foo",
			"2020-01-06T14:31:00.135Z",
		},
	}, {
		Name: "Binary", Func: Binary,
		pass: []string{
			"0",
			"1",
			"0001",
			"1110",
			"0001011010101",
			"0b010",
			"0B010",
		},
		fail: []string{
			"",
			"0b",
			"0B",
			"B010",
			"b010",
			"0a010",
			"0A010",
			"012",
		},
	}, {
		Name: "Octal", Func: Octal,
		pass: []string{
			"076543210",
			"0o01234567",
		},
		fail: []string{
			"abcdefg",
			"012345678",
			"012345670c",
			"00c12345670c",
			"",
			"..",
		},
	}, {
		Name: "Hex", Func: Hex,
		pass: []string{
			"deadBEEF",
			"ff0044",
			"0xff0044",
			"0XfF0044",
			"0x0123456789abcDEF",
			"0X0123456789abcDEF",
			"0hfedCBA9876543210",
			"0HfedCBA9876543210",
			"0123456789abcDEF",
		},
		fail: []string{
			"abcdefg",
			"",
			"..",
			"0xa2h",
			"0xa20x",
			"0x0123456789abcDEFq",
			"0hfedCBA9876543210q",
			"01234q56789abcDEF",
		},
	}, {
		Name: "Numeric", Func: Numeric,
		pass: []string{
			"123",
			"00123",
			"-00123",
			"0",
			"-0",
			"+123",
			"123.123",
			"+000000",
		},
		fail: []string{
			" ",
			"",
			".",
		},
	}, {
		Name: "Digits", Func: Digits,
		pass: []string{
			"123",
			"00123",
			"0",
			"0987654321",
		},
		fail: []string{
			"12.3",
			"12e3",
			"-123",
			"+123",
			"xAF",
			" ",
			"",
			".",
		},
	}, {
		Name: "HexColor", Func: HexColor,
		pass: []string{
			"#ff0000ff",
			"#ff0034",
			"#CCCCCC",
			"0f38",
			"fff",
			"#f00",
		},
		fail: []string{
			"#ff",
			"fff0a",
			"#ff12FG",
		},
	}}

	for _, tab := range testtable {
		t.Run(tab.Name, func(t *testing.T) {
			for _, val := range tab.pass {
				want := true
				t.Run(`"`+val+`"`, func(t *testing.T) {
					got := tab.Func(val)
					if got != want {
						t.Errorf("got=%t; want=%t", got, want)
					}
				})
			}

			for _, val := range tab.fail {
				want := false
				t.Run(`"`+val+`"`, func(t *testing.T) {
					got := tab.Func(val)
					if got != want {
						t.Errorf("got=%t; want=%t", got, want)
					}
				})
			}
		})
	}
}
