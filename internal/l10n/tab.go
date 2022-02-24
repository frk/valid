package l10n

import (
	"regexp"
)

var ISO31661A_2 = make(map[string]Country)
var ISO31661A_3 = make(map[string]Country)

func init() {
	for _, c := range cctab {
		ISO31661A_2[c.A2] = c
		ISO31661A_3[c.A3] = c
	}
}

var (
	// common zip patterns
	rxZip3Digits = regexp.MustCompile(`^[0-9]{3}$`)
	rxZip4Digits = regexp.MustCompile(`^[0-9]{4}$`)
	rxZip5Digits = regexp.MustCompile(`^[0-9]{5}$`)
	rxZip6Digits = regexp.MustCompile(`^[0-9]{6}$`)
)

var cctab = []Country{{
	A2: "AD", A3: "AND", Num: "020",
	// NOTE(mkopriva): For the "AD5XX" post codes any digit between 0-9 is
	// allowed in the place of the Xs because I can't find anything substantial
	// to better handle the following: "PO Boxes in Andorra la Vella have separate
	// postcodes allocated to each group of 50 boxes - e.g., boxes 1001 to 1050
	// have a code of AD551, 1051 to 1100 a code of AD552 etc."
	// (from: https://en.wikipedia.org/wiki/Postal_codes_in_Andorra)
	Zip:   regexp.MustCompile(`^AD(?:[1-7]00|5[0-9]{2})$`),
	Phone: regexp.MustCompile(`^(?:\+376)?[346][0-9]{5}$`),
}, {
	A2: "AE", A3: "ARE", Num: "784",
	Phone: regexp.MustCompile(`^(?:(?:\+?971)|0)?5[024568][0-9]{7}$`),
}, {
	A2: "AF", A3: "AFG", Num: "004",
	Zip: regexp.MustCompile(`^(?:[1-3][0-9]|4[0-3])(?:[0-9][1-9])$`),
}, {
	A2: "AG", A3: "ATG", Num: "028",
}, {
	A2: "AI", A3: "AIA", Num: "660",
	Zip: regexp.MustCompile(`^AI-2640$`),
}, {
	A2: "AL", A3: "ALB", Num: "008", Zip: rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+355|0)6[789][0-9]{6}$`),
	// 10 characters, the first position following the prefix
	// is "J" or "K" or "L", and the last character is a letter
	// – e.g. K99999999L or L99999999G
	VAT: regexp.MustCompile(`^[JKL][0-9]{8}[A-Z]$`),
}, {
	A2: "AM", A3: "ARM", Num: "051", Zip: rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?374|0)(?:(?:10|[9|7][0-9])[0-9]{6}|[2-4][0-9]{7})$`),
}, {
	A2: "AO", A3: "AGO", Num: "024",
	Phone: regexp.MustCompile(`^(?:\+244)[0-9]{9}$`),
}, {
	A2: "AQ", A3: "ATA", Num: "010",
	Zip: regexp.MustCompile(`^BIQQ 1ZZ$`),
}, {
	A2: "AR", A3: "ARG", Num: "032",
	Zip:   regexp.MustCompile(`^(?:[0-9]{4})|(?:[A-Z][0-9]{4}[A-Z]{3})$`),
	Phone: regexp.MustCompile(`^\+?549(?:11|[2368][0-9])[0-9]{8}$`),
	VAT:   regexp.MustCompile(`^[0-9]{11}$`),
}, {
	A2: "AS", A3: "ASM", Num: "016",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-?[0-9]{4})?$`),
}, {
	A2: "AT", A3: "AUT", Num: "040",
	Zip:   rxZip3Digits,
	Phone: regexp.MustCompile(`^(?:\+43|0)[0-9]{1,4}[0-9]{3,12}$`),
	// 'AT'+U+8 digits, – e.g. ATU99999999
	VAT: regexp.MustCompile(`^ATU[0-9]{8}$`),
}, {
	A2: "AU", A3: "AUS", Num: "036", Zip: rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?61|0)4[0-9]{8}$`),
	VAT:   StringMatcherFunc(auVAT),
}, {
	A2: "AW", A3: "ABW", Num: "533",
}, {
	A2: "AX", A3: "ALA", Num: "248",
	Zip: regexp.MustCompile(`^(?:AX-)?[0-9]{5}$`),
}, {
	A2: "AZ", A3: "AZE", Num: "031",
	Zip:   regexp.MustCompile(`^AZ[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+994|0)(?:5[015]|7[07]|99)[0-9]{7}$`),
}, {
	A2: "BA", A3: "BIH", Num: "070",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:(?:\+|00)3876)|06)(?:(?:(?:[0-3]|[5-6])[0-9]{6})|4[0-9]{7})$`),
}, {
	A2: "BB", A3: "BRB", Num: "052",
	Zip: regexp.MustCompile(`^BB[0-9]{5}$`),
}, {
	A2: "BD", A3: "BGD", Num: "050",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?880|0)1[13456789][0-9]{8}$`),
}, {
	A2: "BE", A3: "BEL", Num: "056",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?32|0)4?[0-9]{8}$`),
	VAT:   StringMatcherFunc(beVAT),
}, {
	A2: "BF", A3: "BFA", Num: "854",
}, {
	A2: "BG", A3: "BGR", Num: "100",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?359|0)?8[789][0-9]{7}$`),
	VAT:   regexp.MustCompile(`^BG[0-9]{9,10}$`),
}, {
	A2: "BH", A3: "BHR", Num: "048",
	Zip:   regexp.MustCompile(`^[0-9]{3,4}$`),
	Phone: regexp.MustCompile(`^(?:\+?973)?(?:3|6)[0-9]{7}$`),
}, {
	A2: "BI", A3: "BDI", Num: "108",
}, {
	A2: "BJ", A3: "BEN", Num: "204",
}, {
	A2: "BL", A3: "BLM", Num: "652",
	Zip: regexp.MustCompile(`^97133$`),
}, {
	A2: "BM", A3: "BMU", Num: "060",
}, {
	A2: "BN", A3: "BRN", Num: "096",
	Zip: regexp.MustCompile(`^[A-Z]{2}[0-9]{4}$`),
}, {
	A2: "BO", A3: "BOL", Num: "068",
	Phone: regexp.MustCompile(`^(?:\+?591)?(?:6|7)[0-9]{7}$`),
	VAT:   regexp.MustCompile(`^[0-9]{7}$`),
}, {
	A2: "BQ", A3: "BES", Num: "535",
}, {
	A2: "BR", A3: "BRA", Num: "076",
	Zip:   regexp.MustCompile(`^[0-9]{5}-[0-9]{3}$`),
	Phone: regexp.MustCompile(`^(?:(?:\+?55[ ]?[1-9]{2}[ ]?)|(?:\+?55[ ]?\([1-9]{2}\)[ ]?)|(?:0[1-9]{2}[ ]?)|(?:\([1-9]{2}\)[ ]?)|(?:[1-9]{2}[ ]?))(?:(?:[0-9]{4}-?[0-9]{4})|(?:9[2-9]{1}[0-9]{3}-?[0-9]{4}))$`),
}, {
	A2: "BS", A3: "BHS", Num: "044",
}, {
	A2: "BT", A3: "BTN", Num: "064",
	Zip: rxZip5Digits,
}, {
	A2: "BV", A3: "BVT", Num: "074",
}, {
	A2: "BW", A3: "BWA", Num: "072",
}, {
	A2: "BY", A3: "BLR", Num: "112",
	Zip:   regexp.MustCompile(`^2[1-4]{1}[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+?375)?(?:24|25|29|33|44)[0-9]{7}$`),
	// 9 digit number
	VAT: regexp.MustCompile(`^(?:УНП[ ]?)?[0-9]{9}$`),
}, {
	A2: "BZ", A3: "BLZ", Num: "084",
}, {
	A2: "CA", A3: "CAN", Num: "124",
	Zip:   regexp.MustCompile(`^(?i)[ABCEGHJKLMNPRSTVXY][0-9][ABCEGHJ-NPRSTV-Z][\s\-]?[0-9][ABCEGHJ-NPRSTV-Z][0-9]$`),
	Phone: regexp.MustCompile(`^(?:(?:\+1|1)?(?: |-)?)?(?:\([2-9][0-9]{2}\)|[2-9][0-9]{2})(?: |-)?(?:[2-9][0-9]{2}(?: |-)?[0-9]{4})$`),
	// 9 digit number (same as BN or GST/HST number)
	VAT: regexp.MustCompile(`^[0-9]{9}$`),
}, {
	A2: "CC", A3: "CCK", Num: "166",
	Zip: rxZip4Digits,
}, {
	A2: "CD", A3: "COD", Num: "180",
}, {
	A2: "CF", A3: "CAF", Num: "140",
}, {
	A2: "CG", A3: "COG", Num: "178",
}, {
	A2: "CH", A3: "CHE", Num: "756",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+41|0)7[5-9][0-9]{1,7}$`),
	VAT:   StringMatcherFunc(chVAT),
}, {
	A2: "CI", A3: "CIV", Num: "384",
}, {
	A2: "CK", A3: "COK", Num: "184",
}, {
	A2: "CL", A3: "CHL", Num: "152",
	Zip:   regexp.MustCompile(`^[0-9]{3}-?[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+?56|0)[2-9][0-9]{8}$`),
}, {
	A2: "CM", A3: "CMR", Num: "120",
}, {
	A2: "CN", A3: "CHN", Num: "156",
	Zip:   regexp.MustCompile(`^(?:0[1-7]|1[012356]|2[0-7]|3[0-6]|4[0-7]|5[1-7]|6[1-7]|7[1-5]|8[1345]|9[09])[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:(?:\+|00)86)?1(?:[3568][0-9]|4[579]|6[67]|7[01235678]|9[012356789])[0-9]{8}$`),
}, {
	A2: "CO", A3: "COL", Num: "170",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?57)?(?:[1-8]{1}|3[0-9]{2})?[0-9]{7}$`),
}, {
	A2: "CR", A3: "CRI", Num: "188",
	Zip:   regexp.MustCompile(`^[0-9]{5}(-[0-9]{4})?$`),
	Phone: regexp.MustCompile(`^(?:\+506)?[2-8][0-9]{7}$`),
}, {
	A2: "CU", A3: "CUB", Num: "192",
	Zip: rxZip5Digits,
}, {
	A2: "CV", A3: "CPV", Num: "132",
	Zip: rxZip4Digits,
}, {
	A2: "CW", A3: "CUW", Num: "531",
}, {
	A2: "CX", A3: "CXR", Num: "162",
	Zip: rxZip4Digits,
}, {
	A2: "CY", A3: "CYP", Num: "196",
	Zip: regexp.MustCompile(`^[0-9]{4,5}$`),
	// 9 characters, last one must be a letter – e.g. CY99999999L
	VAT: regexp.MustCompile(`^CY[0-9]{8}[A-Z]$`),
}, {
	A2: "CZ", A3: "CZE", Num: "203",
	Zip:   regexp.MustCompile(`^[0-9]{3}[ ]?[0-9]{2}$`),
	Phone: regexp.MustCompile(`^(?:\+?420)?[ ]?[1-9][0-9]{2}[ ]?[0-9]{3}[ ]?[0-9]{3}$`),
	// 8, 9 or 10 characters -- i.e. CZ12345678, CZ123456789, CZ1234567890
	VAT: regexp.MustCompile(`^CZ[0-9]{8,10}$`),
}, {
	A2: "DE", A3: "DEU", Num: "276",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+49)?0?[1|3](?:[0|5][0-9]{2}|6(?:[23]|0[0-9]?)|7(?:[0-57-9]|6[0-9]))[0-9]{7}$`),
	// 9 digits, e.g. DE999999999
	VAT: regexp.MustCompile(`^DE[0-9]{9}$`),
}, {
	A2: "DJ", A3: "DJI", Num: "262",
}, {
	A2: "DK", A3: "DNK", Num: "208",
	Zip:   regexp.MustCompile(`^(?:DK-)?[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+?45)?(?:[ ]?[0-9]{2}){4}$`),
	VAT:   StringMatcherFunc(dkVAT),
}, {
	A2: "DM", A3: "DMA", Num: "212",
}, {
	A2: "DO", A3: "DOM", Num: "214",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?1)?8[024]9[0-9]{7}$`),
	VAT:   regexp.MustCompile(`^[0-9]{9}|[0-9]{11}$`),
}, {
	A2: "DZ", A3: "DZA", Num: "012",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?213|0)(?:5|6|7)[0-9]{8}$`),
}, {
	A2: "EC", A3: "ECU", Num: "218",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?593|0)(?:[2-7]|9[2-9])[0-9]{7}$`),
	VAT:   regexp.MustCompile(`^[0-9]{13}$`),
}, {
	A2: "EE", A3: "EST", Num: "233",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?372)?[ ]?(?:5|8[1-4])[ ]?(?:[0-9][ ]?){6,7}$`),
	VAT:   regexp.MustCompile(`^EE[0-9]{9}$`),
}, {
	A2: "EG", A3: "EGY", Num: "818",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+?20)|0)?1[0125][0-9]{8}$`),
}, {
	A2: "EH", A3: "ESH", Num: "732",
}, {
	A2: "ER", A3: "ERI", Num: "232",
}, {
	A2: "ES", A3: "ESP", Num: "724",
	Zip:   regexp.MustCompile(`^(?:5[0-2]{1}|[0-4]{1}[0-9]{1})[0-9]{3}$`),
	Phone: regexp.MustCompile(`^(?:\+?34)?[6|7][0-9]{8}$`),
	// 'ES'+letter+8 digits; or 'ES'+letter+7 digits+letter; or 'ES'+8 digits+letter
	VAT: regexp.MustCompile(`^ES(?:[A-Z][0-9]{8}|[A-Z][0-9]{7}[A-Z]|[0-9]{8}[A-Z])$`),
}, {
	A2: "ET", A3: "ETH", Num: "231",
	Zip: rxZip4Digits,
}, {
	A2: "FI", A3: "FIN", Num: "246",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?358|0)[ ]?(?:4(?:0|1|2|4|5|6)?|50)[ ]?(?:[0-9][ ]?){4,8}[0-9]$`),
	VAT:   StringMatcherFunc(fiVAT),
}, {
	A2: "FJ", A3: "FJI", Num: "242",
	Phone: regexp.MustCompile(`^(?:\+?679)?[ ]?[0-9]{3}[ ]?[0-9]{4}$`),
}, {
	A2: "FK", A3: "FLK", Num: "238",
	Zip: regexp.MustCompile(`^FIQQ 1ZZ$`),
}, {
	A2: "FM", A3: "FSM", Num: "583",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "FO", A3: "FRO", Num: "234",
	Zip:   regexp.MustCompile(`^FO-[0-9]{3}$`),
	Phone: regexp.MustCompile(`^(?:\+?298)?(?:[ ]?[0-9]{2}){3}$`),
}, {
	A2: "FR", A3: "FRA", Num: "250",
	Zip:   regexp.MustCompile(`^[0-9]{2}\s?[0-9]{3}$`),
	Phone: regexp.MustCompile(`^(?:\+?33|0)[67][0-9]{8}$`),
	VAT:   StringMatcherFunc(frVAT),
}, {
	A2: "GA", A3: "GAB", Num: "266",
}, {
	A2: "GB", A3: "GBR", Num: "826",
	Zip:   regexp.MustCompile(`^(?i)(?:gir\s?0aa|[a-z]{1,2}[0-9][0-9a-z]?\s?(?:[0-9][a-z]{2})?)$`),
	Phone: regexp.MustCompile(`^(?:\+?44|0)7[0-9]{9}$`),
}, {
	A2: "GD", A3: "GRD", Num: "308",
}, {
	A2: "GE", A3: "GEO", Num: "268",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?995)?(?:5|79)[0-9]{7}$`),
}, {
	A2: "GF", A3: "GUF", Num: "254",
	Zip:   regexp.MustCompile(`^973(?:[0-8][0-9]|90)$`),
	Phone: regexp.MustCompile(`^(?:\+?594|0|00594)[67][0-9]{8}$`),
}, {
	A2: "GG", A3: "GGY", Num: "831",
	Zip:   regexp.MustCompile(`^GY[0-9]{1,2} [0-9][A-Z]{2}$`),
	Phone: regexp.MustCompile(`^(?:\+?44|0)1481[0-9]{6}$`),
}, {
	A2: "GH", A3: "GHA", Num: "288",
	Zip:   regexp.MustCompile(`^[A-Z][A-Z0-9]-[0-9]{4}-[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+233|0)(?:20|50|24|54|27|57|26|56|23|28)[0-9]{7}$`),
}, {
	A2: "GI", A3: "GIB", Num: "292",
	Zip: regexp.MustCompile(`^GX11 1AA$`),
}, {
	A2: "GL", A3: "GRL", Num: "304",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?299)?(?:[ ]?[0-9]{2}){3}$`),
}, {
	A2: "GM", A3: "GMB", Num: "270",
}, {
	A2: "GN", A3: "GIN", Num: "324", Zip: rxZip3Digits,
}, {
	A2: "GP", A3: "GLP", Num: "312",
	Zip:   regexp.MustCompile(`^971(?:[0-8][0-9]|90)$`),
	Phone: regexp.MustCompile(`^(?:\+?590|0|00590)[67][0-9]{8}$`),
}, {
	A2: "GQ", A3: "GNQ", Num: "226",
}, {
	A2: "GR", A3: "GRC", Num: "300",
	Zip:   regexp.MustCompile(`^[0-9]{3}[ ]?[0-9]{2}$`),
	Phone: regexp.MustCompile(`^(?:\+?30|0)?(?:69[0-9]{8})$`),
	VAT:   StringMatcherFunc(grVAT),
}, {
	A2: "GS", A3: "SGS", Num: "239",
	Zip: regexp.MustCompile(`^SIQQ 1ZZ$`),
}, {
	A2: "GT", A3: "GTM", Num: "320",
	Zip: rxZip5Digits,
	// seven digits, one dash (-); one digit (like 1234567-1)
	VAT: regexp.MustCompile(`^[0-9]{7}-[0-9]$`),
}, {
	A2: "GU", A3: "GUM", Num: "316",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "GW", A3: "GNB", Num: "624", Zip: rxZip4Digits,
}, {
	A2: "GY", A3: "GUY", Num: "328",
}, {
	A2: "HK", A3: "HKG", Num: "344",
	Phone: regexp.MustCompile(`^(?:\+?852[\- ]?)?[456789][0-9]{3}[\- ]?[0-9]{4}$`),
}, {
	A2: "HM", A3: "HMD", Num: "334",
}, {
	A2: "HN", A3: "HND", Num: "340",
	Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:[A-Z]{2}[0-9]{4})$`),
	Phone: regexp.MustCompile(`^(?:\+?504)?[9|8][0-9]{7}$`),
}, {
	A2: "HR", A3: "HRV", Num: "191",
	Zip: regexp.MustCompile(`^(?:[1-5][0-9]{4}$)`),
	VAT: StringMatcherFunc(hrVAT),
}, {
	A2: "HT", A3: "HTI", Num: "332",
	Zip: regexp.MustCompile(`^HT[0-9]{4}$`),
}, {
	A2: "HU", A3: "HUN", Num: "348",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?36)(?:20|30|70)[0-9]{7}$`),
	// 8 digits (the first 8 digits of the national tax number) – e.g. HU12345678
	VAT: regexp.MustCompile(`^HU[0-9]{8}$`),
}, {
	A2: "ID", A3: "IDN", Num: "360",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?62|0)8(?:1[123456789]|2[1238]|3[1238]|5[12356789]|7[78]|9[56789]|8[123456789])[ ?|0-9]{5,11}$`),
	// 15 digit number (ex. 02.271.824.1-413.000)
	VAT: regexp.MustCompile(`^(?:[0-9]{15})|(?:[0-9]{2}.[0-9]{3}.[0-9]{3}.[0-9][\-–][0-9]{3}.[0-9]{3})$`),
}, {
	A2: "IE", A3: "IRL", Num: "372",
	// References:
	// - https://stackoverflow.com/questions/33391412/validation-for-irish-eircode
	// - https://www.eircode.ie/docs/default-source/Common/prepareyourbusinessforeircode-edition3published.pdf
	// - https://en.wikipedia.org/wiki/Postal_addresses_in_the_Republic_of_Ireland
	Zip:   regexp.MustCompile(`^(?:[AC-FHKNPRTV-Y][0-9]{2}|D6W)[ -]?[0-9AC-FHKNPRTV-Y]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+?353|0)8[356789][0-9]{7}$`),
	// 'IE'+7 digits and one letter, optionally followed by a 'W' for married women, e.g. IE1234567T or IE1234567TW
	// or 'IE'+7 digits and two letters, e.g. IE1234567FA (since January 2013)
	// or 'IE'+one digit, one letter/"+"/"*", 5 digits and one letter (old style, currently being phased out)
	VAT: regexp.MustCompile(`^IE(?:[0-9]{7}[A-Z]{1,2}|[0-9][A-Z+*][0-9]{5}[A-Z])$`),
}, {
	A2: "IL", A3: "ISR", Num: "376",
	Zip:   regexp.MustCompile(`^(?:[0-9]{5}|[0-9]{7})$`),
	Phone: regexp.MustCompile(`^(?:\+972|0)(?:[23489]|5[012345689]|77)[1-9][0-9]{6}$`),
	VAT:   StringMatcherFunc(ilVAT),
}, {
	A2: "IM", A3: "IMN", Num: "833",
	Zip: regexp.MustCompile(`^IM[0-9]{1,2} [0-9][A-Z]{2}$`),
}, {
	A2: "IN", A3: "IND", Num: "356",
	// References:
	// - https://en.wikipedia.org/wiki/Postal_Index_Number
	// - https://en.youbianku.com/India
	Zip:   StringMatcherFunc(inZIP),
	Phone: regexp.MustCompile(`^(?:\+?91|0)?[6789][0-9]{9}$`),
}, {
	A2: "IO", A3: "IOT", Num: "086",
	Zip: regexp.MustCompile(`^BBND 1ZZ$`),
}, {
	A2: "IQ", A3: "IRQ", Num: "368",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?964|0)?7[0-9]{9}$`),
}, {
	A2: "IR", A3: "IRN", Num: "364",
	Zip:   regexp.MustCompile(`^[0-9]{10}$`),
	Phone: regexp.MustCompile(`^(?:\+?98[\- ]?|0)9[0-39][0-9][\- ]?[0-9]{3}[\- ]?[0-9]{4}$`),
}, {
	A2: "IS", A3: "ISL", Num: "352",
	Zip: rxZip3Digits,
	// 5 or 6 characters depending on age of the company
	VAT: regexp.MustCompile(`^[0-9]{5,6}$`),
}, {
	A2: "IT", A3: "ITA", Num: "380",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?39)?[ ]?3[0-9]{2}[ ]?[0-9]{6,7}$`),
	VAT:   StringMatcherFunc(itVAT),
}, {
	A2: "JE", A3: "JEY", Num: "832",
	Zip: regexp.MustCompile(`^JE[0-9]{1,2} [0-9][A-Z]{2}$`),
}, {
	A2: "JM", A3: "JAM", Num: "388",
	Zip: regexp.MustCompile(`^[1-9]|1[0-9]|20$`),
}, {
	A2: "JO", A3: "JOR", Num: "400",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?962|0)?7[789][0-9]{7}$`),
}, {
	A2: "JP", A3: "JPN", Num: "392",
	Zip:   regexp.MustCompile(`^[0-9]{3}\-[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+81[ \-]?(?:\(0\))?|0)[6789]0(?:[ \-]?[0-9]{4}){2}$`),
}, {
	A2: "KE", A3: "KEN", Num: "404",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?254|0)(?:7|1)[0-9]{8}$`),
}, {
	A2: "KG", A3: "KGZ", Num: "417",
	Zip: rxZip6Digits,
}, {
	A2: "KH", A3: "KHM", Num: "116",
	Zip: rxZip6Digits,
}, {
	A2: "KI", A3: "KIR", Num: "296",
}, {
	A2: "KM", A3: "COM", Num: "174",
}, {
	A2: "KN", A3: "KNA", Num: "659",
}, {
	A2: "KP", A3: "PRK", Num: "408",
}, {
	A2: "KR", A3: "KOR", Num: "410",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+?82)[ \-]?)?0?1(?:[0|1|6|7|8|9]{1})[ \-]?[0-9]{3,4}[ \-]?[0-9]{4}$`),
}, {
	A2: "KW", A3: "KWT", Num: "414",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?965)[569][0-9]{7}$`),
}, {
	A2: "KY", A3: "CYM", Num: "136",
	Zip: regexp.MustCompile(`^KY[0-9]-[0-9]{4}$`),
}, {
	A2: "KZ", A3: "KAZ", Num: "398",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?7|8)?7[0-9]{9}$`),
	// 12 digits
	VAT: regexp.MustCompile(`^[0-9]{12}$`),
}, {
	A2: "LA", A3: "LAO", Num: "418",
	Zip: rxZip5Digits,
}, {
	A2: "LB", A3: "LBN", Num: "422",
	Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:[0-9]{4} [0-9]{4})$`),
	Phone: regexp.MustCompile(`^(?:\+?961)?(?:(?:3|81)[0-9]{6}|7[0-9]{7})$`),
}, {
	A2: "LC", A3: "LCA", Num: "662",
	// Reference:
	// - https://stluciapostal.com/postal-codes-2/
	Zip: regexp.MustCompile(`^LC[0-9]{2}[ ]{0,2}[0-9]{3}$`),
}, {
	A2: "LI", A3: "LIE", Num: "438",
	Zip: regexp.MustCompile(`^(?:948[5-9]|949[0-7])$`),
}, {
	A2: "LK", A3: "LKA", Num: "144",
	Zip: rxZip5Digits,
}, {
	A2: "LR", A3: "LBR", Num: "430",
	Zip: rxZip4Digits,
}, {
	A2: "LS", A3: "LSO", Num: "426",
	Zip: regexp.MustCompile(`^[0-9]{3}$`),
}, {
	A2: "LT", A3: "LTU", Num: "440",
	Zip:   regexp.MustCompile(`^LT\-[0-9]{5}$`),
	Phone: regexp.MustCompile(`^(?:\+370|8)[0-9]{8}$`),
	VAT:   regexp.MustCompile(`^LT[0-9]{9}(?:[0-9]{3})?$`),
}, {
	A2: "LU", A3: "LUX", Num: "442",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+352)?(?:(?:6[0-9]1)[0-9]{6})$`),
	VAT:   regexp.MustCompile(`^LU[0-9]{8}$`),
}, {
	A2: "LV", A3: "LVA", Num: "428",
	Zip: regexp.MustCompile(`^LV\-[0-9]{4}$`),
	VAT: regexp.MustCompile(`^LV[0-9]{11}$`),
}, {
	A2: "LY", A3: "LBY", Num: "434",
	Phone: regexp.MustCompile(`^(?:(?:\+?218)|0)?(?:9[1-6][0-9]{7}|[1-8][0-9]{7,9})$`),
}, {
	A2: "MA", A3: "MAR", Num: "504",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+|00)212|0)[5-7][0-9]{8}$`),
}, {
	A2: "MC", A3: "MCO", Num: "492",
	Zip: regexp.MustCompile(`^MC980(?:[0-9]{2})$`),
	VAT: StringMatcherFunc(frVAT),
}, {
	A2: "MD", A3: "MDA", Num: "498",
	Zip: regexp.MustCompile(`^MD-?[0-9]{4}$`),
}, {
	A2: "ME", A3: "MNE", Num: "499",
	Zip: rxZip5Digits,
}, {
	A2: "MF", A3: "MAF", Num: "663",
	Zip: regexp.MustCompile(`^97150$`),
}, {
	A2: "MG", A3: "MDG", Num: "450",
	Zip: regexp.MustCompile(`^[0-9]{3}$`),
}, {
	A2: "MH", A3: "MHL", Num: "584",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "MK", A3: "MKD", Num: "807",
	Zip: rxZip4Digits,
	// 15 characters, the first two positions are for the prefix
	// "MK", followed by 13 numbers – e.g. MK4032013544513
	VAT: regexp.MustCompile(`^MK[0-9]{13}$`),
}, {
	A2: "ML", A3: "MLI", Num: "466",
}, {
	A2: "MM", A3: "MMR", Num: "104",
	Zip: rxZip5Digits,
}, {
	A2: "MN", A3: "MNG", Num: "496",
	Zip: rxZip5Digits,
}, {
	A2: "MO", A3: "MAC", Num: "446",
	Phone: regexp.MustCompile(`^(?:\+?853[\- ]?)?[6][0-9]{3}[\- ]?[0-9]{4}$`),
}, {
	A2: "MP", A3: "MNP", Num: "580",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "MQ", A3: "MTQ", Num: "474",
	Zip:   regexp.MustCompile(`^972(?:[0-8][0-9]|90)$`),
	Phone: regexp.MustCompile(`^(?:\+?596|0|00596)[67][0-9]{8}$`),
}, {
	A2: "MR", A3: "MRT", Num: "478",
}, {
	A2: "MS", A3: "MSR", Num: "500",
	Zip: regexp.MustCompile(`^MSR 1[1-3][0-9]{2}$`),
}, {
	A2: "MT", A3: "MLT", Num: "470",
	Zip:   regexp.MustCompile(`^(?i)[a-z]{3}\s{0,1}[0-9]{4}$`),
	Phone: regexp.MustCompile(`^(?:\+?356|0)?(?:99|79|77|21|27|22|25)[0-9]{6}$`),
	VAT:   regexp.MustCompile(`^MT[0-9]{8}$`),
}, {
	A2: "MU", A3: "MUS", Num: "480",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?230|0)?[0-9]{8}$`),
}, {
	A2: "MV", A3: "MDV", Num: "462",
	Zip: rxZip5Digits,
}, {
	A2: "MW", A3: "MWI", Num: "454",
}, {
	A2: "MX", A3: "MEX", Num: "484",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?52)?(?:1|01)?[0-9]{10,11}$`),
}, {
	A2: "MY", A3: "MYS", Num: "458",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?6?01){1}(?:(?:[0145]{1}(?:-| )?[0-9]{7,8})|(?:[236789]{1}(?:-| )?[0-9]{7}))$`),
}, {
	A2: "MZ", A3: "MOZ", Num: "508",
	Zip: rxZip4Digits,
}, {
	A2: "NA", A3: "NAM", Num: "516",
}, {
	A2: "NC", A3: "NCL", Num: "540",
	Zip: regexp.MustCompile(`^988(?:[0-8][0-9]|90)$`),
}, {
	A2: "NE", A3: "NER", Num: "562",
	Zip: rxZip4Digits,
}, {
	A2: "NF", A3: "NFK", Num: "574",
	Zip: rxZip4Digits,
}, {
	A2: "NG", A3: "NGA", Num: "566",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?234|0)?[789][0-9]{9}$`),
	// 12 digits in the format 01012345-0001
	VAT: regexp.MustCompile(`^[0-9]{8}-[0-9]{4}$`),
}, {
	A2: "NI", A3: "NIC", Num: "558",
	Zip: rxZip5Digits,
	// 3 digits, 1 dash, 6 digits, 1 dash, 4 digits followed by 1 letter
	VAT: regexp.MustCompile(`^[0-9]{3}-[0-9]{6}-[0-9]{4}[A-Z]$`),
}, {
	A2: "NL", A3: "NLD", Num: "528",
	Zip:   regexp.MustCompile(`^(?i)[0-9]{4}\s?[a-z]{2}$`),
	Phone: regexp.MustCompile(`^(?:(?:(?:\+|00)?31\(0\))|(?:(?:\+|00)?31)|0)6{1}[0-9]{8}$`),
	// 'NL'+9 digits+B+2-digit company index – e.g. NL999999999B01
	VAT: regexp.MustCompile(`^NL[0-9]{9}B[0-9]{2}$`),
}, {
	A2: "NO", A3: "NOR", Num: "578",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?47)?[49][0-9]{7}$`),
}, {
	A2: "NP", A3: "NPL", Num: "524",
	Zip:   regexp.MustCompile(`^(?:10|21|22|32|33|34|44|45|56|57)[0-9]{3}$|^(?:977)$`),
	Phone: regexp.MustCompile(`^(?:\+?977)?9[78][0-9]{8}$`),
}, {
	A2: "NR", A3: "NRU", Num: "520",
}, {
	A2: "NU", A3: "NIU", Num: "570",
}, {
	A2: "NZ", A3: "NZL", Num: "554",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?64|0)[28][0-9]{7,9}$`),
	VAT:   regexp.MustCompile(`^[0-9]{9}$`),
}, {
	A2: "OM", A3: "OMN", Num: "512",
	Zip:   rxZip3Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+|00)968)?(?:9[1-9])[0-9]{6}$`),
}, {
	A2: "PA", A3: "PAN", Num: "591",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?507)[0-9]{7,8}$`),
}, {
	A2: "PE", A3: "PER", Num: "604",
	Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:PE [0-9]{4})$`),
	Phone: regexp.MustCompile(`^(?:\+?51)?9[0-9]{8}$`),
	VAT:   regexp.MustCompile(`^[0-9]{11}$`),
}, {
	A2: "PF", A3: "PYF", Num: "258",
	Zip: regexp.MustCompile(`^987(?:[0-8][0-9]|90)$`),
}, {
	A2: "PG", A3: "PNG", Num: "598",
	Zip: regexp.MustCompile(`^[0-9]{3}$`),
}, {
	A2: "PH", A3: "PHL", Num: "608",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:09|\+639)[0-9]{9}$`),
}, {
	A2: "PK", A3: "PAK", Num: "586",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+92)|(?:0092))-?[0-9]{3}-?[0-9]{7}$|^[0-9]{11}$|^[0-9]{4}-[0-9]{7}$`),
}, {
	A2: "PL", A3: "POL", Num: "616",
	Zip:   regexp.MustCompile(`^[0-9]{2}-[0-9]{3}$`),
	Phone: regexp.MustCompile(`^(?:\+?48)?[ ]?[5-8][0-9][ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2}$`),
	VAT:   StringMatcherFunc(plVAT),
}, {
	A2: "PM", A3: "SPM", Num: "666",
	Zip: regexp.MustCompile(`^97500$`),
}, {
	A2: "PN", A3: "PCN", Num: "612",
	Zip: regexp.MustCompile(`^PCRN 1ZZ$`),
}, {
	A2: "PR", A3: "PRI", Num: "630",
	Zip: regexp.MustCompile(`^00[679][0-9]{2}(?:[ -][0-9]{4})?$`),
}, {
	A2: "PS", A3: "PSE", Num: "275",
}, {
	A2: "PT", A3: "PRT", Num: "620",
	Zip:   regexp.MustCompile(`^[0-9]{4}\-[0-9]{3}?$`),
	Phone: regexp.MustCompile(`^(?:\+?351)?9[1236][0-9]{7}$`),
	VAT:   StringMatcherFunc(ptVAT),
}, {
	A2: "PW", A3: "PLW", Num: "585",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "PY", A3: "PRY", Num: "600",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?595|0)9[9876][0-9]{7}$`),
	VAT:   StringMatcherFunc(pyVAT),
}, {
	A2: "QA", A3: "QAT", Num: "634",
}, {
	A2: "RE", A3: "REU", Num: "638",
	Zip:   regexp.MustCompile(`^974(?:[0-8][0-9]|90)$`),
	Phone: regexp.MustCompile(`^(?:\+?262|0|00262)[67][0-9]{8}$`),
}, {
	A2: "RO", A3: "ROU", Num: "642",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?4?0)[ ]?7[0-9]{2}(?:\/| |\.|\-)?[0-9]{3}(?: |\.|\-)?[0-9]{3}$`),
	VAT:   regexp.MustCompile(`^RO[0-9]{8}$`),
}, {
	A2: "RS", A3: "SRB", Num: "688",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+3816|06)[\- 0-9]{5,9}$`),
	VAT:   StringMatcherFunc(rsVAT),
}, {
	A2: "RU", A3: "RUS", Num: "643",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?7|8)?9[0-9]{9}$`),
}, {
	A2: "RW", A3: "RWA", Num: "646",
	Phone: regexp.MustCompile(`^(?:\+?250|0)?[7][0-9]{8}$`),
}, {
	A2: "SA", A3: "SAU", Num: "682",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:(?:\+?966)|0)?5[0-9]{8}$`),
}, {
	A2: "SB", A3: "SLB", Num: "090",
}, {
	A2: "SC", A3: "SYC", Num: "690",
}, {
	A2: "SD", A3: "SDN", Num: "729", Zip: rxZip5Digits,
}, {
	A2: "SE", A3: "SWE", Num: "752",
	Zip:   regexp.MustCompile(`^[1-9][0-9]{2}\s?[0-9]{2}$`),
	Phone: regexp.MustCompile(`^(?:\+?46|0)[ \-]?7[ \-]?[02369](?:[ \-]?[0-9]){7}$`),
	VAT:   StringMatcherFunc(seVAT),
}, {
	A2: "SG", A3: "SGP", Num: "702",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+65)?[689][0-9]{7}$`),
}, {
	A2: "SH", A3: "SHN", Num: "654",
	Zip: regexp.MustCompile(`^(?:STHL|ASCN|TDCU) 1ZZ$`),
}, {
	A2: "SI", A3: "SVN", Num: "705",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+386[ ]?|0)(?:(?:[0-9]{1}[ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2})|(?:[0-9]{2}(?:[ ]?[0-9]{3}){2}))$`),
	VAT:   StringMatcherFunc(siVAT),
}, {
	A2: "SJ", A3: "SJM", Num: "744", Zip: rxZip4Digits,
}, {
	A2: "SK", A3: "SVK", Num: "703",
	Zip:   regexp.MustCompile(`^[0-9]{3}\s?[0-9]{2}$`),
	Phone: regexp.MustCompile(`^(?:\+?421)?[ ]?[1-9][0-9]{2}(?:[ ]?[0-9]{3}){2}$`),
	VAT:   StringMatcherFunc(skVAT),
}, {
	A2: "SL", A3: "SLE", Num: "694",
	Phone: regexp.MustCompile(`^(?:0|94|\+94)?(?:7(?:0|1|2|5|6|7|8)(?: |-)?[0-9])[0-9]{6}$`),
}, {
	A2: "SM", A3: "SMR", Num: "674",
	Zip:   regexp.MustCompile(`^4789[0-9]$`),
	Phone: regexp.MustCompile(`^(?:(?:\+378)|(?:0549)|(?:\+390549)|(?:\+3780549))?6[0-9]{5,9}$`),
	VAT:   regexp.MustCompile(`^[0-9]{5}$`),
}, {
	A2: "SN", A3: "SEN", Num: "686",
	Zip: rxZip5Digits,
}, {
	A2: "SO", A3: "SOM", Num: "706",
	Zip: regexp.MustCompile(`^[A-Z]{2} [0-9]{5}$`),
}, {
	A2: "SR", A3: "SUR", Num: "740",
}, {
	A2: "SS", A3: "SSD", Num: "728",
}, {
	A2: "ST", A3: "STP", Num: "678",
}, {
	A2: "SV", A3: "SLV", Num: "222",
	Zip: rxZip4Digits,
}, {
	A2: "SX", A3: "SXM", Num: "534",
}, {
	A2: "SY", A3: "SYR", Num: "760",
	Phone: regexp.MustCompile(`^(?:(?:\+?963)|0)?9[0-9]{8}$`),
}, {
	A2: "SZ", A3: "SWZ", Num: "748",
	Zip: regexp.MustCompile(`^[HMSL][0-9]{3}$`),
}, {
	A2: "TC", A3: "TCA", Num: "796",
	Zip: regexp.MustCompile(`^TKCA 1ZZ$`),
}, {
	A2: "TD", A3: "TCD", Num: "148",
}, {
	A2: "TF", A3: "ATF", Num: "260",
}, {
	A2: "TG", A3: "TGO", Num: "768",
}, {
	A2: "TH", A3: "THA", Num: "764",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+66|66|0)[0-9]{9}$`),
}, {
	A2: "TJ", A3: "TJK", Num: "762",
	Zip: rxZip6Digits,
}, {
	A2: "TK", A3: "TKL", Num: "772",
}, {
	A2: "TL", A3: "TLS", Num: "626",
}, {
	A2: "TM", A3: "TKM", Num: "795",
	Zip: rxZip6Digits,
}, {
	A2: "TN", A3: "TUN", Num: "788",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?216)?[2459][0-9]{7}$`),
}, {
	A2: "TO", A3: "TON", Num: "776",
}, {
	A2: "TR", A3: "TUR", Num: "792",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?90|0)?5[0-9]{9}$`),
	VAT:   regexp.MustCompile(`^[0-9]{10}$`),
}, {
	A2: "TT", A3: "TTO", Num: "780",
	Zip: rxZip6Digits,
}, {
	A2: "TV", A3: "TUV", Num: "798",
}, {
	A2: "TW", A3: "TWN", Num: "158",
	Zip:   regexp.MustCompile(`^[0-9]{3}(?:[0-9]{2})?$`),
	Phone: regexp.MustCompile(`^(?:\+?886-?|0)?9[0-9]{8}$`),
}, {
	A2: "TZ", A3: "TZA", Num: "834",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?255|0)?[67][0-9]{8}$`),
}, {
	A2: "UA", A3: "UKR", Num: "804",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?38|8)?0[0-9]{9}$`),
	VAT:   regexp.MustCompile(`^[0-9]{12}$`),
}, {
	A2: "UG", A3: "UGA", Num: "800",
	Phone: regexp.MustCompile(`^(?:\+?256|0)?7[0-9]{8}$`),
}, {
	A2: "UM", A3: "UMI", Num: "581",
	Zip: regexp.MustCompile(`^96898$`),
}, {
	A2: "US", A3: "USA", Num: "840",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-?[0-9]{4})?$`),
	Phone: regexp.MustCompile(`^(?:(?:\+?1)?[ -]?)?` +
		`(?:\([2-9][0-9]{2}\)|[2-9][0-9]{2})` +
		`[ -]?(?:[2-9][0-9]{2}[ -]?[0-9]{4})$`),
}, {
	A2: "UY", A3: "URY", Num: "858",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+598|0)9[1-9][0-9]{6}$`),
	VAT:   regexp.MustCompile(`^[0-9]{12}$`),
}, {
	A2: "UZ", A3: "UZB", Num: "860",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?998)?(?:6[125-79]|7[1-69]|88|9[0-9])[0-9]{7}$`),
	VAT:   StringMatcherFunc(uzVAT),
}, {
	A2: "VA", A3: "VAT", Num: "336",
	Zip: regexp.MustCompile(`^00120$`),
}, {
	A2: "VC", A3: "VCT", Num: "670",
	Zip: regexp.MustCompile(`^VC[0-9]{4}$`),
}, {
	A2: "VE", A3: "VEN", Num: "862",
	Zip: regexp.MustCompile(`^[0-9]{4}(?:-[A-Z])?$`),
	// First digit must be (J, G, V, E), one dash (-), next 9 (nine) numbers
	// like J-305959918, in some cases can be written like J-30595991-8
	VAT: regexp.MustCompile(`^[JGVE]-[0-9]{8}-?[0-9]$`),
}, {
	A2: "VG", A3: "VGB", Num: "092",
	Zip: regexp.MustCompile(`^VG11(?:[1-5][0-9]|60)$`),
}, {
	A2: "VI", A3: "VIR", Num: "850",
	Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
}, {
	A2: "VN", A3: "VNM", Num: "704",
	Zip:   rxZip6Digits,
	Phone: regexp.MustCompile(`^(?:\+?84|0)(?:3[2-9]|5[2689]|7[0|6-9]|8[1-6|89]|9[0-9])(?:[0-9]{7})$`),
}, {
	A2: "VU", A3: "VUT", Num: "548",
}, {
	A2: "WF", A3: "WLF", Num: "876",
	Zip: regexp.MustCompile(`^986(?:[0-8][0-9]|90)$`),
}, {
	A2: "WS", A3: "WSM", Num: "882",
	Zip: regexp.MustCompile(`^WS[0-9]{4}$`),
}, {
	A2: "YE", A3: "YEM", Num: "887",
}, {
	A2: "YT", A3: "MYT", Num: "175",
	Zip: regexp.MustCompile(`^976(?:[0-8][0-9]|90)$`),
}, {
	A2: "ZA", A3: "ZAF", Num: "710",
	Zip:   rxZip4Digits,
	Phone: regexp.MustCompile(`^(?:\+?27|0)[0-9]{9}$`),
}, {
	A2: "ZM", A3: "ZMB", Num: "894",
	Zip:   rxZip5Digits,
	Phone: regexp.MustCompile(`^(?:\+?26)?09[567][0-9]{7}$`),
}, {
	A2: "ZW", A3: "ZWE", Num: "716",
	Phone: regexp.MustCompile(`^(?:\+263)[0-9]{9}$`),
}}
