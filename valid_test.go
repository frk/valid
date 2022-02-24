package valid

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {

	// for testing the Match validator
	RegisterRegexp(`^(?i)testing$`)

	// convenience types (shorter to type)
	type args [][]interface{}
	type vals []interface{}

	type Cases []struct {
		args args
		pass vals
		fail vals
	}

	testtable := []struct {
		Name  string
		Func  interface{}
		Cases Cases
	}{{
		Name: "ASCII", Func: ASCII, Cases: Cases{{
			pass: vals{
				"foobar",
				"0987654321",
				"test@example.com",
				"1234abcDEF",
			},
			fail: vals{
				"ｆｏｏbar",
				"ｘｙｚ０９８",
				"１２３456",
				"ｶﾀｶﾅ",
			},
		}},
	}, {
		Name: "Alpha", Func: Alpha, Cases: Cases{{
			args: args{{""}},
			pass: vals{
				"abc",
				"ABC",
				"FoObar",
			},
			fail: vals{
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"FÜübar",
				"Jön",
				"Heiß",
			},
		}, {
			args: args{{"aze"}, {"az"}},
			pass: vals{
				"Azərbaycan",
				"Bakı",
				"üöğıəçş",
				"sizAzərbaycanlaşdırılmışlardansınızmı",
				"dahaBirDüzgünString",
				"abcçdeəfgğhxıijkqlmnoöprsştuüvyz",
			},
			fail: vals{
				"rəqəm1",
				"  foo  ",
				"",
				"ab(cd)",
				"simvol@",
				"wəkil",
			},
		}, {
			args: args{{"bul"}, {"bg"}},
			pass: vals{
				"абв",
				"АБВ",
				"жаба",
				"яГоДа",
			},
			fail: vals{
				"abc1",
				"  foo  ",
				"",
				"ЁЧПС",
				"_аз_обичам_обувки_",
				"ехо!",
			},
		}, {
			args: args{{"ces"}, {"cs"}},
			pass: vals{
				"žluťoučký",
				"KŮŇ",
				"Pěl",
				"Ďábelské",
				"ódy",
			},
			fail: vals{
				"ábc1",
				"  fůj  ",
				"",
			},
		}, {
			args: args{{"slk"}, {"sk"}},
			pass: vals{
				"môj",
				"ľúbím",
				"mäkčeň",
				"stĹp",
				"vŕba",
				"ňorimberk",
				"ťava",
				"žanéta",
				"Ďábelské",
				"ódy",
			},
			fail: vals{
				"1moj",
				"你好世界",
				"  Привет мир  ",
				"مرحبا العا ",
			},
		}, {
			args: args{{"dan"}, {"da"}},
			pass: vals{
				"aøå",
				"Ære",
				"Øre",
				"Åre",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"",
			},
		}, {
			args: args{{"nld"}, {"nl"}},
			pass: vals{
				"Kán",
				"één",
				"vóór",
				"nú",
				"héél",
			},
			fail: vals{
				"äca ",
				"abcß",
				"Øre",
			},
		}, {
			args: args{{"deu"}, {"de"}},
			pass: vals{
				"äbc",
				"ÄBC",
				"FöÖbär",
				"Heiß",
			},
			fail: vals{
				"äbc1",
				"  föö  ",
				"",
			},
		}, {
			args: args{{"hun"}, {"hu"}},
			pass: vals{
				"árvíztűrőtükörfúrógép",
				"ÁRVÍZTŰRŐTÜKÖRFÚRÓGÉP",
			},
			fail: vals{
				"äbc1",
				"  fäö  ",
				"Heiß",
				"",
			},
		}, {
			args: args{{"por"}, {"pt"}},
			pass: vals{
				"palíndromo",
				"órgão",
				"qwértyúão",
				"àäãcëüïÄÏÜ",
			},
			fail: vals{
				"12abc",
				"Heiß",
				"Øre",
				"æøå",
				"",
			},
		}, {
			args: args{{"ita"}, {"it"}},
			pass: vals{
				"àéèìîóòù",
				"correnti",
				"DEFINIZIONE",
				"compilazione",
				"metró",
				"pèsca",
				"PÉSCA",
				"genî",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"æøå",
				"",
			},
		}, {
			args: args{{"vie"}, {"vi"}},
			pass: vals{
				"thiến",
				"nghiêng",
				"xin",
				"chào",
				"thế",
				"giới",
			},
			fail: vals{
				"thầy3",
				"Ba gà",
				"",
			},
		}, {
			args: args{{"ara"}, {"ar"}},
			pass: vals{
				"أبت",
				"اَبِتَثّجً",
			},
			fail: vals{
				"١٢٣أبت",
				"١٢٣",
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"FÜübar",
				"Jön",
				"Heiß",
			},
		}, {
			args: args{{"fas"}, {"fa"}},
			pass: vals{
				"پدر",
				"مادر",
				"برادر",
				"خواهر",
			},
			fail: vals{
				"فارسی۱۲۳",
				"۱۶۴",
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"FÜübar",
				"Jön",
				"Heiß",
			},
		}, {
			args: args{{"kur"}, {"ku"}},
			pass: vals{
				"ئؤڤگێ",
				"کوردستان",
			},
			fail: vals{
				"ئؤڤگێ١٢٣",
				"١٢٣",
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"FÜübar",
				"Jön",
				"Heiß",
			},
		}, {
			args: args{{"nob"}, {"nb"}},
			pass: vals{
				"aøå",
				"Ære",
				"Øre",
				"Åre",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"",
			},
		}, {
			args: args{{"pol"}, {"pl"}},
			pass: vals{
				"kreską",
				"zamknięte",
				"zwykłe",
				"kropką",
				"przyjęły",
				"święty",
				"Pozwól",
			},
			fail: vals{
				"12řiď ",
				"blé!!",
				"föö!2!",
			},
		}, {
			args: args{{"srp"}, {"sr"}},
			pass: vals{
				"ШћжЂљЕ",
				"ЧПСТЋЏ",
			},
			fail: vals{
				"řiď ",
				"blé33!!",
				"föö!!",
			},
		}, {
			args: args{{"spa"}, {"es"}},
			pass: vals{
				"ábcó",
				"ÁBCÓ",
				"dormís",
				"volvés",
				"español",
			},
			fail: vals{
				"äca ",
				"abcß",
				"föö!!",
			},
		}, {
			args: args{{"swe"}, {"sv"}},
			pass: vals{
				"religiös",
				"stjäla",
				"västgöte",
				"Åre",
			},
			fail: vals{
				"AİıÖöÇçŞşĞğÜüZ",
				"religiös23",
				"",
			},
		}, {
			args: args{{"tur"}, {"tr"}},
			pass: vals{
				"AİıÖöÇçŞşĞğÜüZ",
			},
			fail: vals{
				"0AİıÖöÇçŞşĞğÜüZ1",
				"  AİıÖöÇçŞşĞğÜüZ  ",
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"Heiß",
			},
		}, {
			args: args{{"ukr"}, {"uk"}},
			pass: vals{
				"АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦШЩЬЮЯ",
			},
			fail: vals{
				"0AİıÖöÇçŞşĞğÜüZ1",
				"  AİıÖöÇçŞşĞğÜüZ  ",
				"abc1",
				"  foo  ",
				"",
				"ÄBC",
				"Heiß",
				"ЫыЪъЭэ",
			},
		}, {
			args: args{{"ell"}, {"el"}},
			pass: vals{
				"αβγδεζηθικλμνξοπρςστυφχψω",
				"ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ",
				"άέήίΰϊϋόύώ",
				"ΆΈΉΊΪΫΎΏ",
			},
			fail: vals{
				"0AİıÖöÇçŞşĞğÜüZ1",
				"  AİıÖöÇçŞşĞğÜüZ  ",
				"ÄBC",
				"Heiß",
				"ЫыЪъЭэ",
				"120",
				"jαckγ",
			},
		}, {
			args: args{{"heb"}, {"he"}},
			pass: vals{
				"בדיקה",
				"שלום",
			},
			fail: vals{
				"בדיקה123",
				"  foo  ",
				"abc1",
				"",
			},
		}, {
			args: args{{"fas"}, {"fa"}},
			pass: vals{
				"تست",
				"عزیزم",
				"ح",
			},
			fail: vals{
				"تست 1",
				"  عزیزم  ",
				"",
			},
		}, {
			args: args{{"tha"}, {"th"}},
			pass: vals{
				"สวัสดี",
				"ยินดีต้อนรับเทสเคส",
			},
			fail: vals{
				"สวัสดีHi",
				"123 ยินดีต้อนรับ",
				"ยินดีต้อนรับ-๑๒๓",
			},
		}},
	}, {
		Name: "Alnum", Func: Alnum, Cases: Cases{{
			args: args{{""}, {"eng"}, {"en"}},
			pass: vals{
				"abc123",
				"ABC11",
			},
			fail: vals{
				"abc ",
				"foo!!",
				"ÄBC",
				"FÜübar",
				"Jön",
			},
		}, {
			args: args{{"aze"}, {"az"}},
			pass: vals{
				"Azərbaycan",
				"Bakı",
				"abc1",
				"abcç2",
				"3kərə4kərə",
			},
			fail: vals{
				"  foo1  ",
				"",
				"ab(cd)",
				"simvol@",
				"wəkil",
			},
		}, {
			args: args{{"bul"}, {"bg"}},
			pass: vals{
				"абв1",
				"4АБ5В6",
				"жаба",
				"яГоДа2",
				"йЮя",
				"123",
			},
			fail: vals{
				" ",
				"789  ",
				"hello000",
			},
		}, {
			args: args{{"ces"}, {"cs"}},
			pass: vals{
				"řiť123",
				"KŮŇ11",
			},
			fail: vals{
				"řiď ",
				"blé!!",
			},
		}, {
			args: args{{"slk"}, {"sk"}},
			pass: vals{
				"1môj",
				"2ľúbím",
				"3mäkčeň",
				"4stĹp",
				"5vŕba",
				"6ňorimberk",
				"7ťava",
				"8žanéta",
				"9Ďábelské",
				"10ódy",
			},
			fail: vals{
				"1moj!",
				"你好世界",
				"  Привет мир  ",
			},
		}, {
			args: args{{"dan"}, {"da"}},
			pass: vals{
				"ÆØÅ123",
				"Ære321",
				"321Øre",
				"123Åre",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"",
			},
		}, {
			args: args{{"nld"}, {"nl"}},
			pass: vals{
				"Kán123",
				"één354",
				"v4óór",
				"nú234",
				"hé54él",
			},
			fail: vals{
				"1äca ",
				"ab3cß",
				"Øre",
			},
		}, {
			args: args{{"deu"}, {"de"}},
			pass: vals{
				"äbc123",
				"ÄBC11",
			},
			fail: vals{
				"äca ",
				"föö!!",
			},
		}, {
			args: args{{"hun"}, {"hu"}},
			pass: vals{
				"0árvíztűrőtükörfúrógép123",
				"0ÁRVÍZTŰRŐTÜKÖRFÚRÓGÉP123",
			},
			fail: vals{
				"1időúr!",
				"äbc1",
				"  fäö  ",
				"Heiß!",
				"",
			},
		}, {
			args: args{{"por"}, {"pt"}},
			pass: vals{
				"palíndromo",
				"2órgão",
				"qwértyúão9",
				"àäãcë4üïÄÏÜ",
			},
			fail: vals{
				"!abc",
				"Heiß",
				"Øre",
				"æøå",
				"",
			},
		}, {
			args: args{{"ita"}, {"it"}},
			pass: vals{
				"123àéèìîóòù",
				"123correnti",
				"DEFINIZIONE321",
				"compil123azione",
				"met23ró",
				"pès56ca",
				"PÉS45CA",
				"gen45î",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"æøå",
				"",
			},
		}, {
			args: args{{"spa"}, {"es"}},
			pass: vals{
				"ábcó123",
				"ÁBCÓ11",
			},
			fail: vals{
				"äca ",
				"abcß",
				"föö!!",
			},
		}, {
			args: args{{"vie"}, {"vi"}},
			pass: vals{
				"Thầy3",
				"3Gà",
			},
			fail: vals{
				"toang!",
				"Cậu Vàng",
			},
		}, {
			args: args{{"ara"}, {"ar"}},
			pass: vals{
				"أبت123",
				"أبتَُِ١٢٣",
			},
			fail: vals{
				"äca ",
				"abcß",
				"föö!!",
			},
		}, {
			args: args{{"fas"}, {"fa"}},
			pass: vals{
				"پارسی۱۲۳",
				"۱۴۵۶",
			},
			fail: vals{
				"äca ",
				"abcßة",
				"föö!!",
				"٤٥٦",
				"مژگان9",
			},
		}, {
			args: args{{"kur"}, {"ku"}},
			pass: vals{
				"ئؤڤگێ١٢٣",
			},
			fail: vals{
				"äca ",
				"abcß",
				"föö!!",
			},
		}, {
			args: args{{"nob"}, {"nb"}},
			pass: vals{
				"ÆØÅ123",
				"Ære321",
				"321Øre",
				"123Åre",
			},
			fail: vals{
				"äbc123",
				"ÄBC11",
				"",
			},
		}, {
			args: args{{"pol"}, {"pl"}},
			pass: vals{
				"kre123ską",
				"zam21knięte",
				"zw23ykłe",
				"123",
				"prz23yjęły",
				"świ23ęty",
				"Poz1322wól",
			},
			fail: vals{
				"12řiď ",
				"blé!!",
				"föö!2!",
			},
		}, {
			args: args{{"srp"}, {"sr"}},
			pass: vals{
				"ШћжЂљЕ123",
				"ЧПСТ132ЋЏ",
			},
			fail: vals{
				"řiď ",
				"blé!!",
				"föö!!",
			},
		}, {
			args: args{{"swe"}, {"sv"}},
			pass: vals{
				"religiös13",
				"st23jäla",
				"västgöte123",
				"123Åre",
			},
			fail: vals{
				"AİıÖöÇçŞşĞğÜüZ",
				"foo!!",
				"",
			},
		}, {
			args: args{{"tur"}, {"tr"}},
			pass: vals{
				"AİıÖöÇçŞşĞğÜüZ123",
			},
			fail: vals{
				"AİıÖöÇçŞşĞğÜüZ ",
				"foo!!",
				"ÄBC",
			},
		}, {
			args: args{{"ukr"}, {"uk"}},
			pass: vals{
				"АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦШЩЬЮЯ123",
			},
			fail: vals{
				"éeoc ",
				"foo!!",
				"ÄBC",
				"ЫыЪъЭэ",
			},
		}, {
			args: args{{"ell"}, {"el"}},
			pass: vals{
				"αβγδεζηθικλμνξοπρςστυφχψω",
				"ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ",
				"20θ",
				"1234568960",
			},
			fail: vals{
				"0AİıÖöÇçŞşĞğÜüZ1",
				"  AİıÖöÇçŞşĞğÜüZ  ",
				"ÄBC",
				"Heiß",
				"ЫыЪъЭэ",
				"jαckγ",
			},
		}, {
			args: args{{"heb"}, {"he"}},
			pass: vals{
				"אבג123",
				"שלום11",
			},
			fail: vals{
				"אבג ",
				"לא!!",
				"abc",
				"  foo  ",
			},
		}, {
			args: args{{"tha"}, {"th"}},
			pass: vals{
				"สวัสดี๑๒๓",
				"ยินดีต้อนรับทั้ง๒คน",
			},
			fail: vals{
				"1.สวัสดี",
				"ยินดีต้อนรับทั้ง 2 คน",
			},
		}},
	}, {
		Name: "BIC", Func: BIC, Cases: Cases{{
			pass: vals{
				"SBICKEN1345",
				"SBICKEN1",
				"SBICKENY",
				"SBICKEN1YYP",
			},

			fail: vals{
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
			pass: vals{
				"1MUz4VMYui5qY1mxUiG8BQ1Luv6tqkvaiL",
				"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
				"bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
			},

			fail: vals{
				"4J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
				"0x56F0B8A998425c53c75C4A303D4eF987533c5597",
				"pp8skudq3x5hzw8ew7vzsw8tn4k8wxsqsv0lt0mf3g",
			},
		}},
	}, {
		Name: "Base32", Func: Base32, Cases: Cases{{
			pass: vals{
				"ZG======",
				"JBSQ====",
				"JBSWY===",
				"JBSWY3A=",
				"JBSWY3DP",
				"JBSWY3DPEA======",
				"K5SWYY3PNVSSA5DPEBXG6ZA=",
				"K5SWYY3PNVSSA5DPEBXG6===",
			},
			fail: vals{
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
			pass: vals{
				"BukQL",
				"3KMUV89zab",
				"91GHkLMNtyo98",
				"YyjKm3H",
				"Mkhss145TRFg",
				"7678765677",
				"abcodpq",
				"AAVHJKLPY",
			},
			fail: vals{
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
			pass: vals{
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
			fail: vals{
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
			pass: vals{
				"",
				"bGFkaWVzIGFuZCBnZW50bGVtZW4sIHdlIGFyZSBmbG9hdGluZyBpbiBzcGFjZQ",
				"1234",
				"bXVtLW5ldmVyLXByb3Vk",
				"PDw_Pz8-Pg",
				"VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw",
			},
			fail: vals{
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
			pass: vals{
				"0",
				"1",
				"0001",
				"1110",
				"0001011010101",
				"0b010",
				"0B010",
			},
			fail: vals{
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
			pass: vals{
				"true",
				"false",
				"TRUE",
				"FALSE",
			},
			fail: vals{
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
			pass: vals{
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
			fail: vals{
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
			pass: vals{
				"123",
				"1234",
			},
			fail: vals{
				"",
				"12",
				"12345",
				"abc",
				"abcd",
			},
		}},
	}, {
		Name: "Currency", Func: Currency, Cases: Cases{{
			args: args{{"gbp", (*CurrencyOpts)(nil)}},
			pass: vals{
				"£1",
				"1234·05",
				"1,234·05",
				"1 234·05",
				"£1234·05",
				"£1,234·05",
				"£1 234·05",
				"£\u20091234·05",
				"£\u20091,234·05",
				"£\u20091 234·05",
				"£\u202F1234·05",
				"£\u202F1,234·05",
				"£\u202F1 234·05",
				"£2.99",
			},
			fail: vals{
				"",
				"$1",
				"£1,00",
				"£ 1,234·07",
				"1,234·05£",
				"£1,234·0",
				"£1,234·999",
				"£1,234·9a",
				"£1,2E4·99",
				"£2.",
				"£.99",
				"£02.99",
				"£1,234 567,890.99",
				"£1,23,456,789.99",
				"£1,234,,567.99",
				"£,234.99",
			},
		}},
	}, {
		Name: "DataURI", Func: DataURI, Cases: Cases{{
			pass: vals{
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
			fail: vals{
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
			args: args{{"en"}},
			pass: vals{
				"123",
				"0",
				"-0",
				"-00123",
				"00123",
				"+123",
				"0.01",
				".1",
				"1.0",
				"-.25",
				"-0",
				"0.0000000000001",
			},
			fail: vals{
				"0,01",
				",1",
				"1,0",
				"-,25",
				"0,0000000000001",
				"0٫01",
				"٫1",
				"1٫0",
				"-٫25",
				"0٫0000000000001",
				"....",
				" ",
				"",
				"-",
				"+",
				".",
				"0.1a",
				"a",
				"\n",
			},
		}, {
			args: args{{"en_AU"}},
			pass: vals{
				"123",
				"00123",
				"-00123",
				"0",
				"-0",
				"+123",
				"0.01",
				".1",
				"1.0",
				"-.25",
				"-0",
				"0.0000000000001",
			},
			fail: vals{
				"0,01",
				",1",
				"1,0",
				"-,25",
				"0,0000000000001",
				"0٫01",
				"٫1",
				"1٫0",
				"-٫25",
				"0٫0000000000001",
				"....",
				" ",
				"",
				"-",
				"+",
				".",
				"0.1a",
				"a",
				"\n",
			},
		}, {
			args: args{{"bg_BG"}},
			pass: vals{
				"123",
				"00123",
				"-00123",
				"0",
				"-0",
				"+123",
				"0,01",
				",1",
				"1,0",
				"-,25",
				"-0",
				"0,0000000000001",
			},
			fail: vals{
				"0.0000000000001",
				"0.01",
				".1",
				"1.0",
				"-.25",
				"0٫01",
				"٫1",
				"1٫0",
				"-٫25",
				"0٫0000000000001",
				"....",
				" ",
				"",
				"-",
				"+",
				".",
				"0.1a",
				"a",
				"\n",
			},
		}, {
			args: args{{"cs_CZ"}},
			pass: vals{
				"123",
				"00123",
				"-00123",
				"0",
				"-0",
				"+123",
				"0,01",
				",1",
				"1,0",
				"-,25",
				"-0",
				"0,0000000000001",
			},
			fail: vals{
				"0.0000000000001",
				"0.01",
				".1",
				"1.0",
				"-.25",
				"0٫01",
				"٫1",
				"1٫0",
				"-٫25",
				"0٫0000000000001",
				"....",
				" ",
				"",
				"-",
				"+",
				".",
				"0.1a",
				"a",
				"\n",
			},
		}, {
			args: args{{"ar_JO"}},
			pass: vals{
				"٠٫٠١",
			},
			fail: vals{
				"0,01",
			},
		}, {
			args: args{{"ar_EG"}},
			pass: vals{
				"٠٫٠١",
			},
			fail: vals{
				"0,01",
			},
		}, {
			args: args{{"en_ZM"}},
			pass: vals{
				"0.01",
			},
			fail: vals{
				"0,01",
			},
		}},
	}, {
		Name: "Digits", Func: Digits, Cases: Cases{{
			pass: vals{
				"123",
				"00123",
				"0",
				"0987654321",
			},
			fail: vals{
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
			pass: vals{
				"9421023610112",
				"1234567890128",
				"4012345678901",
				"9771234567003",
				"9783161484100",
				"73513537",
			},
			fail: vals{
				"5901234123451",
				"079777681629",
				"0705632085948",
			},
		}},
	}, {
		Name: "EIN", Func: todo_EIN, Cases: Cases{{
			pass: vals{
				// TODO
			},
			fail: vals{
				// TODO
			},
		}},
	}, {
		Name: "ETH", Func: ETH, Cases: Cases{{
			pass: vals{
				"0x0000000000000000000000000000000000000001",
				"0x683E07492fBDfDA84457C16546ac3f433BFaa128",
				"0x88dA6B6a8D3590e88E0FcadD5CEC56A7C9478319",
				"0x8a718a84ee7B1621E63E680371e0C03C417cCaF6",
				"0xFCb5AFB808b5679b4911230Aa41FfCD0cd335b42",
			},
			fail: vals{
				"0xGHIJK05pwm37asdf5555QWERZCXV2345AoEuIdHt",
				"0xFCb5AFB808b5679b4911230Aa41FfCD0cd335b422222",
				"0xFCb5AFB808b5679b4911230Aa41FfCD0cd33",
				"0b0110100001100101011011000110110001101111",
				"683E07492fBDfDA84457C16546ac3f433BFaa128",
				"1C6o5CDkLxjsVpnLSuqRs1UBFozXLEwYvU",
			},
		}},
	}, {
		Name: "Email", Func: Email, Cases: Cases{{
			pass: vals{
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
			fail: vals{
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
		Name: "FQDN", Func: FQDN, Cases: Cases{{
			pass: vals{
				"domain.com",
				"dom.plato",
				"a.domain.co",
				"foo--bar.com",
				"xn--froschgrn-x9a.com",
				"rebecca.blackfriday",
			},
			fail: vals{
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
			pass: vals{
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
			fail: vals{
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
			pass: vals{
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
			fail: vals{
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
			pass: vals{
				"d94f3f016ae679c3008de268209132f2",
				"751adbc511ccbe8edf23d486fa4581cd",
				"88dae00e614d8f24cfd5a8b3f8002e93",
				"0bf1c35032a71a14c2f719e5a14c1e96",
				"d94f3F016Ae679C3008de268209132F2",
				"88DAE00e614d8f24cfd5a8b3f8002E93",
			},
			fail: vals{
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
				"KYT0bf1c35032a71a14c2f719e5a1",
			},
		}, {
			args: args{{"crc32"}, {"crc32b"}},
			pass: vals{
				"d94f3f01",
				"751adbc5",
				"88dae00e",
				"0bf1c350",
				"88DAE00e",
				"751aDBc5",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"q94375dj93458w34",
				"q943",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha1"}, {"tiger160"}, {"ripemd160"}},
			pass: vals{
				"3ca25ae354e192b26879f651a51d92aa8a34d8d3",
				"aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
				"beb8c3f30da46be179b8df5f5ecb5e4b10508230",
				"efd5d3b190e893ed317f38da2420d63b7ae0d5ed",
				"AAF4c61ddCC5e8a2dabede0f3b482cd9AEA9434D",
				"3ca25AE354e192b26879f651A51d92aa8a34d8D3",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha256"}},
			pass: vals{
				"2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
				"1d996e033d612d9af2b44b70061ee0e868bfd14c2dd90b129e1edeb7953e7985",
				"80f70bfeaed5886e33536bcfa8c05c60afef5a0e48f699a7912d5e399cdcc441",
				"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898c",
				"2CF24dba5FB0a30e26E83b2AC5b9E29E1b161e5C1fa7425E73043362938b9824",
				"80F70bFEAed5886e33536bcfa8c05c60aFEF5a0e48f699a7912d5e399cdCC441",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha384"}},
			pass: vals{
				"3fed1f814d28dc5d63e313f8a601ecc4836d1662a19365cbdcf6870f6b56388850b58043f7ebf2418abb8f39c3a42e31",
				"b330f4e575db6e73500bd3b805db1a84b5a034e5d21f0041d91eec85af1dfcb13e40bb1c4d36a72487e048ac6af74b58",
				"bf547c3fc5841a377eb1519c2890344dbab15c40ae4150b4b34443d2212e5b04aa9d58865bf03d8ae27840fef430b891",
				"fc09a3d11368386530f985dacddd026ae1e44e0e297c805c3429d50744e6237eb4417c20ffca8807b071823af13a3f65",
				"3fed1f814d28dc5d63e313f8A601ecc4836d1662a19365CBDCf6870f6b56388850b58043f7ebf2418abb8f39c3a42e31",
				"b330f4E575db6e73500bd3b805db1a84b5a034e5d21f0041d91EEC85af1dfcb13e40bb1c4d36a72487e048ac6af74b58",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"sha512"}},
			pass: vals{
				"9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
				"83c586381bf5ba94c8d9ba8b6b92beb0997d76c257708742a6c26d1b7cbb9269af92d527419d5b8475f2bb6686d2f92a6649b7f174c1d8306eb335e585ab5049",
				"45bc5fa8cb45ee408c04b6269e9f1e1c17090c5ce26ffeeda2af097735b29953ce547e40ff3ad0d120e5361cc5f9cee35ea91ecd4077f3f589b4d439168f91b9",
				"432ac3d29e4f18c7f604f7c3c96369a6c5c61fc09bf77880548239baffd61636d42ed374f41c261e424d20d98e320e812a6d52865be059745fdb2cb20acff0ab",
				"9B71D224bd62f3785D96d46ad3ea3d73319bFBC2890CAAdae2dff72519673CA72323C3d99ba5c11d7c7ACC6e14b8c5DA0c4663475c2E5c3adef46f73bcDEC043",
				"432AC3d29E4f18c7F604f7c3c96369A6C5c61fC09Bf77880548239baffd61636d42ed374f41c261e424d20d98e320e812a6d52865be059745fdb2cb20acff0ab",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}, {
			args: args{{"tiger192"}},
			pass: vals{
				"6281a1f098c5e7290927ed09150d43ff3990a0fe1a48267c",
				"56268f7bc269cf1bc83d3ce42e07a85632394737918f4760",
				"46fc0125a148788a3ac1d649566fc04eb84a746f1a6e4fa7",
				"7731ea1621ae99ea3197b94583d034fdbaa4dce31a67404a",
				"6281A1f098c5e7290927ed09150d43ff3990a0fe1a48267C",
				"46FC0125a148788a3AC1d649566fc04eb84A746f1a6E4fa7",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"KYT0bf1c35032a71a14c2f719e5a14c1dsjkjkjkjkkjk",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}},
	}, {
		Name: "Hex", Func: Hex, Cases: Cases{{
			pass: vals{
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
			fail: vals{
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
		Name: "HexColor", Func: HexColor, Cases: Cases{{
			pass: vals{
				"#ff0000ff",
				"#ff0034",
				"#CCCCCC",
				"0f38",
				"fff",
				"#f00",
			},
			fail: vals{
				"#ff",
				"fff0a",
				"#ff12FG",
			},
		}},
	}, {
		Name: "IBAN", Func: IBAN, Cases: Cases{{
			pass: vals{
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
			fail: vals{
				"XX22YYY1234567890123",
				"FR14 2004 1010 0505 0001 3",
				"FR7630006000011234567890189@",
				"FR7630006000011234567890189😅",
				"FR763000600001123456!!🤨7890189@",
			},
		}},
	}, {
		Name: "IMEI", Func: IMEI, Cases: Cases{{
			pass: vals{
				"352099001761481",
				"868932036356090",
				"490154203237518",
				"546918475942169",
				"998227667144730",
				"532729766805999",
				"35-209900-176148-1",
				"86-893203-635609-0",
				"49-015420-323751-8",
				"54-691847-594216-9",
				"99-822766-714473-0",
				"53-272976-680599-9",
			},
			fail: vals{
				"490154203237517",
				"3568680000414120",
				"3520990017614823",
				"49-015420-323751-7",
				"35-686800-0041412-0",
				"35-209900-1761482-3",
			},
		}},
	}, {
		Name: "IP", Func: IP, Cases: Cases{{
			args: args{{0}},
			pass: vals{
				"127.0.0.1",
				"0.0.0.0",
				"255.255.255.255",
				"1.2.3.4",
				"::1",
				"2001:db8:0000:1:1:1:1:1",
				"2001:41d0:2:a141::1",
				"::ffff:127.0.0.1",
				"::0000",
				"0000::",
				"1::",
				"1111:1:1:1:1:1:1:1",
				"fe80::a6db:30ff:fe98:e946",
				"::",
				"::ffff:127.0.0.1",
				"0:0:0:0:0:ffff:127.0.0.1",
			},
			fail: vals{
				"abc",
				"256.0.0.0",
				"0.0.0.256",
				"26.0.0.256",
				"0200.200.200.200",
				"200.0200.200.200",
				"200.200.0200.200",
				"200.200.200.0200",
				"::banana",
				"banana::",
				"::1banana",
				"::1::",
				"1:",
				":1",
				":1:1:1::2",
				"1:1:1:1:1:1:1:1:1:1:1:1:1:1:1:1",
				"::11111",
				"11111:1:1:1:1:1:1:1",
				"2001:db8:0000:1:1:1:1::1",
				"0:0:0:0:0:0:ffff:127.0.0.1",
				"0:0:0:0:ffff:127.0.0.1",
			},
		}, {
			args: args{{4}},
			pass: vals{
				"127.0.0.1",
				"0.0.0.0",
				"255.255.255.255",
				"1.2.3.4",
				"255.0.0.1",
				"0.0.1.1",
			},
			fail: vals{
				"::1",
				"2001:db8:0000:1:1:1:1:1",
				"::ffff:127.0.0.1",
				"137.132.10.01",
				"0.256.0.256",
				"255.256.255.256",
			},
		}, {
			args: args{{6}},
			pass: vals{
				"::1",
				"2001:db8:0000:1:1:1:1:1",
				"::ffff:127.0.0.1",
				"fe80::1234%1",
				"ff08::9abc%10",
				"ff08::9abc%interface10",
				"ff02::5678%pvc1.3",
			},
			fail: vals{
				"127.0.0.1",
				"0.0.0.0",
				"fe80:::a6db:30ff:fe98:e946",
				"255.255.255.255",
				"1.2.3.4",
				"::ffff:287.0.0.1",
				"2001::kexp:1:1:1:1",
				"%",
				"fe80::1234%",
				"fe80::1234%1%3%4",
				"fe80%fe80%",
			},
		}},
	}, {
		Name: "IPRange", Func: IPRange, Cases: Cases{{
			pass: vals{
				"127.0.0.1/24",
				"0.0.0.0/0",
				"255.255.255.0/32",
			},
			fail: vals{
				"127.200.230.1/35",
				"127.200.230.1/-1",
				"1.1.1.1/011",
				"::1/64",
				"1.1.1/24.1",
				"1.1.1.1/01",
				"1.1.1.1/1.1",
				"1.1.1.1/1.",
				"1.1.1.1/1/1",
				"1.1.1.1",
			},
		}},
	}, {
		Name: "ISBN", Func: ISBN, Cases: Cases{{
			args: args{{10}},
			pass: vals{
				"3836221195", "3-8362-2119-5", "3 8362 2119 5",
				"1617290858", "1-61729-085-8", "1 61729 085-8",
				"0007269706", "0-00-726970-6", "0 00 726970 6",
				"3423214120", "3-423-21412-0", "3 423 21412 0",
				"340101319X", "3-401-01319-X", "3 401 01319 X",
			},
			fail: vals{
				"3423214121", "3-423-21412-1", "3 423 21412 1",
				"978-3836221191", "9783836221191",
				"123456789a", "foo", "",
			},
		}, {
			args: args{{13}},
			pass: vals{
				"9783836221191", "978-3-8362-2119-1", "978 3 8362 2119 1",
				"9783401013190", "978-3401013190", "978 3401013190",
				"9784873113685", "978-4-87311-368-5", "978 4 87311 368 5",
			},
			fail: vals{
				"9783836221190", "978-3-8362-2119-0", "978 3 8362 2119 0",
				"3836221195", "3-8362-2119-5", "3 8362 2119 5",
				"01234567890ab", "foo", "",
			},
		}, {
			args: args{{0}},
			pass: vals{
				"340101319X",
				"9784873113685",
			},
			fail: vals{
				"3423214121",
				"9783836221190",
			},
		}},
	}, {
		Name: "ISIN", Func: ISIN, Cases: Cases{{
			pass: vals{
				"AU0000XVGZA3",
				"DE000BAY0017",
				"BE0003796134",
				"SG1G55870362",
				"GB0001411924",
				"DE000WCH8881",
				"PLLWBGD00016",
			},
			fail: vals{
				"DE000BAY0018",
				"PLLWBGD00019",
				"foo",
				"5398228707871528",
			},
		}},
	}, {
		Name: "ISO31661A", Func: ISO31661A, Cases: Cases{{
			args: args{{0}},
			pass: vals{
				"FR",
				"fR",
				"GB",
				"PT",
				"CM",
				"JP",
				"ABW",
				"HND",
				"KHM",
				"RWA",
				"PM",
				"ZW",
				"MM",
				"cc",
				"GG",
			},
			fail: vals{
				"",
				"AA",
				"PI",
				"RP",
				"WV",
				"WL",
				"UK",
				"ZZ",
			},
		}, {
			args: args{{2}},
			pass: vals{
				"FR",
				"fR",
				"GB",
				"PT",
				"CM",
				"JP",
				"PM",
				"ZW",
				"MM",
				"cc",
				"GG",
			},
			fail: vals{
				"",
				"FRA",
				"AA",
				"PI",
				"RP",
				"WV",
				"WL",
				"UK",
				"ZZ",
			},
		}, {
			args: args{{3}},
			pass: vals{
				"ABW",
				"HND",
				"KHM",
				"RWA",
			},
			fail: vals{
				"",
				"FR",
				"fR",
				"GB",
				"PT",
				"CM",
				"JP",
				"PM",
				"ZW",
			},
		}},
	}, {
		Name: "ISO4217", Func: ISO4217, Cases: Cases{{
			pass: vals{
				// TODO
			},
			fail: vals{
				// TODO
			},
		}},
	}, {
		Name: "ISRC", Func: ISRC, Cases: Cases{{
			pass: vals{
				"USAT29900609",
				"GBAYE6800011",
				"USRC15705223",
				"USCA29500702",
			},
			fail: vals{
				"USAT2990060",
				"SRC15705223",
				"US-CA29500702",
				"USARC15705223",
			},
		}},
	}, {
		Name: "ISSN", Func: ISSN, Cases: Cases{{
			args: args{{false, false}},
			pass: vals{
				"0378-5955",
				"0000-0000",
				"2434-561X",
				"2434-561x",
				"01896016",
				"20905076",
			},
			fail: vals{
				"0378-5954",
				"0000-0001",
				"0378-123",
				"037-1234",
				"0",
				"2434-561c",
				"1684-5370",
				"19960791",
				"",
			},
		}, {
			args: args{{false, true}},
			pass: vals{
				"2434-561X",
				"2434561X",
				"0378-5955",
				"03785955",
			},
			fail: vals{
				"2434-561x",
				"2434561x",
			},
		}, {
			args: args{{true, false}},
			pass: vals{
				"2434-561X",
				"2434-561x",
				"0378-5955",
			},
			fail: vals{
				"2434561X",
				"2434561x",
				"03785955",
			},
		}, {
			args: args{{true, true}},
			pass: vals{
				"2434-561X",
				"0378-5955",
			},
			fail: vals{
				"2434-561x",
				"2434561X",
				"2434561x",
				"03785955",
			},
		}},
	}, {
		Name: "In", Func: In, Cases: Cases{{
			args: args{},
			pass: vals{},
			fail: vals{"foo"},
		}, {
			args: args{{true, 123, "bar", 23.0, "foo"}},
			pass: vals{
				"foo",
				"bar",
			},
			fail: vals{
				"baz",
				"123",
			},
		}},
	}, {
		Name: "Int", Func: Int, Cases: Cases{{
			pass: vals{
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
			fail: vals{
				"",
				"0.1",
				".01",
				"123e45",
				"abcdef",
				"      ",
			},
		}},
	}, {
		Name: "JSON", Func: JSON, Cases: Cases{{
			pass: vals{
				[]byte(`{ "key": "value" }`),
				[]byte(`{}`),
				[]byte(`[123, "hey"]`),
				[]byte(`[]`),
				[]byte(`null`),
				[]byte(`1234`),
				[]byte(`"yes"`),
				[]byte(`true`),
				[]byte(`false`),
			},
			fail: vals{
				[]byte(`{ key: "value" }`),
				[]byte(`{ \'key\': \'value\' }`),
				[]byte(`nope`),
				[]byte(`nil`),
				[]byte(`{ "key": nil }`),
			},
		}},
	}, {
		Name: "JWT", Func: JWT, Cases: Cases{{
			pass: vals{
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRJbkFzIjoiYWRtaW4iLCJpYXQiOjE0MjI3Nzk2Mzh9.gzSraSYS8EXBxLN_oWnFSRgCzcmJmMjLiuyu5CSpyHI",
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb3JlbSI6Imlwc3VtIn0.ymiJSsMJXR6tMSr8G9usjQ15_8hKPDv_CArLhxw28MI",
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb2xvciI6InNpdCIsImFtZXQiOlsibG9yZW0iLCJpcHN1bSJdfQ.rRpe04zbWbbJjwM43VnHzAboDzszJtGrNsUxaqQ-GQ8",
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqb2huIjp7ImFnZSI6MjUsImhlaWdodCI6MTg1fSwiamFrZSI6eyJhZ2UiOjMwLCJoZWlnaHQiOjI3MH19.YRLPARDmhGMC3BBk_OhtwwK21PIkVCqQe8ncIRPKo-E",
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ", // No signature

			},
			fail: vals{
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				"$Zs.ewu.su84",
				"ks64$S/9.dy$§kz.3sd73b",
			},
		}},
	}, {
		Name: "LatLong", Func: LatLong, Cases: Cases{{
			args: args{{false}},
			pass: vals{
				"(-17.738223, 85.605469)",
				"(-12.3456789, +12.3456789)",
				"(-60.978437, -0.175781)",
				"(77.719772, -37.529297)",
				"(7.264394, 165.058594)",
				"0.955766, -19.863281",
				"(31.269161,164.355469)",
				"+12.3456789, -12.3456789",
				"-15.379543, -137.285156",
				"(11.770570, -162.949219)",
				"-55.034319, 113.027344",
				"58.025555, 36.738281",
				"55.720923,-28.652344",
				"-90.00000,-180.00000",
				"(-71, -146)",
				"(-71.616864, -146.616864)",
				"-0.55, +0.22",
				"90, 180",
				"+90, -180",
				"-90,+180",
				"90,180",
				"0, 0",
			},
			fail: vals{
				"(020.000000, 010.000000000)",
				"89.9999999989, 360.0000000",
				"90.1000000, 180.000000",
				"+90.000000, -180.00001",
				"090.0000, 0180.0000",
				"126, -158",
				"(-126.400010, -158.400010)",
				"-95, -96",
				"-95.738043, -96.738043",
				"137, -148",
				"(-137.5942, -148.5942)",
				"(-120, -203)",
				"(-119, -196)",
				"+119.821728, -196.821728",
				"(-110, -223)",
				"-110.369532, 223.369532",
				"(-120.969949, +203.969949)",
				"-116, -126",
				"-116.894222, -126.894222",
				"-112, -160",
				"-112.96381, -160.96381",
				"-90., -180.",
				"+90.1, -180.1",
				"(-17.738223, 85.605469",
				"0.955766, -19.863281)",
				"+,-",
				"(,)",
				",",
				" ",
			},
		}, {
			args: args{{true}},
			pass: vals{
				"40° 26′ 46″ N, 79° 58′ 56″ W",
				"40° 26′ 46″ S, 79° 58′ 56″ E",
				"90° 0′ 0″ S, 180° 0′ 0″ E",
				"40° 26′ 45.9996″ N, 79° 58′ 55.2″ E",
				"40° 26′ 46″ n, 79° 58′ 56″ w",
				"40°26′46″s, 79°58′56″e",
				"11° 0′ 0.005″ S, 180° 0′ 0″ E",
				"40°26′45.9996″N, 79°58′55.2″E",
			},
			fail: vals{
				"100° 26′ 46″ N, 79° 70′ 56″ W",
				"40° 89′ 46″ S, 79° 58′ 100″ E",
				"40° 26.445′ 45″ N, 79° 58′ 55.2″ E",
				"40° 46″ N, 79° 58′ 56″ W",
			},
		}},
	}, {
		Name: "Locale", Func: Locale, Cases: Cases{{
			pass: vals{
				"uz_Latn_UZ",
				"en",
				"gsw",
				"es_ES",
				"sw_KE",
				"am_ET",
				"ca_ES_VALENCIA",
				"en_US_POSIX",
			},
			fail: vals{
				"lo_POP",
				"12",
				"12_DD",
			},
		}},
	}, {
		Name: "LowerCase", Func: LowerCase, Cases: Cases{{
			pass: vals{
				"abc",
				"abc123",
				"this is lowercase.",
				"tr竪s 端ber",
			},
			fail: vals{
				"fooBar",
				"123A",
			},
		}},
	}, {
		Name: "MAC", Func: MAC, Cases: Cases{{
			args: args{{0}},
			pass: vals{
				"08:00:2b:01:02:03",
				"08-00-2b-01-02-03",
				"01:AB:03:04:05:06",
				"01-02-03-04-05-ab",
				"08002b010203",
				"08:00:2b:01:02:03:04:05",
				"08-00-2b-01-02-03-04-05",
				"08002b0102030405",
			},
			fail: vals{
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
			pass: vals{
				"08:00:2b:01:02:03",
				"08-00-2b-01-02-03",
				"01:AB:03:04:05:06",
				"01-02-03-04-05-ab",
				"08002b010203",
			},
			fail: vals{
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
			pass: vals{
				"08:00:2b:01:02:03:04:05",
				"08-00-2b-01-02-03-04-05",
				"08002b0102030405",
			},
			fail: vals{
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
		Name: "MIME", Func: MIME, Cases: Cases{{
			pass: vals{
				"application/json",
				"application/xhtml+xml",
				"audio/mp4",
				"image/bmp",
				"font/woff2",
				"message/http",
				"model/vnd.gtw",
				"multipart/form-data",
				"multipart/form-data; boundary=something",
				"multipart/form-data; charset=utf-8; boundary=something",
				"multipart/form-data; boundary=something; charset=utf-8",
				"multipart/form-data; boundary=something; charset=\"utf-8\"",
				"multipart/form-data; boundary=\"something\"; charset=utf-8",
				"multipart/form-data; boundary=\"something\"; charset=\"utf-8\"",
				"text/css",
				"text/plain; charset=utf8",
				"Text/HTML;Charset=\"utf-8\"",
				"text/html;charset=UTF-8",
				"Text/html;charset=UTF-8",
				"text/html; charset=us-ascii",
				"text/html; charset=us-ascii (Plain text)",
				"text/html; charset=\"us-ascii\"",
				"video/mp4",
			},
			fail: vals{
				"",
				" ",
				"/",
				"f/b",
				"application",
				"application\\json",
				"application/json/text",
				"application/json; charset=utf-8",
				"audio/mp4; charset=utf-8",
				"image/bmp; charset=utf-8",
				"font/woff2; charset=utf-8",
				"message/http; charset=utf-8",
				"model/vnd.gtw; charset=utf-8",
				"video/mp4; charset=utf-8",
			},
		}},
	}, {
		Name: "MD5", Func: MD5, Cases: Cases{{
			pass: vals{
				"d94f3f016ae679c3008de268209132f2",
				"751adbc511ccbe8edf23d486fa4581cd",
				"88dae00e614d8f24cfd5a8b3f8002e93",
				"0bf1c35032a71a14c2f719e5a14c1e96",
			},
			fail: vals{
				"KYT0bf1c35032a71a14c2f719e5a14c1",
				"q94375dj93458w34",
				"39485729348",
				"%&FHKJFvk",
			},
		}},
	}, {
		Name: "MagnetURI", Func: MagnetURI, Cases: Cases{{
			pass: vals{
				"magnet:?xt=urn:btih:06E2A9683BF4DA92C73A661AC56F0ECC9C63C5B4&dn=helloword2000&tr=udp://helloworld:1337/announce",
				"magnet:?xt=urn:btih:3E30322D5BFC7444B7B1D8DD42404B75D0531DFB&dn=world&tr=udp://world.com:1337",
				"magnet:?xt=urn:btih:4ODKSDJBVMSDSNJVBCBFYFBKNRU875DW8D97DWC6&dn=helloworld&tr=udp://helloworld.com:1337",
				"magnet:?xt=urn:btih:1GSHJVBDVDVJFYEHKFHEFIO8573898434JBFEGHD&dn=foo&tr=udp://foo.com:1337",
				"magnet:?xt=urn:btih:MCJDCYUFHEUD6E2752T7UJNEKHSUGEJFGTFHVBJS&dn=bar&tr=udp://bar.com:1337",
				"magnet:?xt=urn:btih:LAKDHWDHEBFRFVUFJENBYYTEUY837562JH2GEFYH&dn=foobar&tr=udp://foobar.com:1337",
				"magnet:?xt=urn:btih:MKCJBHCBJDCU725TGEB3Y6RE8EJ2U267UNJFGUID&dn=test&tr=udp://test.com:1337",
				"magnet:?xt=urn:btih:UHWY2892JNEJ2GTEYOMDNU67E8ICGICYE92JDUGH&dn=baz&tr=udp://baz.com:1337",
				"magnet:?xt=urn:btih:HS263FG8U3GFIDHWD7829BYFCIXB78XIHG7CWCUG&dn=foz&tr=udp://foz.com:1337",
			},
			fail: vals{
				"",
				":?xt=urn:btih:06E2A9683BF4DA92C73A661AC56F0ECC9C63C5B4&dn=helloword2000&tr=udp://helloworld:1337/announce",
				"magnett:?xt=urn:btih:3E30322D5BFC7444B7B1D8DD42404B75D0531DFB&dn=world&tr=udp://world.com:1337",
				"xt=urn:btih:4ODKSDJBVMSDSNJVBCBFYFBKNRU875DW8D97DWC6&dn=helloworld&tr=udp://helloworld.com:1337",
				"magneta:?xt=urn:btih:1GSHJVBDVDVJFYEHKFHEFIO8573898434JBFEGHD&dn=foo&tr=udp://foo.com:1337",
				"magnet:?xt=uarn:btih:MCJDCYUFHEUD6E2752T7UJNEKHSUGEJFGTFHVBJS&dn=bar&tr=udp://bar.com:1337",
				"magnet:?xt=urn:btihz&dn=foobar&tr=udp://foobar.com:1337",
				"magnet:?xat=urn:btih:MKCJBHCBJDCU725TGEB3Y6RE8EJ2U267UNJFGUID&dn=test&tr=udp://test.com:1337",
				"magnet::?xt=urn:btih:UHWY2892JNEJ2GTEYOMDNU67E8ICGICYE92JDUGH&dn=baz&tr=udp://baz.com:1337",
				"magnet:?xt:btih:HS263FG8U3GFIDHWD7829BYFCIXB78XIHG7CWCUG&dn=foz&tr=udp://foz.com:1337",
			},
		}},
	}, {
		Name: "Match", Func: Match, Cases: Cases{{
			args: args{{`^(?i)testing$`}},
			pass: vals{
				"TESTING",
				"TeSTiNG",
				"testing",
			},
			fail: vals{
				"T3ST1NG",
				"bad testing",
				"foo bar",
				"",
			},
		}, {
			args: args{{`(?i)testing`}}, // unregistered
			fail: vals{
				"testing",
				"TESTING",
				"foo bar",
				"",
			},
		}},
	}, {
		Name: "MongoId", Func: MongoId, Cases: Cases{{
			pass: vals{
				"507f1f77bcf86cd799439011",
			},
			fail: vals{
				"507f1f77bcf86cd7994390",
				"507f1f77bcf86cd79943901z",
				"",
				"507f1f77bcf86cd799439011 ",
			},
		}},
	}, {
		Name: "Numeric", Func: Numeric, Cases: Cases{{
			pass: vals{
				"123",
				"00123",
				"-00123",
				"0",
				"-0",
				"+123",
				"123.123",
				"+000000",
			},
			fail: vals{
				" ",
				"",
				".",
			},
		}},
	}, {
		Name: "Octal", Func: Octal, Cases: Cases{{
			pass: vals{
				"076543210",
				"0o01234567",
			},
			fail: vals{
				"abcdefg",
				"012345678",
				"012345670c",
				"00c12345670c",
				"",
				"..",
			},
		}},
	}, {
		Name: "PAN", Func: PAN, Cases: Cases{{
			pass: vals{
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
			fail: vals{
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
	}, {
		Name: "PassportNumber", Func: todo_PassportNumber, Cases: Cases{{
			pass: vals{},
			fail: vals{},
		}},
	}, {
		Name: "Port", Func: Port, Cases: Cases{{
			pass: vals{
				"0",
				"22",
				"80",
				"443",
				"3000",
				"8080",
				"65535",
			},
			fail: vals{
				"",
				"-1",
				"65536",
			},
		}},
	}, {
		Name: "RGB", Func: RGB, Cases: Cases{{
			pass: vals{
				"rgb(0,0,0)",
				"rgb(255,255,255)",
				"rgba(0,0,0,0)",
				"rgba(255,255,255,1)",
				"rgba(255,255,255,.1)",
				"rgba(255,255,255,0.1)",
				"rgb(5%,5%,5%)",
				"rgba(5%,5%,5%,.3)",
			},
			fail: vals{
				"rgb(0,0,0,)",
				"rgb(0,0,)",
				"rgb(0,0,256)",
				"rgb()",
				"rgba(0,0,0)",
				"rgba(255,255,255,2)",
				"rgba(255,255,255,.12)",
				"rgba(255,255,256,0.1)",
				"rgb(4,4,5%)",
				"rgba(5%,5%,5%)",
				"rgba(3,3,3%,.3)",
				"rgb(101%,101%,101%)",
				"rgba(3%,3%,101%,0.3)",
			},
		}},
	}, {
		Name: "SSN", Func: SSN, Cases: Cases{{
			pass: vals{
				"001-01-0001",
				"012345678",
				"123-45-6789",
			},
			fail: vals{
				"000-01-0001",
				"001-00-0001",
				"001-01-0000",
				"0123456789",
				"12a-45-6789",
				"123-4b-6789",
				"123-45-678c",
				"123-45-67-89",
				"12345-6789",
			},
		}},
	}, {
		Name: "SemVer", Func: SemVer, Cases: Cases{{
			pass: vals{
				"0.0.4",
				"1.2.3",
				"10.20.30",
				"1.1.2-prerelease+meta",
				"1.1.2+meta",
				"1.1.2+meta-valid",
				"1.0.0-alpha",
				"1.0.0-beta",
				"1.0.0-alpha.beta",
				"1.0.0-alpha.beta.1",
				"1.0.0-alpha.1",
				"1.0.0-alpha0.valid",
				"1.0.0-alpha.0valid",
				"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
				"1.0.0-rc.1+build.1",
				"2.0.0-rc.1+build.123",
				"1.2.3-beta",
				"10.2.3-DEV-SNAPSHOT",
				"1.2.3-SNAPSHOT-123",
				"1.0.0",
				"2.0.0",
				"1.1.7",
				"2.0.0+build.1848",
				"2.0.1-alpha.1227",
				"1.0.0-alpha+beta",
				"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
				"1.2.3----R-S.12.9.1--.12+meta",
				"1.2.3----RC-SNAPSHOT.12.9.1--.12",
				"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
				"99999999999999999999999.999999999999999999.99999999999999999",
				"1.0.0-0A.is.legal",
			},
			fail: vals{
				"-invalid+invalid",
				"-invalid.01",
				"alpha",
				"alpha.beta",
				"alpha.beta.1",
				"alpha.1",
				"alpha+beta",
				"alpha_beta",
				"alpha.",
				"alpha..",
				"beta",
				"1.0.0-alpha_beta",
				"-alpha.",
				"1.0.0-alpha..",
				"1.0.0-alpha..1",
				"1.0.0-alpha...1",
				"1.0.0-alpha....1",
				"1.0.0-alpha.....1",
				"1.0.0-alpha......1",
				"1.0.0-alpha.......1",
				"01.1.1",
				"1.01.1",
				"1.1.01",
				"1.2",
				"1.2.3.DEV",
				"1.2-SNAPSHOT",
				"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
				"1.2-RC-SNAPSHOT",
				"-1.0.3-gamma+b7718",
				"+justmeta",
				"9.8.7+meta+meta",
				"9.8.7-whatever+meta+meta",
				"99999999999999999999999.999999999999999999.99999999999999999-",
				"---RC-SNAPSHOT.12.09.1--------------------------------..12",
			},
		}},
	}, {
		Name: "Slug", Func: Slug, Cases: Cases{{
			pass: vals{
				"cs-cz",
				"cscz",
			},
			fail: vals{
				"not-----------slug",
				"@#_$@",
				"-not-slug",
				"not-slug-",
				"_not-slug",
				"not-slug_",
				"not slug",
			},
		}},
	}, {
		Name: "StrongPassword", Func: StrongPassword, Cases: Cases{{
			args: args{{(*StrongPasswordOpts)(nil)}},
			pass: vals{
				`%2%k{7BsL"M%Kd6e`,
				`EXAMPLE of very long_password123!`,
				`mxH_+2vs&54_+H3P`,
				`+&DxJ=X7-4L8jRCD`,
				`etV*p%Nr6w&H%FeF`,
			},
			fail: vals{
				``,
				`password`,
				`hunter2`,
				`hello world`,
				`passw0rd`,
				`password!`,
				`PASSWORD!`,
			},
		}, {
			args: args{{&StrongPasswordOpts{
				MinLen:     12,
				MinLower:   3,
				MinUpper:   3,
				MinNumbers: 3,
				MinSymbols: 3,
			}}},
			pass: vals{
				`%2%k{7BsL"M%Kd6e`,
				`EXAMPLE of very long_password123!`,
				`mxH_+2vs&54_+H3P`,
			},
			fail: vals{
				``,
				`password`,
				`hunter2`,
				`hello world`,
				`passw0rd`,
				`password!`,
				`PASSWORD!`,
				`+&DxJ=X7-4L8jRCD`,
				`etV*p%Nr6w&H%FeF`,
			},
		}},
	}, {
		Name: "URL", Func: todo_URL, Cases: Cases{{
			pass: vals{},
			fail: vals{},
		}},
	}, {
		Name: "UUID", Func: UUID, Cases: Cases{{
			args: args{{3}},
			pass: vals{
				"A987FBC9-4BED-3078-CF07-9141BA07C9F3",
			},
			fail: vals{
				"",
				"xxxA987FBC9-4BED-3078-CF07-9141BA07C9F3",
				"934859",
				"AAAAAAAA-1111-1111-AAAG-111111111111",
				"A987FBC9-4BED-4078-8F07-9141BA07C9F3",
				"A987FBC9-4BED-5078-AF07-9141BA07C9F3",
			},
		}, {
			args: args{{4}},
			pass: vals{
				"713ae7e3-cb32-45f9-adcb-7c4fa86b90c1",
				"625e63f3-58f5-40b7-83a1-a72ad31acffb",
				"57b73598-8764-4ad0-a76a-679bb6640eb1",
				"9c858901-8a57-4791-81fe-4c455b099bc9",
			},
			fail: vals{
				"",
				"xxxA987FBC9-4BED-3078-CF07-9141BA07C9F3",
				"934859",
				"AAAAAAAA-1111-1111-AAAG-111111111111",
				"A987FBC9-4BED-5078-AF07-9141BA07C9F3",
				"A987FBC9-4BED-3078-CF07-9141BA07C9F3",
			},
		}, {
			args: args{{5}},
			pass: vals{
				"987FBC97-4BED-5078-AF07-9141BA07C9F3",
				"987FBC97-4BED-5078-BF07-9141BA07C9F3",
				"987FBC97-4BED-5078-8F07-9141BA07C9F3",
				"987FBC97-4BED-5078-9F07-9141BA07C9F3",
			},
			fail: vals{
				"",
				"xxxA987FBC9-4BED-3078-CF07-9141BA07C9F3",
				"934859",
				"AAAAAAAA-1111-1111-AAAG-111111111111",
				"9c858901-8a57-4791-81fe-4c455b099bc9",
				"A987FBC9-4BED-3078-CF07-9141BA07C9F3",
			},
		}},
	}, {
		Name: "Uint", Func: Uint, Cases: Cases{{
			pass: vals{
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
			fail: vals{
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
		Name: "UpperCase", Func: UpperCase, Cases: Cases{{
			pass: vals{},
			fail: vals{},
		}},
	}, {
		Name: "VAT", Func: VAT, Cases: Cases{{
			pass: vals{},
			fail: vals{},
		}},
	}}

	argstostr := func(args []reflect.Value) (str string) {
		for _, a := range args {
			str += fmt.Sprintf("%v, ", a.Interface())
		}
		if len(str) > 1 {
			str = str[:len(str)-2] // drop last ", "
		}
		return "[" + str + "]"
	}

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
						name := fmt.Sprintf("\"%v\"", val)
						t.Run(name, func(t *testing.T) {
							rv := reflect.ValueOf(val)
							vv := append([]reflect.Value{rv}, args...)

							got := fn.Call(vv)[0].Bool()
							if got != want {
								t.Errorf("got=%t; want=%t; args=%s", got, want, argstostr(args))
							}
						})
					}

					for _, val := range cases.fail {
						want := false
						name := fmt.Sprintf("\"%v\"", val)
						t.Run(name, func(t *testing.T) {
							rv := reflect.ValueOf(val)
							vv := append([]reflect.Value{rv}, args...)

							got := fn.Call(vv)[0].Bool()
							if got != want {
								t.Errorf("got=%t; want=%t; args=%s", got, want, argstostr(args))
							}
						})
					}
				}
			}
		})
	}
}
