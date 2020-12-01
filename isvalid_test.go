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
		Name: "ASCII", Func: ASCII,
		pass: []string{
			"foobar",
			"0987654321",
			"test@example.com",
			"1234abcDEF",
		},

		fail: []string{
			"ｆｏｏbar",
			"ｘｙｚ０９８",
			"１２３456",
			"ｶﾀｶﾅ",
		},
	}, {
		Name: "BIC", Func: BIC,
		pass: []string{
			"SBICKEN1345",
			"SBICKEN1",
			"SBICKENY",
			"SBICKEN1YYP",
		},

		fail: []string{
			"SBIC23NXXX",
			"S23CKENXXXX",
			"SBICKENXX",
			"SBICKENXX9",
			"SBICKEN13458",
			"SBICKEN",
		},
	}, {
		Name: "BTC", Func: BTC,
		pass: []string{
			"1MUz4VMYui5qY1mxUiG8BQ1Luv6tqkvaiL",
			"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
			"bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
		},

		fail: []string{
			"4J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
			"0x56F0B8A998425c53c75C4A303D4eF987533c5597",
			"pp8skudq3x5hzw8ew7vzsw8tn4k8wxsqsv0lt0mf3g",
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
		Name: "Bool", Func: Bool,
		pass: []string{
			"true",
			"false",
			"TRUE",
			"FALSE",
		},
		fail: []string{
			"",
			"True",
			"False",
			"0",
			"1",
			"t",
			"f",
		},
	}, {
		Name: "CIDR", Func: CIDR,
		pass: []string{
			"135.104.0.0/32",
			"0.0.0.0/24",
			"135.104.0.0/24",
			"135.104.0.1/32",
			"135.104.0.1/24",
			"::1/128",
			"abcd:2345::/127",
			"abcd:2345::/65",
			"abcd:2345::/64",
			"abcd:2345::/63",
			"abcd:2345::/33",
			"abcd:2345::/32",
			"abcd:2344::/31",
			"abcd:2300::/24",
			"abcd:2345::/24",
			"2001:DB8::/48",
			"2001:DB8::1/48",
		},
		fail: []string{
			"192.168.1.1/255.255.255.0",
			"192.168.1.1/35",
			"2001:db8::1/-1",
			"2001:db8::1/-0",
			"-0.0.0.0/32",
			"0.-1.0.0/32",
			"0.0.-2.0/32",
			"0.0.0.-3/32",
			"0.0.0.0/-0",
			"",
		},
	}, {
		Name: "CVV", Func: CVV,
		pass: []string{
			"123",
			"1234",
		},
		fail: []string{
			"",
			"12",
			"12345",
			"abc",
			"abcd",
		},
	}, {
		Name: "Currency", Func: Currency,
		pass: []string{
			// TODO
		},
		fail: []string{
			// TODO
		},
	}, {
		Name: "DataURI", Func: DataURI,
		pass: []string{
			"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC",
			"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIBAMAAAA2IaO4AAAAFVBMVEXk5OTn5+ft7e319fX29vb5+fn///++GUmVAAAALUlEQVQIHWNICnYLZnALTgpmMGYIFWYIZTA2ZFAzTTFlSDFVMwVyQhmAwsYMAKDaBy0axX/iAAAAAElFTkSuQmCC",
			"   data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIBAMAAAA2IaO4AAAAFVBMVEXk5OTn5+ft7e319fX29vb5+fn///++GUmVAAAALUlEQVQIHWNICnYLZnALTgpmMGYIFWYIZTA2ZFAzTTFlSDFVMwVyQhmAwsYMAKDaBy0axX/iAAAAAElFTkSuQmCC   ",
			"data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22100%22%20height%3D%22100%22%3E%3Crect%20fill%3D%22%2300B1FF%22%20width%3D%22100%22%20height%3D%22100%22%2F%3E%3C%2Fsvg%3E",
			"data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIj48cmVjdCBmaWxsPSIjMDBCMUZGIiB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIvPjwvc3ZnPg==",
			" data:,Hello%2C%20World!",
			" data:,Hello World!",
			" data:text/plain;base64,SGVsbG8sIFdvcmxkIQ%3D%3D",
			" data:text/html,%3Ch1%3EHello%2C%20World!%3C%2Fh1%3E",
			"data:,A%20brief%20note",
			"data:text/html;charset=US-ASCII,%3Ch1%3EHello!%3C%2Fh1%3E",
		},
		fail: []string{
			"dataxbase64",
			"data:HelloWorld",
			"data:,A%20brief%20invalid%20[note",
			"file:text/plain;base64,SGVsbG8sIFdvcmxkIQ%3D%3D",
			"data:text/html;charset=,%3Ch1%3EHello!%3C%2Fh1%3E",
			"data:text/html;charset,%3Ch1%3EHello!%3C%2Fh1%3E", "data:base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC",
			"",
			"http://wikipedia.org",
			"base64",
			"iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC",
		},
	}, {
		Name: "Decimal", Func: Decimal,
		pass: []string{
			// TODO
		},
		fail: []string{
			// TODO
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
		Name: "EAN", Func: EAN,
		pass: []string{
			"9421023610112",
			"1234567890128",
			"4012345678901",
			"9771234567003",
			"9783161484100",
			"73513537",
		},
		fail: []string{
			"5901234123451",
			"079777681629",
			"0705632085948",
		},
	}, {
		Name: "EIN", Func: EIN,
		pass: []string{
			// TODO
		},
		fail: []string{
			// TODO
		},
	}, {
		Name: "ETH", Func: ETH,
		pass: []string{
			"0x0000000000000000000000000000000000000001",
			"0x683E07492fBDfDA84457C16546ac3f433BFaa128",
			"0x88dA6B6a8D3590e88E0FcadD5CEC56A7C9478319",
			"0x8a718a84ee7B1621E63E680371e0C03C417cCaF6",
			"0xFCb5AFB808b5679b4911230Aa41FfCD0cd335b42",
		},
		fail: []string{
			"0xGHIJK05pwm37asdf5555QWERZCXV2345AoEuIdHt",
			"0xFCb5AFB808b5679b4911230Aa41FfCD0cd335b422222",
			"0xFCb5AFB808b5679b4911230Aa41FfCD0cd33",
			"0b0110100001100101011011000110110001101111",
			"683E07492fBDfDA84457C16546ac3f433BFaa128",
			"1C6o5CDkLxjsVpnLSuqRs1UBFozXLEwYvU",
		},
	}, {
		Name: "FQDN", Func: FQDN,
		pass: []string{
			"domain.com",
			"dom.plato",
			"a.domain.co",
			"foo--bar.com",
			"xn--froschgrn-x9a.com",
			"rebecca.blackfriday",
		},
		fail: []string{
			"abc",
			"256.0.0.0",
			"_.com",
			"*.some.com",
			"s!ome.com",
			"domain.com/",
			"/more.com",
			"domain.com�",
			"domain.com©",
			"example.0",
			"192.168.0.9999",
			"192.168.0",
		},
	}, {
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
	}, {
		Name: "PAN", Func: PAN,
		pass: []string{
			"375556917985515",
			"36050234196908",
			"4716461583322103",
			"4716-2210-5188-5662",
			"4929 7226 5379 7141",
			"5398228707871527",
			"6283875070985593",
			"6263892624162870",
			"6234917882863855",
			"6234698580215388",
			"6226050967750613",
			"6246281879460688",
			"2222155765072228",
			"2225855203075256",
			"2720428011723762",
			"2718760626256570",
			"6765780016990268",
			"4716989580001715211",
		},
		fail: []string{
			"foo",
			"foo",
			"5398228707871528",
			"2718760626256571",
			"2721465526338453",
			"2220175103860763",
			"375556917985515999999993",
			"899999996234917882863855",
			"prefix6234917882863855",
			"623491788middle2863855",
			"6234917882863855suffix",
			"4716989580001715213",
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
