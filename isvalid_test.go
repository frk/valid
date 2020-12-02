package isvalid

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	type Cases []struct {
		args [][]interface{}
		pass []string
		fail []string
	}

	// convenience type, shorter to type
	type args [][]interface{}

	testtable := []struct {
		Name  string
		Func  interface{}
		Cases Cases
	}{{
		Name: "ASCII", Func: ASCII, Cases: Cases{{
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
		}},
	}, {
		Name: "BIC", Func: BIC, Cases: Cases{{
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
		}},
	}, {
		Name: "BTC", Func: BTC, Cases: Cases{{
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
		}},
	}, {
		Name: "Base32", Func: Base32, Cases: Cases{{
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
		}},
	}, {
		Name: "Base58", Func: Base58, Cases: Cases{{
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
		}},
	}, {
		Name: "Base64", Func: Base64, Cases: Cases{{
			args: args{{false}},
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
			args: args{{true}},
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
		}},
	}, {
		Name: "Binary", Func: Binary, Cases: Cases{{
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
		}},
	}, {
		Name: "Bool", Func: Bool, Cases: Cases{{
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
		}},
	}, {
		Name: "CIDR", Func: CIDR, Cases: Cases{{
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
		}},
	}, {
		Name: "CVV", Func: CVV, Cases: Cases{{
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
		}},
	}, {
		Name: "Currency", Func: Currency, Cases: Cases{{
			pass: []string{
				// TODO
			},
			fail: []string{
				// TODO
			},
		}},
	}, {
		Name: "DataURI", Func: DataURI, Cases: Cases{{
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
		}},
	}, {
		Name: "Decimal", Func: Decimal, Cases: Cases{{
			pass: []string{
				// TODO
			},
			fail: []string{
				// TODO
			},
		}},
	}, {
		Name: "Digits", Func: Digits, Cases: Cases{{
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
		}},
	}, {
		Name: "EAN", Func: EAN, Cases: Cases{{
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
		}},
	}, {
		Name: "EIN", Func: EIN, Cases: Cases{{
			pass: []string{
				// TODO
			},
			fail: []string{
				// TODO
			},
		}},
	}, {
		Name: "ETH", Func: ETH, Cases: Cases{{
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
		}},
	}, {
		Name: "FQDN", Func: FQDN, Cases: Cases{{
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
		}},
	}, {
		Name: "Float", Func: Float, Cases: Cases{{
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
		}},
	}, {
		Name: "HSL", Func: HSL, Cases: Cases{{
			pass: []string{
				"hsl(360,0000000000100%,000000100%)",
				"hsl(000010, 00000000001%, 00000040%)",
				"HSL(00000,0000000000100%,000000100%)",
				"hsL(0, 0%, 0%)",
				"hSl(  360  , 100%  , 100%   )",
				"Hsl(  00150  , 000099%  , 01%   )",
				"hsl(01080, 03%, 4%)",
				"hsl(-540, 03%, 4%)",
				"hsla(+540, 03%, 4%)",
				"hsla(+540, 03%, 4%, 500)",
				"hsl(+540deg, 03%, 4%, 500)",
				"hsl(+540gRaD, 03%, 4%, 500)",
				"hsl(+540.01e-98rad, 03%, 4%, 500)",
				"hsl(-540.5turn, 03%, 4%, 500)",
				"hsl(+540, 03%, 4%, 500e-01)",
				"hsl(+540, 03%, 4%, 500e+80)",
				"hsl(4.71239rad, 60%, 70%)",
				"hsl(270deg, 60%, 70%)",
				"hsl(200, +.1%, 62%, 1)",
				"hsl(270 60% 70%)",
				"hsl(200, +.1e-9%, 62e10%, 1)",
				"hsl(.75turn, 60%, 70%)",
				// "hsl(200grad+.1%62%/1)", //supposed to pass, but need to handle delimiters
				"hsl(200grad +.1% 62% / 1)",
				"hsl(270, 60%, 50%, .15)",
				"hsl(270, 60%, 50%, 15%)",
				"hsl(270 60% 50% / .15)",
				"hsl(270 60% 50% / 15%)",
			},
			fail: []string{
				"hsl (360,0000000000100%,000000100%)",
				"hsl(0260, 100 %, 100%)",
				"hsl(0160, 100%, 100%, 100 %)",
				"hsl(-0160, 100%, 100a)",
				"hsl(-0160, 100%, 100)",
				"hsl(-0160 100%, 100%, )",
				"hsl(270 deg, 60%, 70%)",
				"hsl( deg, 60%, 70%)",
				"hsl(, 60%, 70%)",
				"hsl(3000deg, 70%)",
			},
		}},
	}, {
		Name: "Hash", Func: Hash, Cases: Cases{{
			args: args{{"md5"}, {"md4"}, {"ripemd128"}, {"tiger128"}},
			pass: []string{
				"d94f3f016ae679c3008de268209132f2",
				"751adbc511ccbe8edf23d486fa4581cd",
				"88dae00e614d8f24cfd5a8b3f8002e93",
				"0bf1c35032a71a14c2f719e5a14c1e96",
				"d94f3F016Ae679C3008de268209132F2",
				"88DAE00e614d8f24cfd5a8b3f8002E93",
			},
			fail: []string{
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
				"KYT0bf1c35032a71a14c2f719e5a1",
			},
		}, {
			args: args{{"crc32"}, {"crc32b"}},
			pass: []string{
				"d94f3f01",
				"751adbc5",
				"88dae00e",
				"0bf1c350",
				"88DAE00e",
				"751aDBc5",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"q94375dj93458w34",
				"q943",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha1"}, {"tiger160"}, {"ripemd160"}},
			pass: []string{
				"3ca25ae354e192b26879f651a51d92aa8a34d8d3",
				"aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
				"beb8c3f30da46be179b8df5f5ecb5e4b10508230",
				"efd5d3b190e893ed317f38da2420d63b7ae0d5ed",
				"AAF4c61ddCC5e8a2dabede0f3b482cd9AEA9434D",
				"3ca25AE354e192b26879f651A51d92aa8a34d8D3",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha256"}},
			pass: []string{
				"2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
				"1d996e033d612d9af2b44b70061ee0e868bfd14c2dd90b129e1edeb7953e7985",
				"80f70bfeaed5886e33536bcfa8c05c60afef5a0e48f699a7912d5e399cdcc441",
				"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898c",
				"2CF24dba5FB0a30e26E83b2AC5b9E29E1b161e5C1fa7425E73043362938b9824",
				"80F70bFEAed5886e33536bcfa8c05c60aFEF5a0e48f699a7912d5e399cdCC441",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha384"}},
			pass: []string{
				"3fed1f814d28dc5d63e313f8a601ecc4836d1662a19365cbdcf6870f6b56388850b58043f7ebf2418abb8f39c3a42e31",
				"b330f4e575db6e73500bd3b805db1a84b5a034e5d21f0041d91eec85af1dfcb13e40bb1c4d36a72487e048ac6af74b58",
				"bf547c3fc5841a377eb1519c2890344dbab15c40ae4150b4b34443d2212e5b04aa9d58865bf03d8ae27840fef430b891",
				"fc09a3d11368386530f985dacddd026ae1e44e0e297c805c3429d50744e6237eb4417c20ffca8807b071823af13a3f65",
				"3fed1f814d28dc5d63e313f8A601ecc4836d1662a19365CBDCf6870f6b56388850b58043f7ebf2418abb8f39c3a42e31",
				"b330f4E575db6e73500bd3b805db1a84b5a034e5d21f0041d91EEC85af1dfcb13e40bb1c4d36a72487e048ac6af74b58",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha512"}},
			pass: []string{
				"9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
				"83c586381bf5ba94c8d9ba8b6b92beb0997d76c257708742a6c26d1b7cbb9269af92d527419d5b8475f2bb6686d2f92a6649b7f174c1d8306eb335e585ab5049",
				"45bc5fa8cb45ee408c04b6269e9f1e1c17090c5ce26ffeeda2af097735b29953ce547e40ff3ad0d120e5361cc5f9cee35ea91ecd4077f3f589b4d439168f91b9",
				"432ac3d29e4f18c7f604f7c3c96369a6c5c61fc09bf77880548239baffd61636d42ed374f41c261e424d20d98e320e812a6d52865be059745fdb2cb20acff0ab",
				"9B71D224bd62f3785D96d46ad3ea3d73319bFBC2890CAAdae2dff72519673CA72323C3d99ba5c11d7c7ACC6e14b8c5DA0c4663475c2E5c3adef46f73bcDEC043",
				"432AC3d29E4f18c7F604f7c3c96369A6C5c61fC09Bf77880548239baffd61636d42ed374f41c261e424d20d98e320e812a6d52865be059745fdb2cb20acff0ab",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"tiger192"}},
			pass: []string{
				"6281a1f098c5e7290927ed09150d43ff3990a0fe1a48267c",
				"56268f7bc269cf1bc83d3ce42e07a85632394737918f4760",
				"46fc0125a148788a3ac1d649566fc04eb84a746f1a6e4fa7",
				"7731ea1621ae99ea3197b94583d034fdbaa4dce31a67404a",
				"6281A1f098c5e7290927ed09150d43ff3990a0fe1a48267C",
				"46FC0125a148788a3AC1d649566fc04eb84A746f1a6E4fa7",
			},
			fail: []string{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}},
	}, {
		Name: "IBAN", Func: IBAN, Cases: Cases{{
			pass: []string{
				"SC52BAHL01031234567890123456USD",
				"LC14BOSL123456789012345678901234",
				"MT31MALT01100000000000000000123",
				"SV43ACAT00000000000000123123",
				"EG800002000156789012345180002",
				"BE71 0961 2345 6769",
				"FR76 3000 6000 0112 3456 7890 189",
				"DE91 1000 0000 0123 4567 89",
				"GR96 0810 0010 0000 0123 4567 890",
				"RO09 BCYP 0000 0012 3456 7890",
				"SA44 2000 0001 2345 6789 1234",
				"ES79 2100 0813 6101 2345 6789",
				"CH56 0483 5012 3456 7800 9",
				"GB98 MIDL 0700 9312 3456 78",
				"IL170108000000012612345",
				"IT60X0542811101000000123456",
				"JO71CBJO0000000000001234567890",
				"TR320010009999901234567890",
				"BR1500000000000010932840814P2",
				"LB92000700000000123123456123",
				"IR200170000000339545727003",
			},
			fail: []string{
				"XX22YYY1234567890123",
				"FR14 2004 1010 0505 0001 3",
				"FR7630006000011234567890189@",
				"FR7630006000011234567890189😅",
				"FR763000600001123456!!🤨7890189@",
			},
		}},
	}, {
		Name: "Int", Func: Int, Cases: Cases{{
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
		}},
	}, {
		Name: "Uint", Func: Uint, Cases: Cases{{
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
		}},
	}, {
		Name: "Octal", Func: Octal, Cases: Cases{{
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
		}},
	}, {
		Name: "Hex", Func: Hex, Cases: Cases{{
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
		}},
	}, {
		Name: "Numeric", Func: Numeric, Cases: Cases{{
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
		}},
	}, {
		Name: "HexColor", Func: HexColor, Cases: Cases{{
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
		}},
	}, {
		Name: "Email", Func: Email, Cases: Cases{{
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
		}},
	}, {
		Name: "MAC", Func: MAC, Cases: Cases{{
			args: args{{0}},
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
			args: args{{6}},
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
			args: args{{8}},
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
		}},
	}, {
		Name: "PAN", Func: PAN, Cases: Cases{{
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
		}},
	}}

	for _, tab := range testtable {
		t.Run(tab.Name, func(t *testing.T) {
			fn := reflect.ValueOf(tab.Func)

			for _, cases := range tab.Cases {
				argvals := [][]reflect.Value{}
				for _, arglist := range cases.args {
					vals := make([]reflect.Value, len(arglist))
					for i, arg := range arglist {
						vals[i] = reflect.ValueOf(arg)
					}
					argvals = append(argvals, vals)
				}
				if len(argvals) == 0 {
					argvals = append(argvals, []reflect.Value{})
				}

				for _, args := range argvals {
					for _, val := range cases.pass {
						want := true
						t.Run(`"`+val+`"`, func(t *testing.T) {
							rv := reflect.ValueOf(val)
							vv := append([]reflect.Value{rv}, args...)

							got := fn.Call(vv)[0].Bool()
							if got != want {
								t.Errorf("got=%t; want=%t; args=%v", got, want, args)
							}
						})
					}

					for _, val := range cases.fail {
						want := false
						t.Run(`"`+val+`"`, func(t *testing.T) {
							rv := reflect.ValueOf(val)
							vv := append([]reflect.Value{rv}, args...)

							got := fn.Call(vv)[0].Bool()
							if got != want {
								t.Errorf("got=%t; want=%t; args=%v", got, want, args)
							}
						})
					}
				}
			}
		})
	}
}
