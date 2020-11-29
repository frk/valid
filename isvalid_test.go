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
	}, {
		Name: "Base32", Func: Base32,
		pass: []string{
			"ZG======",
			"JBSQ====",
			"JBSWY===",
			"JBSWY3A=",
			"JBSWY3DP",
			"JBSWY3DPEA======",
			"K5SWYY3PNVSSA5DPEBXG6ZA=",
			"K5SWYY3PNVSSA5DPEBXG6===",
		},
		fail: []string{
			"12345",
			"",
			"JBSWY3DPtesting123",
			"ZG=====",
			"Z======",
			"Zm=8JBSWY3DP",
			"=m9vYg==",
			"Zm9vYm/y====",
		},
	}, {
		Name: "Base58", Func: Base58,
		pass: []string{
			"BukQL",
			"3KMUV89zab",
			"91GHkLMNtyo98",
			"YyjKm3H",
			"Mkhss145TRFg",
			"7678765677",
			"abcodpq",
			"AAVHJKLPY",
		},
		fail: []string{
			"0OPLJH",
			"IMKLP23",
			"KLMOmk986",
			"LL1l1985hG",
			"*MP9K",
			"Zm=8JBSWY3DP",
			")()(=9292929MKL",
		},
	}, {
		Name: "Base64-normal", Func: func(v string) bool {
			return Base64(v, false)
		},
		pass: []string{
			"",
			"Zg==",
			"Zm8=",
			"Zm9v",
			"Zm9vYg==",
			"Zm9vYmE=",
			"Zm9vYmFy",
			"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=",
			"Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==",
			"U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==",
			"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
				"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
				"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
				"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
				"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
				"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" +
				"HQIDAQAB",
		},
		fail: []string{
			"12345",
			"Vml2YW11cyBmZXJtZtesting123",
			"Zg=",
			"Z===",
			"Zm=8",
			"=m9vYg==",
			"Zm9vYmFy====",
		},
	}, {
		Name: "Base64-urlsafe", Func: func(v string) bool {
			return Base64(v, true)
		},
		pass: []string{
			"",
			"bGFkaWVzIGFuZCBnZW50bGVtZW4sIHdlIGFyZSBmbG9hdGluZyBpbiBzcGFjZQ",
			"1234",
			"bXVtLW5ldmVyLXByb3Vk",
			"PDw_Pz8-Pg",
			"VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw",
		},
		fail: []string{
			" AA",
			"\tAA",
			"\rAA",
			"\nAA",
			"This+isa/bad+base64Url==",
			"0K3RgtC+INC30LDQutC+0LTQuNGA0L7QstCw0L3QvdCw0Y8g0YHRgtGA0L7QutCw",
		},
	}, {
		Name: "Email", Func: Email,
		pass: []string{
			"foo@bar.com",
			"x@x.au",
			"foo@bar.com.au",
			"foo+bar@bar.com",
			"hans.m端ller@test.com",
			"hans@m端ller.com",
			"test|123@m端ller.com",
			"test123+ext@gmail.com",
			"some.name.midd.leNa.me+extension@GoogleMail.com",
			`"foobar"@example.com`,
			`"  foo  m端ller "@example.com`,
			`"foo\\@bar"@example.com`,
			"test@gmail.com",
			"test.1@gmail.com",
		},
		fail: []string{
			`invalidemail@`,
			`invalid.com`,
			`@invalid.com`,
			`foo@bar.com.`,
			`foo@bar.co.uk.`,
			`multiple..dots@stillinvalid.com`,
			`test123+invalid! sub_address@gmail.com`,
			`gmail...ignores...dots...@gmail.com`,
			`ends.with.dot.@gmail.com`,
			`multiple..dots@gmail.com`,
			`wrong()[]",:;<>@@gmail.com`,
			`"wrong()[]",:;<>@@gmail.com`,
		},
	}, {
		Name: "MAC-0", Func: func(v string) bool {
			return MAC(v, 0)
		},
		pass: []string{
			"08:00:2b:01:02:03",
			"08-00-2b-01-02-03",
			"01:AB:03:04:05:06",
			"01-02-03-04-05-ab",
			"08002b010203",
			"08:00:2b:01:02:03:04:05",
			"08-00-2b-01-02-03-04-05",
			"08002b0102030405",
		},
		fail: []string{
			"A9 C5 D4 9F EB D3",
			"01 02 03 04 05 ab",
			"0102.0304.05ab",
			"08002b:010203",
			"08002b-010203",
			"0800.2b01.0203",
			"0800.2b01.0203.0405",
			"abc",
			"01:02:03:04:05",
			"01:02:03:04::ab",
			"1:2:3:4:5:6",
			"AB:CD:EF:GH:01:02",
			"A9C5 D4 9F EB D3",
			"01-02 03:04 05 ab",
			"0102.03:04.05ab",
		},
	}, {
		Name: "MAC-6", Func: func(v string) bool {
			return MAC(v, 6)
		},
		pass: []string{
			"08:00:2b:01:02:03",
			"08-00-2b-01-02-03",
			"01:AB:03:04:05:06",
			"01-02-03-04-05-ab",
			"08002b010203",
		},
		fail: []string{
			"08:00:2b:01:02:03:04:05",
			"08-00-2b-01-02-03-04-05",
			"08002b0102030405",
			"A9 C5 D4 9F EB D3",
			"01 02 03 04 05 ab",
			"0102.0304.05ab",
			"08002b:010203",
			"08002b-010203",
			"0800.2b01.0203",
			"0800.2b01.0203.0405",
			"abc",
			"01:02:03:04:05",
			"01:02:03:04::ab",
			"1:2:3:4:5:6",
			"AB:CD:EF:GH:01:02",
			"A9C5 D4 9F EB D3",
			"01-02 03:04 05 ab",
			"0102.03:04.05ab",
		},
	}, {
		Name: "MAC-8", Func: func(v string) bool {
			return MAC(v, 8)
		},
		pass: []string{
			"08:00:2b:01:02:03:04:05",
			"08-00-2b-01-02-03-04-05",
			"08002b0102030405",
		},
		fail: []string{
			"08:00:2b:01:02:03",
			"08-00-2b-01-02-03",
			"01:AB:03:04:05:06",
			"01-02-03-04-05-ab",
			"08002b010203",
			"A9 C5 D4 9F EB D3",
			"01 02 03 04 05 ab",
			"0102.0304.05ab",
			"08002b:010203",
			"08002b-010203",
			"0800.2b01.0203",
			"0800.2b01.0203.0405",
			"abc",
			"01:02:03:04:05",
			"01:02:03:04::ab",
			"1:2:3:4:5:6",
			"AB:CD:EF:GH:01:02",
			"A9C5 D4 9F EB D3",
			"01-02 03:04 05 ab",
			"0102.03:04.05ab",
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
