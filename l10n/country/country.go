package country

import (
	"regexp"
	"strconv"
	"strings"
)

// References:
// - https://wikipedia.org/wiki/ISO_3166-1_alpha-2
// - https://en.wikipedia.org/wiki/List_of_postal_codes
// - https://en.youbianku.com
// - https://en.wikipedia.org/wiki/VAT_identification_number
// ...

var ISO31661A_2 = make(map[string]Country)
var ISO31661A_3 = make(map[string]Country)

func init() {
	// NOTE only map's keys are populated here, the fields of the Country
	// value are populated individually by the country/internal/xx packages.

	for _, a2 := range alpha2 {
		ISO31661A_2[a2] = Country{}
	}
	for _, a3 := range alpha3 {
		ISO31661A_3[a3] = Country{}
	}
}

var alpha2 = []string{"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR",
	"AS", "AT", "AU", "AW", "AX", "AZ", "BA", "BB", "BD", "BE", "BF", "BG",
	"BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS", "BT", "BV",
	"BW", "BY", "BZ", "CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL",
	"CM", "CN", "CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ", "DE", "DJ",
	"DK", "DM", "DO", "DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI",
	"FJ", "FK", "FM", "FO", "FR", "GA", "GB", "GD", "GE", "GF", "GG", "GH",
	"GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY",
	"HK", "HM", "HN", "HR", "HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO",
	"IQ", "IR", "IS", "IT", "JE", "JM", "JO", "JP", "KE", "KG", "KH", "KI",
	"KM", "KN", "KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC", "LI", "LK",
	"LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "ME", "MF", "MG",
	"MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU",
	"MV", "MW", "MX", "MY", "MZ", "NA", "NC", "NE", "NF", "NG", "NI", "NL",
	"NO", "NP", "NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK",
	"PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS",
	"RU", "RW", "SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK",
	"SL", "SM", "SN", "SO", "SR", "SS", "ST", "SV", "SX", "SY", "SZ", "TC",
	"TD", "TF", "TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO", "TR", "TT",
	"TV", "TW", "TZ", "UA", "UG", "UM", "US", "UY", "UZ", "VA", "VC", "VE",
	"VG", "VI", "VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM", "ZW",
}

var alpha3 = []string{"ABW", "AFG", "AGO", "AIA", "ALA", "ALB", "AND", "ARE", "ARG",
	"ARM", "ASM", "ATA", "ATF", "ATG", "AUS", "AUT", "AZE", "BDI", "BEL", "BEN",
	"BES", "BFA", "BGD", "BGR", "BHR", "BHS", "BIH", "BLM", "BLR", "BLZ", "BMU",
	"BOL", "BRA", "BRB", "BRN", "BTN", "BVT", "BWA", "CAF", "CAN", "CCK", "CHE",
	"CHL", "CHN", "CIV", "CMR", "COD", "COG", "COK", "COL", "COM", "CPV", "CRI",
	"CUB", "CUW", "CXR", "CYM", "CYP", "CZE", "DEU", "DJI", "DMA", "DNK", "DOM",
	"DZA", "ECU", "EGY", "ERI", "ESH", "ESP", "EST", "ETH", "FIN", "FJI", "FLK",
	"FRA", "FRO", "FSM", "GAB", "GBR", "GEO", "GGY", "GHA", "GIB", "GIN", "GLP",
	"GMB", "GNB", "GNQ", "GRC", "GRD", "GRL", "GTM", "GUF", "GUM", "GUY", "HKG",
	"HMD", "HND", "HRV", "HTI", "HUN", "IDN", "IMN", "IND", "IOT", "IRL", "IRN",
	"IRQ", "ISL", "ISR", "ITA", "JAM", "JEY", "JOR", "JPN", "KAZ", "KEN", "KGZ",
	"KHM", "KIR", "KNA", "KOR", "KWT", "LAO", "LBN", "LBR", "LBY", "LCA", "LIE",
	"LKA", "LSO", "LTU", "LUX", "LVA", "MAC", "MAF", "MAR", "MCO", "MDA", "MDG",
	"MDV", "MEX", "MHL", "MKD", "MLI", "MLT", "MMR", "MNE", "MNG", "MNP", "MOZ",
	"MRT", "MSR", "MTQ", "MUS", "MWI", "MYS", "MYT", "NAM", "NCL", "NER", "NFK",
	"NGA", "NIC", "NIU", "NLD", "NOR", "NPL", "NRU", "NZL", "OMN", "PAK", "PAN",
	"PCN", "PER", "PHL", "PLW", "PNG", "POL", "PRI", "PRK", "PRT", "PRY", "PSE",
	"PYF", "QAT", "REU", "ROU", "RUS", "RWA", "SAU", "SDN", "SEN", "SGP", "SGS",
	"SHN", "SJM", "SLB", "SLE", "SLV", "SMR", "SOM", "SPM", "SRB", "SSD", "STP",
	"SUR", "SVK", "SVN", "SWE", "SWZ", "SXM", "SYC", "SYR", "TCA", "TCD", "TGO",
	"THA", "TJK", "TKL", "TKM", "TLS", "TON", "TTO", "TUN", "TUR", "TUV", "TWN",
	"TZA", "UGA", "UKR", "UMI", "URY", "USA", "UZB", "VAT", "VCT", "VEN", "VGB",
	"VIR", "VNM", "VUT", "WLF", "WSM", "YEM", "ZAF", "ZMB", "ZWE",
}

type Country struct {
	// ISO 3166-1 Alpha-2
	A2 string
	// ISO 3166-1 Alpha-3
	A3 string
	// ISO 3166-1 numeric
	Num string
	// The validator for the country's zip / postal code, will be nil
	// for countries that don't use postal codes.
	Zip StringMatcher
	// The validator for the country's phone numbers, may be nil.
	Phone StringMatcher
	//
	VAT StringMatcher
}

func Add(c Country) {
	ISO31661A_2[c.A2] = c
	ISO31661A_3[c.A3] = c
}

func Get(cc string) (c Country, ok bool) {
	if len(cc) == 2 {
		cc = strings.ToUpper(cc)
		c, ok = ISO31661A_2[cc]
		return c, ok
	}
	if len(cc) == 3 {
		cc = strings.ToUpper(cc)
		c, ok = ISO31661A_3[cc]
		return c, ok
	}
	return c, false
}

type StringMatcher interface {
	MatchString(v string) bool
}

type StringMatcherFunc func(v string) bool

func (f StringMatcherFunc) MatchString(v string) bool {
	return f(v)
}

var (
	RxZip3Digits = regexp.MustCompile(`^[0-9]{3}$`)
	RxZip4Digits = regexp.MustCompile(`^[0-9]{4}$`)
	RxZip5Digits = regexp.MustCompile(`^[0-9]{5}$`)
	RxZip6Digits = regexp.MustCompile(`^[0-9]{6}$`)
)

// convenience func for byte known to represent a digit
func btoi(v byte) int {
	i, _ := strconv.Atoi(string(v))
	return i
}

// ISO 7064 (MOD 11, 10)
func ISO7064_MOD11_10(v string) bool {
	check := 10
	for i := 0; i < len(v)-1; i++ {
		if mod := (btoi(v[i]) + check) % 10; mod == 0 {
			check = (10 * 2) % 11
		} else {
			check = (mod * 2) % 11
		}
	}

	if check != 1 {
		return btoi(v[10]) == (11 - check)
	}
	return btoi(v[10]) == 0
}
