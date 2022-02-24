package l10n_test

import (
	"testing"

	"github.com/frk/valid"
)

type List []struct {
	Name string
	Func func(v, cc string) bool
	Pass []string
	Fail []string
}

func Run(t *testing.T, ccs []string, list List) {
	for _, tt := range list {
		for _, cc := range ccs {
			name := cc + "/" + tt.Name
			for _, v := range tt.Pass {
				want := true
				t.Run(name, func(t *testing.T) {
					got := tt.Func(v, cc)
					if got != want {
						t.Errorf("got=%t; want=%t; value=%q", got, want, v)
					}
				})
			}
			for _, v := range tt.Fail {
				want := false
				t.Run(name, func(t *testing.T) {
					got := tt.Func(v, cc)
					if got != want {
						t.Errorf("got=%t; want=%t; value=%q", got, want, v)
					}
				})
			}
		}
	}
}

func Test(t *testing.T) {
	Run(t, []string{"AD", "AND"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+376312345",
			"312345",
		},
		Fail: []string{
			"31234",
			"31234567",
			"512345",
			"NotANumber",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"AD100",
			"AD200",
			"AD300",
			"AD400",
			"AD500",
			"AD600",
			"AD700",
		},
	}})

	Run(t, []string{"AE", "ARE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+971502674453",
			"+971521247658",
			"+971541255684",
			"+971555454458",
			"+971561498855",
			"+971585215778",
			"971585215778",
			"0585215778",
			"585215778",
		},
		Fail: []string{
			"12345",
			"+971511498855",
			"+9715614988556",
			"+9745614988556",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
		},
	}})

	Run(t, []string{"AF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})

	Run(t, []string{"AG", "ATG"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{},
		Fail: []string{},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{},
		Fail: []string{},
	}})

	Run(t, []string{"AI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})

	Run(t, []string{"AL", "ALB"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"067123456",
			"+35567123456",
		},
		Fail: []string{
			"67123456",
			"06712345",
			"06712345678",
			"065123456",
			"057123456",
			"NotANumber",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})

	Run(t, []string{"AM", "ARM"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+37410324123",
			"+37422298765",
			"+37431276521",
			"022698763",
			"37491987654",
			"+37494567890",
		},
		Fail: []string{
			"12345",
			"+37411498855",
			"+37411498123",
			"05614988556",
			"",
			"37456789000",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AO", "AGO"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+244911123432",
			"+244123091232",
		},
		Fail: []string{
			"+2449111234321",
			"31234",
			"31234567",
			"512345",
			"NotANumber",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AQ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AR", "ARG"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"5491143214321",
			"+5491143214321",
			"+5492414321432",
			"5498418432143",
		},
		Fail: []string{
			"1143214321",
			"91143214321",
			"+91143214321",
			"549841004321432",
			"549 11 43214321",
			"549111543214321",
			"5714003425432",
			"549114a214321",
			"54 9 11 4321-4321",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AS", "ASM"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AT", "AUT"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+436761234567",
			"06761234567",
			"00436123456789",
			"+436123456789",
			"01999",
			"+4372876",
			"06434908989562345",
		},
		Fail: []string{
			"167612345678",
			"1234",
			"064349089895623459",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AU", "AUS"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"61404111222",
			"+61411222333",
			"0417123456",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"04123456789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"4000",
			"2620",
			"3000",
			"2017",
			"0800",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"51824753556",
		},
		Fail: []string{
			"41824753556",
			"61824753556",
		},
	}})
	Run(t, []string{"AW", "ABW"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AX", "ALA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"AZ", "AZE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+994707007070",
			"0707007070",
			"+994502111111",
			"0505436743",
			"0554328772",
			"0993301022",
			"+994776007139",
		},
		Fail: []string{
			"wronumber",
			"",
			"994707007070",
			"++9945005050",
			"556007070",
			"1234566",
			"+994778008080a",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"AZ0100",
			"AZ0121",
			"AZ3500",
		},
		Fail: []string{
			"",
			" AZ0100",
			"AZ100",
			"AZ34340",
			"EN2020",
			"AY3030",
		},
	}})
	Run(t, []string{"BA", "BIH"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"060123456",
			"061123456",
			"062123456",
			"063123456",
			"0641234567",
			"065123456",
			"066123456",
			"+38760123456",
			"+38761123456",
			"+38762123456",
			"+38763123456",
			"+387641234567",
			"+38765123456",
			"+38766123456",
			"0038760123456",
			"0038761123456",
			"0038762123456",
			"0038763123456",
			"00387641234567",
			"0038765123456",
			"0038766123456",
		},
		Fail: []string{
			"0601234567",
			"0611234567",
			"06212345",
			"06312345",
			"064123456",
			"0651234567",
			"06612345",
			"+3866123456",
			"+3856123456",
			"00038760123456",
			"038761123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BB", "BRB"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BD", "BGD"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+8801794626846",
			"01399098893",
			"8801671163269",
			"01717112029",
			"8801898765432",
			"+8801312345678",
			"01494676946",
		},
		Fail: []string{
			"",
			"0174626346",
			"017943563469",
			"18001234567",
			"0131234567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BE", "BEL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0470123456",
			"+32470123456",
			"32470123456",
			"021234567",
			"+3221234567",
			"3221234567",
		},
		Fail: []string{
			"12345",
			"+3212345",
			"3212345",
			"04701234567",
			"+3204701234567",
			"3204701234567",
			"0212345678",
			"+320212345678",
			"320212345678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"BE0202239951",
			"BE0428759497",
		},
		Fail: []string{
			"BE0102239951",
			"BE0431150351",
		},
	}})
	Run(t, []string{"BF", "BFA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BG", "BGR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+359897123456",
			"+359898888888",
			"0897123123",
		},
		Fail: []string{
			"",
			"0898123",
			"+359212555666",
			"18001234567",
			"12125559999",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"1000",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"BG123456789",
			"BG1234567890",
		},
		Fail: []string{
			"BG12345678",
			"BG12345678901",
		},
	}})
	Run(t, []string{"BH", "BHR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+97335078110",
			"+97339534385",
			"+97366331055",
			"+97333146000",
			"97335078110",
			"35078110",
			"66331055",
		},
		Fail: []string{
			"12345",
			"+973350781101",
			"+97379534385",
			"+973035078110",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
			"035078110",
			"16331055",
			"hello",
			"+9733507811a",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BJ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BL", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BO", "BOL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+59175553635",
			"+59162223685",
			"+59179783890",
			"+59160081890",
			"79783890",
			"60081890",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+59199783890",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BQ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BR", "BRA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+55 12 996551215",
			"+55 15 97661234",
			"+55 (12) 996551215",
			"+55 (15) 97661234",
			"55 (17) 96332-2155",
			"55 (17) 6332-2155",
			"55 15 976612345",
			"55 15 75661234",
			"+5512984567890",
			"+551283456789",
			"5512984567890",
			"551283456789",
			"015994569878",
			"01593456987",
			"022995678947",
			"02299567894",
			"(22)99567894",
			"(22)9956-7894",
			"(22) 99567894",
			"(22) 9956-7894",
			"(22)999567894",
			"(22)99956-7894",
			"(22) 999567894",
			"(22) 99956-7894",
			"(11) 94123-4567",
		},
		Fail: []string{
			"0819876543",
			"+55 15 7566123",
			"+017 123456789",
			"5501599623874",
			"+55012962308",
			"+55 015 1234-3214",
			"+55 11 91431-4567",
			"+55 (11) 91431-4567",
			"+551191431-4567",
			"5511914314567",
			"5511912345678",
			"(11) 91431-4567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"39100-000",
			"22040-020",
			"39400-152",
		},
		Fail: []string{
			"79800A12",
			"13165-00",
			"38175-abc",
			"81470-2763",
			"78908",
			"13010|111",
		},
	}})
	Run(t, []string{"BS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BT", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BW", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BY", "BLR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+375241234567",
			"+375251234567",
			"+375291234567",
			"+375331234567",
			"+375441234567",
			"375331234567",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"12345",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"225320",
			"211120",
			"247710",
			"231960",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"BZ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CA", "CAN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"19876543210",
			"8005552222",
			"+15673628910",
		},
		Fail: []string{
			"564785",
			"0123456789",
			"1437439210",
			"+10345672645",
			"11435213543",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"L4T 0A5",
			"G1A-0A2",
			"A1A 1A1",
			"X0A-0H0",
			"V5K 0A1",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CD", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CH", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CL", "CHL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+56733875615",
			"56928590234",
			"0928590294",
			"0208590294",
		},
		Fail: []string{
			"1234",
			"+5633875615",
			"563875615",
			"56109834567",
			"56069834567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CN", "CHN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"15323456787",
			"13523333233",
			"13898728332",
			"+8613238234822",
			"+8613487234567",
			"+8617823492338",
			"+8617823492338",
			"16637108167",
			"+8616637108167",
			"+8616637108167",
			"+8616712341234",
			"008618812341234",
			"008618812341234",
			"+8619912341234",
			"+8619812341234",
			"+8619712341234",
			"+8619612341234",
			"+8619512341234",
			"+8619312341234",
			"+8619212341234",
			"+8619112341234",
			"17269427292",
			"16565600001",
			"+8617269427292",
			"008617269427292",
		},
		Fail: []string{
			"12345",
			"",
			"+08613811211114",
			"+008613811211114",
			"08613811211114",
			"0086-13811211114",
			"0086-138-1121-1114",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"150237",
			"100000",
		},
		Fail: []string{
			"141234",
			"386789",
			"ab1234",
		},
	}})
	Run(t, []string{"CO", "COL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+573003321235",
			"573003321235",
			"579871235",
			"3003321235",
			"3213321235",
			"3103321235",
			"3253321235",
			"3321235",
			"574321235",
			"5784321235",
			"5784321235",
			"9821235",
			"573011140876",
			"0698345",
		},
		Fail: []string{
			"1234",
			"+57443875615",
			"57309875615",
			"57109834567",
			"5792434567",
			"5702345689",
			"5714003425432",
			"5703013347567",
			"069834567",
			"969834567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CR", "CRI"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+50688888888",
			"+50665408090",
			"+50640895069",
			"25789563",
			"85789563",
		},
		Fail: []string{
			"+5081",
			"+5067777777",
			"+50188888888",
			"+50e987643254",
			"+506e4t4",
			"-50688888888",
			"50688888888",
			"12345678",
			"98765432",
			"01234567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CU", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CW", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CX", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CY", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"CZ", "CZE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+420 123 456 789",
			"+420 123456789",
			"+420123456789",
			"123 456 789",
			"123456789",
		},
		Fail: []string{
			"",
			"+42012345678",
			"+421 123 456 789",
			"+420 023456789",
			"+4201234567892",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"20134",
			"392 90",
			"39919",
			"938 29",
			"39949",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"DE", "DEU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+49015123456789",
			"+4915123456789",
			"+4930405044550",
			"015123456789",
			"15123456789",
			"15623456789",
			"15623456789",
			"1601234567",
			"16012345678",
			"1621234567",
			"1631234567",
			"1701234567",
			"17612345678",
			"15345678910",
			"15412345678",
		},
		Fail: []string{
			"34412345678",
			"14412345678",
			"16212345678",
			"1761234567",
			"16412345678",
			"17012345678",
			"+4912345678910",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"DJ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"DK", "DNK"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"12345678",
			"12 34 56 78",
			"45 12345678",
			"4512345678",
			"45 12 34 56 78",
			"+45 12 34 56 78",
		},
		Fail: []string{
			"",
			"+45010203",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"12 34 56",
			"123 123 12",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"DK31569931",
			"DK31948959",
			"DK38484621",
			"DK31329566",
		},
		Fail: []string{
			"DK31569941",
			"DK31948859",
			"DK38484641",
			"DK31329567",
		},
	}})
	Run(t, []string{"DM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"DO", "DOM"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+18096622563",
			"+18295614488",
			"+18495259567",
			"8492283478",
			"8092324576",
			"8292387713",
		},
		Fail: []string{
			"+18091",
			"+1849777777",
			"-18296643245",
			"+18086643245",
			"+18396643245",
			"8196643245",
			"+38492283478",
			"6492283478",
			"8192283478",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"12345",
		},
		Fail: []string{
			"A1234",
			"123",
			"123456",
		},
	}})
	Run(t, []string{"DZ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"EC", "ECU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+593987654321",
			"593987654321",
			"0987654321",
			"027332615",
			"+59323456789",
		},
		Fail: []string{
			"03321321",
			"+593387561",
			"59312345677",
			"02344635",
			"593123456789",
			"081234567",
			"+593912345678",
			"+593902345678",
			"+593287654321",
			"593287654321",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"EE", "EST"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+372 512 34 567",
			"372 512 34 567",
			"+37251234567",
			"51234567",
			"81234567",
			"+372842345678",
		},
		Fail: []string{
			"12345",
			"",
			"NotANumber",
			"+333 51234567",
			"61234567",
			"+51234567",
			"+372 539 57 4",
			"+372 900 1234",
			"12345678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"EG", "EGY"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+201004513789",
			"+201111453489",
			"+201221204610",
			"+201144621154",
			"+201200124304",
			"+201011201564",
			"+201124679001",
			"+201064790156",
			"+201274652177",
			"+201280134679",
			"+201090124576",
			"+201583728900",
			"201599495596",
			"201090124576",
			"01090124576",
			"01538920744",
			"1593075993",
			"1090124576",
		},
		Fail: []string{
			"+221004513789",
			"+201404513789",
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"EH", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ER", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ES", "ESP"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+34654789321",
			"654789321",
			"+34714789321",
			"714789321",
			"+34744789321",
			"744789321",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+3465478932",
			"65478932",
			"+346547893210",
			"6547893210",
			"+3470478932",
			"7047893210",
			"+34854789321",
			"7547893219",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"01001",
			"52999",
			"27880",
		},
		Fail: []string{
			"123",
			"1234",
			"53000",
			"052999",
			"0123",
			"abcde",
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"ESX9999999R",
			"ESX99999999",
			"ES99999999R",
		},
		Fail: []string{
			"ESX9999999",
			"ES9999999R",
			"ESXR9999999",
			"ES9999999XR",
		},
	}})
	Run(t, []string{"ET", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FI", "FIN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+358505557171",
			"0455571",
			"0505557171",
			"358505557171",
			"04412345",
			"0457 123 45 67",
			"+358457 123 45 67",
			"+358 50 555 7171",
		},
		Fail: []string{
			"12345",
			"",
			"045557",
			"045555717112312332423423421",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"+3-585-0555-7171",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
			"044123",
			"019123456789012345678901",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FJ", "FJI"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+6799898679",
			"6793788679",
			"+679 989 8679",
			"679 989 8679",
			"679 3456799",
			"679908 8909",
		},
		Fail: []string{
			"12345",
			"",
			"04555792",
			"902w99900030900000000099",
			"8uiuiuhhyy&GUU88d",
			"010-38238383",
			"19676338855",
			"679 9 89 8679",
			"6793 45679",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FO", "FRO"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"123456",
			"12 34 56",
			"298 123456",
			"298123456",
			"298 12 34 56",
			"+298 12 34 56",
		},
		Fail: []string{
			"",
			"+4501020304",
			"ASDFGJKLmZXJtZtesting123",
			"12345678",
			"12 34 56 78",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"FR", "FRA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0612457898",
			"+33612457898",
			"33612457898",
			"0712457898",
			"+33712457898",
			"33712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+34612457898",
			"+336124578980",
			"+3361245789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"75008",
			"44 522",
			"98025",
			"38 499",
			"39940",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"FR42813454717",
			"FR30803417153",
			"FR83404833048",
			"FR40303265045",
			"FR23334175221",
		},
		Fail: []string{
			"813454717",
			"FR4281345471",
			"FR428134547171",
			"FR84323140391",
		},
	}})
	Run(t, []string{"GA", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GB", "GBR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"447789345856",
			"+447861235675",
			"07888814488",
		},
		Fail: []string{
			"67699567",
			"0773894868",
			"077389f8688",
			"+07888814488",
			"0152456999",
			"442073456754",
			"+443003434751",
			"05073456754",
			"08001123123",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"TW8 9GS",
			"BS98 1TL",
			"DE99 3GG",
			"DE55 4SW",
			"DH98 1BT",
			"DH99 1NS",
			"GIR0aa",
			"SA99",
			"W1N 4DJ",
			"AA9A 9AA",
			"AA99 9AA",
			"BS98 1TL",
			"DE993GG",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GD", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GE", "GEO"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+99550001111",
			"+99551535213",
			"+995798526662",
			"798526662",
			"50001111",
			"798526662",
			"+995799766525",
		},
		Fail: []string{
			"+995500011118",
			"+9957997665250",
			"+995999766525",
			"20000000000",
			"68129485729",
			"6589394827",
			"298RI89572",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GF", "GUF"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0612457898",
			"+594612457898",
			"594612457898",
			"0712457898",
			"+594712457898",
			"594712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+54612457898",
			"+5946124578980",
			"+59461245789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GG", "GGY"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+441481123456",
			"+441481789123",
			"441481123456",
			"441481789123",
		},
		Fail: []string{
			"999",
			"+441481123456789",
			"+447123456789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GH", "GHA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0202345671",
			"0502345671",
			"0242345671",
			"0542345671",
			"0272345671",
			"0572345671",
			"0262345671",
			"0562345671",
			"0232345671",
			"0282345671",
			"+233202345671",
			"+233502345671",
			"+233242345671",
			"+233542345671",
			"+233272345671",
			"+233572345671",
			"+233262345671",
			"+233562345671",
			"+233232345671",
			"+233282345671",
		},
		Fail: []string{
			"082123",
			"232345671",
			"0292345671",
			"+233292345671",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GL", "GRL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"123456",
			"12 34 56",
			"299 123456",
			"299123456",
			"299 12 34 56",
			"+299 12 34 56",
		},
		Fail: []string{
			"",
			"+4501020304",
			"ASDFGJKLmZXJtZtesting123",
			"12345678",
			"12 34 56 78",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GP", "GLP"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0612457898",
			"+590612457898",
			"590612457898",
			"0712457898",
			"+590712457898",
			"590712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+594612457898",
			"+5906124578980",
			"+59061245789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GQ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GR", "GRC"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+306944848966",
			"6944848966",
			"306944848966",
		},
		Fail: []string{
			"2102323234",
			"+302646041461",
			"120000000",
			"20000000000",
			"68129485729",
			"6589394827",
			"298RI89572",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"022 93",
			"29934",
			"90293",
			"299 42",
			"94944",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{},
		Fail: []string{},
	}})
	Run(t, []string{"GS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GT", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GU", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GW", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"GY", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"HK", "HKG"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"91234567",
			"9123-4567",
			"61234567",
			"51234567",
			"+85291234567",
			"+852-91234567",
			"+852-9123-4567",
			"+852 9123 4567",
			"9123 4567",
			"852-91234567",
		},
		Fail: []string{
			"999",
			"+852-912345678",
			"123456789",
			"+852-1234-56789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"HM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"HN", "HND"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+50495551876",
			"+50488908787",
			"+50493456789",
			"+50489234567",
			"+50488987896",
			"+50497567389",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+34683456543",
			"65478932",
			"+50298787654",
			"+504989874",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"HR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"HT", "HTI"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"HT1234",
		},
		Fail: []string{
			"HT123",
			"HT12345",
			"AA1234",
		},
	}})
	Run(t, []string{"HU", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ID", "IDN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0811 778 998",
			"0811 7785 9983",
			"0812 7784 9984",
			"0813 7782 9982",
			"0821 1234 1234",
			"0822 1234 1234",
			"0823 1234 1234",
			"0852 1234 6764",
			"0853 1234 6764",
			"0851 1234 6764",
			"0814 7782 9982",
			"0815 7782 9982",
			"0816 7782 9982",
			"0855 7782 9982",
			"0856 7782 9982",
			"0857 7782 9982",
			"0858 7782 9982",
			"0817 7785 9983",
			"0818 7784 9984",
			"0819 7782 9982",
			"0859 1234 1234",
			"0877 1234 1234",
			"0878 1234 1234",
			"0895 7785 9983",
			"0896 7784 9984",
			"0897 7782 9982",
			"0898 1234 1234",
			"0899 1234 1234",
			"0881 7785 9983",
			"0882 7784 9984",
			"0883 7782 9982",
			"0884 1234 1234",
			"0886 1234 1234",
			"0887 1234 1234",
			"0888 7785 9983",
			"0889 7784 9984",
			"0828 7784 9984",
			"0838 7784 9984",
			"0831 7784 9984",
			"0832 7784 9984",
			"0833 7784 9984",
			"089931236181900",
			"62811 778 998",
			"62811778998",
			"628993123618190",
			"62898 740123456",
			"62899 7401 2346",
			"+62811 778 998",
			"+62811778998",
			"+62812 9650 3508",
			"08197231819",
			"085361008008",
			"+62811787391",
		},
		Fail: []string{
			"0899312361819001",
			"0217123456",
			"622178878890",
			"6221 740123456",
			"0341 8123456",
			"0778 89800910",
			"0741 123456",
			"+6221740123456",
			"+65740 123 456",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"740123456",
			"+65640123456",
			"+64210123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"10210",
			"40181",
			"55161",
			"60233",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"02.271.824.1-413.000",
			"022718241413000",
			"09.254.294.3-407.000",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"IE", "IRL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+353871234567",
			"353831234567",
			"353851234567",
			"353861234567",
			"353871234567",
			"353881234567",
			"353891234567",
			"0871234567",
			"0851234567",
		},
		Fail: []string{
			"999",
			"+353341234567",
			"+33589484858",
			"353841234567",
			"353811234567",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"A65 TF12",
			"D02 AF30",
		},
		Fail: []string{
			"123",
			"A6W U9U9",
			"75690HG",
			"AW5  TF12",
			"AW5 TF12",
			"756  90HG",
			"A65T F12",
			"O62 O1O2",
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"IE1234567T",
			"IE1234567TW",
			"IE1234567FA",
			"IE1A23456B",
			"IE1+23456B",
			"IE1*23456B",
		},
		Fail: []string{
			"IE123456T",
			"IE12345678",
			"IE12345678W",
			"IE1234567FAZ",
			"IE12345678AZ",
			"IE1A234567",
			"IE1-23456B",
		},
	}})
	Run(t, []string{"IL", "ISR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"10200",
			"10292",
			"10300",
			"10329",
			"3885500",
			"4290500",
			"4286000",
			"7080000",
		},
		Fail: []string{
			"123",
			"012345",
			"011111",
			"101123",
			"291123",
			"351123",
			"541123",
			"551123",
			"651123",
			"661123",
			"861123",
			"871123",
			"881123",
			"891123",
		},
	}})
	Run(t, []string{"IM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"IN", "IND"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"364240",
			"360005",
		},
		Fail: []string{
			"123",
			"012345",
			"011111",
			"101123",
			"291123",
			"351123",
			"541123",
			"551123",
			"651123",
			"661123",
			"861123",
			"871123",
			"881123",
			"891123",
		},
	}})
	Run(t, []string{"IO", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"IQ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"IR", "IRN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+989123456789",
			"989223456789",
			"09323456789",
			"09021456789",
			"+98-990-345-6789",
			"+98 938 345 6789",
			"0938 345 6789",
		},
		Fail: []string{
			"",
			"+989623456789",
			"+981123456789",
			"01234567890",
			"09423456789",
			"09823456789",
			"9123456789",
			"091234567890",
			"0912345678",
			"+98 912 3456 6789",
			"0912 345 678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"4351666456",
			"5614736867",
		},
		Fail: []string{
			"43516 6456",
			"123443516 6456",
			"891123",
		},
	}})
	Run(t, []string{"IS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"IT", "ITA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"370 3175423",
			"333202925",
			"+39 310 7688449",
			"+39 3339847632",
		},
		Fail: []string{
			"011 7387545",
			"12345",
			"+45 345 6782395",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"IT07643520567",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"JE", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"JM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"JO", "JOR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0796477263",
			"0777866254",
			"0786725261",
			"+962796477263",
			"+962777866254",
			"+962786725261",
			"962796477263",
			"962777866254",
			"962786725261",
		},
		Fail: []string{
			"00962786725261",
			"00962796477263",
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"JP", "JPN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"09012345678",
			"08012345678",
			"07012345678",
			"06012345678",
			"090 1234 5678",
			"+8190-1234-5678",
			"+81 (0)90-1234-5678",
			"+819012345678",
			"+81-(0)90-1234-5678",
			"+81 90 1234 5678",
		},
		Fail: []string{
			"12345",
			"",
			"045555717112312332423423421",
			"Vml2YW11cyBmZXJtZtesting123",
			"+3-585-0555-7171",
			"0 1234 5689",
			"16 1234 5689",
			"03_1234_5689",
			"0312345678",
			"0721234567",
			"06 1234 5678",
			"072 123 4567",
			"0729 12 3456",
			"07296 1 2345",
			"072961 2345",
			"03-1234-5678",
			"+81312345678",
			"+816-1234-5678",
			"+81 090 1234 5678",
			"+8109012345678",
			"+81-090-1234-5678",
			"90 1234 5678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"135-0000",
			"874-8577",
			"669-1161",
			"470-0156",
			"672-8031",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KE", "KEN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+254728590432",
			"+254733875610",
			"254728590234",
			"0733346543",
			"0700459022",
			"0110934567",
			"+254110456794",
			"254198452389",
		},
		Fail: []string{
			"999",
			"+25489032",
			"123456789",
			"+254800723845",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KH", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KP", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KR", "KOR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+82-010-1234-5678",
			"+82-10-1234-5678",
			"82-010-1234-5678",
			"82-10-1234-5678",
			"+82 10 1234 5678",
			"010-123-5678",
			"10-1234-5678",
			"+82 10 1234 5678",
			"011 1234 5678",
			"+820112345678",
			"01012345678",
			"+82 016 1234 5678",
			"82 19 1234 5678",
			"+82 010 12345678",
		},
		Fail: []string{
			"abcdefghi",
			"+82 10 1234 567",
			"+82 10o 1234 1234",
			"+82 101 1234 5678",
			"+82 10 12 5678",
			"+011 7766 1234",
			"011_7766_1234",
			"+820 11 7766 1234",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KW", "KWT"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"96550000000",
			"96560000000",
			"96590000000",
			"+96550000000",
			"+96550000220",
			"+96551111220",
		},
		Fail: []string{
			"+96570000220",
			"00962786725261",
			"00962796477263",
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KY", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"KZ", "KAZ"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+77254716212",
			"77254716212",
			"87254716212",
			"7254716212",
		},
		Fail: []string{
			"12345",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LA", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LB", "LBN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+96171234568",
			"+9613123456",
			"3456123",
			"3123456",
			"81978468",
			"77675798",
		},
		Fail: []string{
			"+961712345688888",
			"00912220000",
			"7767579888",
			"+0921110000",
			"+3123456888",
			"021222200000",
			"213333444444",
			"",
			"+212234",
			"+21",
			"02122333",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LI", "LIE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"9485",
			"9497",
			"9491",
			"9489",
			"9496",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LT", "LTU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+37051234567",
			"851234567",
		},
		Fail: []string{
			"+65740 123 456",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"740123456",
			"+65640123456",
			"+64210123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LU", "LUX"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"601123456",
			"+352601123456",
		},
		Fail: []string{
			"NaN",
			"791234",
			"+352791234",
			"26791234",
			"+35226791234",
			"+112039812",
			"+352703123456",
			"1234",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"LY", "LBY"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"912220000",
			"0923330000",
			"218945550000",
			"+218958880000",
			"212220000",
			"0212220000",
			"+218212220000",
		},
		Fail: []string{
			"9122220000",
			"00912220000",
			"09211110000",
			"+0921110000",
			"+2180921110000",
			"021222200000",
			"213333444444",
			"",
			"+212234",
			"+21",
			"02122333",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MA", "MAR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0522714782",
			"0690851123",
			"0708186135",
			"+212522714782",
			"+212690851123",
			"+212708186135",
			"00212522714782",
			"00212690851123",
			"00212708186135",
		},
		Fail: []string{
			"522714782",
			"690851123",
			"708186135",
			"212522714782",
			"212690851123",
			"212708186135",
			"0212522714782",
			"0212690851123",
			"0212708186135",
			"",
			"12345",
			"0922714782",
			"+212190851123",
			"00212408186135",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MD", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ME", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MH", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ML", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MO", "MAC"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"61234567",
			"+85361234567",
			"+853-61234567",
			"+853-6123-4567",
			"+853 6123 4567",
			"6123 4567",
			"853-61234567",
		},
		Fail: []string{
			"999",
			"12345678",
			"612345678",
			"+853-12345678",
			"+853-22345678",
			"+853-82345678",
			"+853-612345678",
			"+853-1234-5678",
			"+853 1234 5678",
			"+853-6123-45678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MP", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MQ", "MTQ"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0612457898",
			"+596612457898",
			"596612457898",
			"0712457898",
			"+596712457898",
			"596712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+594612457898",
			"+5966124578980",
			"+59661245789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MT", "MLT"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+35699000000",
			"+35679000000",
			"99000000",
		},
		Fail: []string{
			"356",
			"+35699000",
			"+35610000000",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"VLT2345",
			"VLT 2345",
			"ATD1234",
			"MSK8723",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MU", "MUS"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+23012341234",
			"12341234",
			"012341234",
		},
		Fail: []string{
			"41234",
			"",
			"+230",
			"+2301",
			"+23012",
			"+230123",
			"+2301234",
			"+23012341",
			"+230123412",
			"+2301234123",
			"+230123412341",
			"+2301234123412",
			"+23012341234123",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MW", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MX", "MEX"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+52019654789321",
			"+52199654789321",
			"+5201965478932",
			"+5219654789321",
			"52019654789321",
			"52199654789321",
			"5201965478932",
			"5219654789321",
			"87654789321",
			"8654789321",
			"0187654789321",
			"18654789321",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+3465478932",
			"65478932",
			"+346547893210",
			"+34704789321",
			"704789321",
			"+34754789321",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MY", "MYS"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+60128228789",
			"+60195830837",
			"+6019-5830837",
			"+6019-5830837",
			"+6010-4357675",
			"+60172012370",
			"0128737867",
			"0172012370",
			"01468987837",
			"01112347345",
			"016-2838768",
			"016 2838768",
		},
		Fail: []string{
			"12345",
			"601238788657",
			"088387675",
			"16-2838768",
			"032551433",
			"6088-387888",
			"088-261987",
			"1800-88-8687",
			"088-320000",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"56000",
			"12000",
			"79502",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"MZ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NA", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NE", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NL", "NLD"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0670123456",
			"0612345678",
			"31612345678",
			"31670123456",
			"+31612345678",
			"+31670123456",
			"+31(0)612345678",
			"0031612345678",
			"0031(0)612345678",
		},
		Fail: []string{
			"12345",
			"+3112345",
			"3112345",
			"06701234567",
			"012345678",
			"+3104701234567",
			"3104701234567",
			"0212345678",
			"021234567",
			"+3121234567",
			"3121234567",
			"+310212345678",
			"310212345678",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"1012 SZ",
			"3432FE",
			"1118 BH",
			"3950IO",
			"3997 GH",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"NL999999999B01",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NO", "NOR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+4796338855",
			"+4746338855",
			"4796338855",
			"4746338855",
			"46338855",
			"96338855",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+4676338855",
			"19676338855",
			"+4726338855",
			"4736338855",
			"66338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NP", "NPL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+9779817385479",
			"+9779717385478",
			"+9779862002615",
			"+9779853660020",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+97796123456789",
			"+9771234567",
			"+977981234",
			"4736338855",
			"66338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"10811",
			"32600",
			"56806",
			"977",
		},
		Fail: []string{
			"11977",
			"asds",
			"13 32",
			"-977",
			"97765",
		},
	}})
	Run(t, []string{"NR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NU", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"NZ", "NZL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+6427987035",
			"642240512347",
			"0293981646",
			"029968425",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+642956696123566",
			"+02119620856",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"7843",
			"3581",
			"0449",
			"0984",
			"4144",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"OM", "OMN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+96891212121",
			"0096899999999",
			"93112211",
			"99099009",
		},
		Fail: []string{
			"+96890212121",
			"0096890999999",
			"0090999999",
			"+9689021212",
			"",
			"+212234",
			"+21",
			"02122333",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PA", "PAN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+5076784565",
			"+5074321557",
			"5073331112",
			"+50723431212",
		},
		Fail: []string{
			"+50755555",
			"+207123456",
			"2001236542",
			"+507987643254",
			"+507jjjghtf",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PE", "PER"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+51912232764",
			"+51923464567",
			"+51968267382",
			"+51908792973",
			"974980472",
			"908792973",
			"+51974980472",
		},
		Fail: []string{
			"999",
			"+51812232764",
			"+5181223276499",
			"+25589032",
			"123456789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PH", "PHL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+639275149120",
			"+639275142327",
			"+639003002023",
			"09275149116",
			"09194877624",
		},
		Fail: []string{
			"12112-13-345",
			"12345678901",
			"sx23YW11cyBmZxxXJt123123",
			"010-38238383",
			"966684123123-2590",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PL", "POL"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+48512689767",
			"+48 56 376 87 47",
			"56 566 78 46",
			"657562855",
			"+48657562855",
			"+48 887472765",
			"+48 56 6572724",
			"+48 67 621 5461",
			"48 67 621 5461",
		},
		Fail: []string{
			"+48  67 621 5461",
			"+55657562855",
			"3454535",
			"teststring",
			"",
			"1800-88-8687",
			"+6019-5830837",
			"357562855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"47-260",
			"12-930",
			"78-399",
			"39-490",
			"38-483",
		},
		Fail: []string{
			"360",
			"90312",
			"399",
			"935",
			"38842",
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"PL123-456-32-18",
			"PL000-000-00-00",
		},
		Fail: []string{
			"PL123-456-78-90",
		},
	}})
	Run(t, []string{"PM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PR", "PRI"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"00979",
			"00631",
			"00786",
			"00987",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PT", "PRT"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"4829-489",
			"0294-348",
			"8156-392",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"PT999999990",
			"PT501442600",
		},
		Fail: []string{
			"PT999999999",
		},
	}})
	Run(t, []string{"PW", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"PY", "PRY"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+595991372649",
			"+595992847352",
			"+595993847593",
			"+595994857473",
			"+595995348532",
			"+595996435231",
			"+595981847362",
			"+595982435452",
			"+595983948502",
			"+595984342351",
			"+595985403481",
			"+595986384012",
			"+595971435231",
			"+595972103924",
			"+595973438542",
			"+595974425864",
			"+595975425843",
			"+595976342546",
			"+595961435234",
			"+595963425043",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"65478932",
			"+59599384712",
			"+5959938471234",
			"+595547893218",
			"+591993546843",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"QA", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"RE", "REU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0612457898",
			"+262612457898",
			"262612457898",
			"0712457898",
			"+262712457898",
			"262712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+264612457898",
			"+2626124578980",
			"+26261245789",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"RO", "ROU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+40740123456",
			"+40 740123456",
			"+40740 123 456",
			"+40740.123.456",
			"+40740-123-456",
			"40740123456",
			"40 740123456",
			"40740 123 456",
			"40740.123.456",
			"40740-123-456",
			"0740123456",
			"0740/123456",
			"0740 123 456",
			"0740.123.456",
			"0740-123-456",
		},
		Fail: []string{
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"123456",
			"740123456",
			"+40640123456",
			"+40210123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"RS", "SRB"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0640133338",
			"063333133",
			"0668888878",
			"+381645678912",
			"+381611314000",
			"0655885010",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"RU", "RUS"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+79676338855",
			"79676338855",
			"89676338855",
			"9676338855",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"RW", "RWA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+250728590432",
			"+250733875610",
			"250738590234",
			"0753346543",
			"0780459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25089032",
			"123456789",
			"+250800723845",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SA", "SAU"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0556578654",
			"+966556578654",
			"966556578654",
			"596578654",
			"572655597",
		},
		Fail: []string{
			"12345",
			"",
			"+9665626626262",
			"+96633221097",
			"0114152198",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SB", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SD", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SE", "SWE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+46701234567",
			"46701234567",
			"0721234567",
			"073-1234567",
			"0761-234567",
			"079-123 45 67",
		},
		Fail: []string{
			"12345",
			"+4670123456",
			"+46301234567",
			"+0731234567",
			"0731234 56",
			"+7312345678",
			"",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"12994",
			"284 39",
			"39556",
			"489 39",
			"499 49",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"SE556293998201",
		},
		Fail: []string{
			"SE556293998301",
			"SE556293998101",
		},
	}})
	Run(t, []string{"SG", "SGP"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"87654321",
			"98765432",
			"+6587654321",
			"+6598765432",
			"+6565241234",
		},
		Fail: []string{
			"987654321",
			"876543219",
			"8765432",
			"9876543",
			"12345678",
			"+98765432",
			"+9876543212",
			"+15673628910",
			"19876543210",
			"8005552222",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"308215",
			"546080",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SH", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SI", "SVN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: valid.VAT,
		Pass: []string{
			"SI99662981",
			"SI19136234",
		},
		Fail: []string{
			"SI99662982",
			"SI19136235",
		},
	}})
	Run(t, []string{"SJ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SK", "SVK"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+421 123 456 789",
			"+421 123456789",
			"+421123456789",
			"123 456 789",
			"123456789",
		},
		Fail: []string{
			"",
			"+42112345678",
			"+422 123 456 789",
			"+421 023456789",
			"+4211234567892",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SL", "SLE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+94766661206",
			"94713114340",
			"0786642116",
			"078 7642116",
			"078-7642116",
		},
		Fail: []string{
			"9912349956789",
			"12345",
			"1678123456",
			"0731234567",
			"0749994567",
			"0797878674",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SM", "SMR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"612345",
			"05496123456",
			"+37861234567",
			"+390549612345678",
			"+37805496123456789",
		},
		Fail: []string{
			"61234567890",
			"6123",
			"1234567",
			"+49123456",
			"NotANumber",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SN", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SO", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ST", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SX", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SY", "SYR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0944549710",
			"+963944549710",
			"956654379",
			"0944549710",
			"0962655597",
		},
		Fail: []string{
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"SZ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TD", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TH", "THA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0912345678",
			"+66912345678",
			"66912345678",
		},
		Fail: []string{
			"99123456789",
			"12345",
			"67812345623",
			"081234567891",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"10250",
			"72170",
			"12140",
		},
		Fail: []string{
			"T1025",
			"T72170",
			"12140TH",
		},
	}})
	Run(t, []string{"TJ", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TK", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TL", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TN", "TUN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"23456789",
			"+21623456789",
			"21623456789",
		},
		Fail: []string{
			"12345",
			"75200123",
			"+216512345678",
			"13520459",
			"85479520",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TO", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TR", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TT", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TV", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TW", "TWN"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0987123456",
			"+886987123456",
			"886987123456",
			"+886-987123456",
			"886-987123456",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"0-987123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"TZ", "TZA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+255728590432",
			"+255733875610",
			"255628590234",
			"0673346543",
			"0600459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25589032",
			"123456789",
			"+255800723845",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"UA", "UKR"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+380982345679",
			"380982345679",
			"80982345679",
			"0982345679",
		},
		Fail: []string{
			"+30982345679",
			"982345679",
			"+380 98 234 5679",
			"+380-98-234-5679",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"740123456",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			"65000",
			"65080",
			"01000",
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"UG", "UGA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+256728590432",
			"+256733875610",
			"256728590234",
			"0773346543",
			"0700459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25489032",
			"123456789",
			"+254800723845",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"UM", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"US", "USA"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"19876543210",
			"8005552222",
			"+15673628910",
			"+1(567)3628910",
			"+1(567)362-8910",
			"+1(567) 362-8910",
			"1(567)362-8910",
			"1(567)362 8910",
			"223-456-7890",
		},
		Fail: []string{
			"564785",
			"0123456789",
			"1437439210",
			"+10345672645",
			"11435213543",
			"1(067)362-8910",
			"1(167)362-8910",
			"+2(267)362-8910",
			"+3365520145",
		},
	}})
	Run(t, []string{"UY", "URY"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+59899123456",
			"099123456",
			"+59894654321",
			"091111111",
		},
		Fail: []string{
			"54321",
			"montevideo",
			"",
			"+598099123456",
			"090883338",
			"099 999 999",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"UZ", "UZB"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+998664835244",
			"998664835244",
			"664835244",
			"+998957124555",
			"998957124555",
			"957124555",
		},
		Fail: []string{
			"+998644835244",
			"998644835244",
			"644835244",
			"+99664835244",
			"ASDFGJKLmZXJtZtesting123",
			"123456789",
			"870123456",
			"",
			"+998",
			"998",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VA", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VC", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VE", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VG", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VI", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VN", "VNM"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0336012403",
			"+84586012403",
			"84981577798",
			"0708001240",
			"84813601243",
			"0523803765",
			"0863803732",
			"0883805866",
			"0892405867",
			"+84888696413",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"260976684590",
			"01678912345",
			"+841698765432",
			"841626543219",
			"0533803765",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"VU", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"WF", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"WS", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"YE", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"YT", ""}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ZA", "ZAF"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0821231234",
			"+27821231234",
			"27821231234",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"12345",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ZM", "ZMB"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"0956684590",
			"0966684590",
			"0976684590",
			"+260956684590",
			"+260966684590",
			"+260976684590",
			"260976684590",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"966684590",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
	Run(t, []string{"ZW", "ZWE"}, List{{
		Name: "Phone", Func: valid.Phone,
		Pass: []string{
			"+263561890123",
			"+263715558041",
			"+263775551112",
			"+263775551695",
			"+263715556633",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+2631234567890",
			"+2641234567",
			"+263981234",
			"4736338855",
			"66338855",
		},
	}, {
		Name: "Zip", Func: valid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
}
