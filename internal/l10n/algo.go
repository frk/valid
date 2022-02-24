package l10n

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/frk/valid/internal/algo"
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

// 11 digit number formed from a 9 digit unique identifier and two prefix
// check digits. The two leading digits (the check digits) will be derived
// from the subsequent 9 digits using a modulus 89 check digit calculation.
// - https://en.wikipedia.org/wiki/VAT_identification_number
// - https://abr.business.gov.au/Help/AbnFormat
var vatAU = regexp.MustCompile(`^[1-9][0-9]{10}$`)
var weightsAU = []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

func auVAT(v string) bool {
	if !vatAU.MatchString(v) {
		return false
	}

	b := []byte(v)
	n, _ := strconv.Atoi(string(b[0]))
	b[0] = strconv.Itoa(n - 1)[0]

	v = string(b)

	sum := 0
	for i := 0; i < len(v); i++ {
		n, _ := strconv.Atoi(string(v[i]))
		sum += n * weightsAU[i]
	}
	return (sum % 89) == 0
}

// 'BE'+ 8 digits + 2 check digits – e.g. BE09999999XX
var vatBE = regexp.MustCompile(`^BE[01][0-9]{9}$`)

func beVAT(v string) bool {
	if !vatBE.MatchString(v) {
		return false
	}
	v = v[2:]

	// mod 97
	num, _ := strconv.Atoi(v[:8])
	chk, _ := strconv.Atoi(v[8:])
	return num%97 == 97-chk
}

// 6 digits (up to 31 December 2013). CHE 9 numeric digits plus TVA/MWST/IVA
// e.g. CHE-123.456.788 TVA[20] The last digit is a MOD11 checksum digit build
// with weighting pattern: 5,4,3,2,7,6,5,4
// - https://en.wikipedia.org/wiki/VAT_identification_number
var vatCH = regexp.MustCompile(``) // TODO

func chVAT(v string) bool {
	if !vatCH.MatchString(v) {
		return false
	}
	// TODO
	return false
}

// 8 digits – e.g. DK99999999, last digit is check digit
var vatDK = regexp.MustCompile(`^DK[0-9]{8}$`)
var weigthsDK = []int{2, 7, 6, 5, 4, 3, 2}

func dkVAT(v string) bool {
	if !vatDK.MatchString(v) {
		return false
	}
	v = v[2:]

	// The check digit is calculated utilizing MOD 11-2.
	// https://web.archive.org/web/20120917151518/http://www.erhvervsstyrelsen.dk/modulus_11
	sum := 0
	for i := 0; i < len(v)-1; i++ {
		num, _ := strconv.Atoi(string(v[i]))
		sum += num * weigthsDK[i]
	}

	mod := (sum % 11)
	if mod == 1 {
		return false
	}

	check, _ := strconv.Atoi(string(v[len(v)-1]))
	if mod == 0 {
		return mod == check
	}
	return check == (11 - mod)
}

// FI + 7 digits + check digit, e.g. FI99999999
var vatFI = regexp.MustCompile(`^FI[0-9]{8}$`)
var weigthsFI = []int{7, 9, 10, 5, 8, 4, 2}

func fiVAT(v string) bool {
	if !vatFI.MatchString(v) {
		return false
	}
	v = v[2:]

	// The check digit is calculated utilizing MOD 11-2.
	// http://tarkistusmerkit.teppovuori.fi/tarkmerk.htm#y-tunnus2
	sum := 0
	for i := 0; i < len(v)-1; i++ {
		num, _ := strconv.Atoi(string(v[i]))
		sum += num * weigthsFI[i]
	}

	mod := (sum % 11)
	if mod == 1 {
		return false
	}

	check, _ := strconv.Atoi(string(v[len(v)-1]))
	if mod == 0 {
		return mod == check
	}
	return check == (11 - mod)
}

// - 'FR'+ 2 digits (as validation key ) + 9 digits (as SIREN), the first and/or the
//   second value can also be a character (any except O or I) - e.g. FRXX999999999
// References:
// - https://en.wikipedia.org/wiki/VAT_identification_number
// - https://www.gov.uk/guidance/vat-eu-country-codes-vat-numbers-and-vat-in-other-languages
var vatFR = regexp.MustCompile(`^FR[A-HJ-NP-Z0-9]{2}[0-9]{9}$`)

func frVAT(v string) bool {
	if !vatFR.MatchString(v) {
		return false
	}
	v = v[2:]

	// The SIREN number ought to be valid Luhn: https://en.wikipedia.org/wiki/SIREN_code
	if !algo.Luhn(v[2:]) {
		return false
	}

	key, err := strconv.ParseInt(v[:2], 10, 32)
	if err != nil {
		return false
	}

	siren, err := strconv.ParseInt(v[2:], 10, 32)
	if err != nil {
		return false
	}

	// The validation key is calculated as follows:
	// [ 12 + 3 * ( SIREN modulo 97 ) ] modulo 97
	return key == (12+(3*(siren%97)))%97

	// NOTE(mkopriva): couldn't find anything official and/or
	// open source that demonstrates how to handle a validation
	// key that contains letters ... this will fail if such a
	// key is provided even if it is valid.
}

// https://tatief.wordpress.com/2008/12/29/αλγόριθμος-του-αφμ-έλεγχος-ορθότητας/
var vatGR = regexp.MustCompile(`^(?:EL|GR)[0-9]{9}$`)
var weightsGR = []int{256, 128, 64, 32, 16, 8, 4, 2}

func grVAT(v string) bool {
	if !vatGR.MatchString(v) {
		return false
	}
	v = v[2:]

	sum := 0
	for i := 0; i < 8; i++ {
		n, _ := strconv.Atoi(string(v[i]))
		sum += n * weightsGR[i]
	}

	chk, _ := strconv.Atoi(string(v[8]))
	return chk == (sum%11)%10
}

// 'HR'+ 11 digit number, e.g. HR12345678901
var vatHR = regexp.MustCompile(`^HR[0-9]{11}$`)

func hrVAT(v string) bool {
	if !vatHR.MatchString(v) {
		return false
	}
	return ISO7064_MOD11_10(v[2:])
}

// 9 digit number. If the number of digits is less than 9, then zeros
// should be padded to the left side. The leftmost digit is 5 for corporations.
// Other leftmost digits are used for individuals. The rightmost digit
// is a check digit (using Luhn algorithm).
var vatIL = regexp.MustCompile(`^[0-9]{9}$`)

func ilVAT(v string) bool {
	if !vatIL.MatchString(v) {
		return false
	}
	return algo.Luhn(v)
}

var zipIN = regexp.MustCompile(`^[1-9][0-9]{2}[ ]?[0-9]{3}$`)
var zipNegIN = regexp.MustCompile(`^(?:10|29|35|54|55|65|66|86|87|88|89)`)

func inZIP(v string) bool {
	return zipIN.MatchString(v) && !zipNegIN.MatchString(v)
}

// 11 digits (the first 7 digits is a progressive number, the following 3
// means the province of residence, the last digit is a check number
// - The check digit is calculated using Luhn's Algorithm.)
var vatIT = regexp.MustCompile(`^IT[0-9]{11}$`)

func itVAT(v string) bool {
	if !vatIT.MatchString(v) {
		return false
	}
	return algo.Luhn(v[2:])
}

// 10 digits, the last one is a check digit; for convenience the
// digits are separated by hyphens (xxx-xxx-xx-xx or xxx-xx-xx-xxx
// for legal people), but formally the number consists only of digits
var vatPL = regexp.MustCompile(`^PL(?:[0-9]{10}|(?:[0-9]{3}-){2}[0-9]{2}-[0-9]{2}|[0-9]{3}-(?:[0-9]{2}-){2}[0-9]{3})$`)
var weightsPL = []int{6, 5, 7, 2, 3, 4, 5, 6, 7}

func plVAT(v string) bool {
	if !vatPL.MatchString(v) {
		return false
	}

	// remove potential hyphens
	v = strings.Map(func(r rune) rune {
		if r == '-' {
			return -1
		}
		return r
	}, v[2:])

	// https://pl.wikipedia.org/wiki/NIP
	sum := 0
	for i := 0; i < len(v)-1; i++ {
		num, _ := strconv.Atoi(string(v[i]))
		sum += num * weightsPL[i]
	}

	check, _ := strconv.Atoi(string(v[len(v)-1]))
	return check == (sum % 11)
}

// 9 digits; the last digit is the check digit. The first digit depends on
// what the number refers to, e.g.: 1-3 are regular people, 5 are companies.
var vat1PT = regexp.MustCompile(`^PT[0-9]{9}$`)
var vat2PT = regexp.MustCompile(`^PT(?:(?:1|2|3|5|6|8)|(?:45|70|71|72|74|75|77|79|90|91|98|99))`)

func ptVAT(v string) bool {
	if !vat1PT.MatchString(v) || !vat2PT.MatchString(v) {
		return false
	}
	v = v[2:]

	// MOD 11-2
	// https://pt.wikipedia.org/wiki/Número_de_identificação_fiscal
	sum := 0
	for i := 0; i < len(v)-1; i++ {
		num, _ := strconv.Atoi(string(v[i]))
		sum += num * (9 - i)
	}

	mod := (sum % 11)
	check, _ := strconv.Atoi(string(v[len(v)-1]))
	if mod < 2 {
		return check == 0
	}
	return check == (11 - mod)
}

// 6 to 8 digits, 1 dash, 1 check sum digit
var vatPY = regexp.MustCompile(`^[0-9]{6,8}-[0-9]$`)

func pyVAT(v string) bool {
	if !vatPY.MatchString(v) {
		return false
	}
	// TODO
	return false
}

// 9 digits (ex. 123456788) of which the first 8 are the actual ID number,
// and the last digit is a checksum digit, calculated according to ISO 7064,
// MOD 11-10
var vatRS = regexp.MustCompile(`^[0-9]{9}$`)

func rsVAT(v string) bool {
	if !vatRS.MatchString(v) {
		return false
	}
	return ISO7064_MOD11_10(v)
}

// VAT:
//	12 digits, of which the last two are most often 01 e.g. SE999999999901.
// 	(For sole proprietors who have several businesses the numbers can be 02,
// 	03 and so on, since sole proprietors only have their personnummer as the
// 	organisationsnummer. The first 10 digits are the same as the Swedish
// 	organisationsnummer.
//
// 	The last digit of the first 10 digit number is the control digit, it is
// 	calculated according to the Luhn algorithm.
//
// 	Reference:
// 	- https://en.wikipedia.org/wiki/VAT_identification_number
// 	- https://sv.wikipedia.org/wiki/Organisationsnummer
var vatSE = regexp.MustCompile(`^SE[0-9]{12}$`)

func seVAT(v string) bool {
	if !vatSE.MatchString(v) {
		return false
	}
	return algo.Luhn(v[2 : len(v)-2])
}

// VAT: (https://web.archive.org/web/20150215110937/http://www.durs.gov.si/si/storitve/vpis_v_davcni_register_in_davcna_stevilka/vpis_v_davcni_register_in_davcna_stevilka_pojasnila/davcna_stevilka_splosno)
//
// Davčna številka je sestavljena iz osmih številk; sedem je naključno izbranih, osma pa je izračunana po modulu 11:
//
// prvih sedem mest je osnovna številka, ki je naključno izbrana iz nabora številk od 1.000.000 do 9.999.999,
// na osmem mestu je kontrolna številka, izračunana po modulu 11.
//
// Postopek izračuna kontrolne številke:
//
// posamezno številko osnovne številke pomnožimo s konstantnimi ponderji 8, 7, 6, 5, 4, 3 in 2,
// zmnožke seštejemo,
// seštevek delimo z 11,
// ostanek deljenja odštejemo od 11 in razlika je kontrolna številka.
//
// Če je ostanek deljenja 1 in je razlika zato 10, je kontrolna številka 0. Če je ostanek deljenja 0 in je razlika zato 11, se ta osnovna številka izključi iz nabora možnih davčnih številk.
//
// Primer:
//
// osnovna številka: 1 5 0 1 2 5 5
// konstantni ponderji: 8,7,6,5,4,3,2
// seštevek zmnožkov: 8 + 35 + 0 + 5 + 8 + 15 + 10 = 81
// seštevek delimo z 11: 7
// ostanek je: 4
// razlika do 11: 7 (kontrolna št.)
// davčna številka: 15012557
var vatSI = regexp.MustCompile(`^SI[1-9][0-9]{7}$`)
var weightsSI = []int{8, 7, 6, 5, 4, 3, 2}

func siVAT(v string) bool {
	if !vatSI.MatchString(v) {
		return false
	}
	v = v[2:]

	sum := 0
	for i := 0; i < len(v)-1; i++ {
		n, _ := strconv.Atoi(string(v[i]))
		sum += n * weightsSI[i]
	}

	res := (11 - sum%11)
	if res == 10 {
		res = 0
	}

	chk, _ := strconv.Atoi(string(v[len(v)-1]))
	return chk == res
}

// 'SK'+10 digits (number must be divisible by 11)
var vatSK = regexp.MustCompile(`^SK[0-9]{10}$`)

func skVAT(v string) bool {
	if !vatSK.MatchString(v) {
		return false
	}

	n, _ := strconv.Atoi(v[2:])
	return (n % 11) == 0
}

// 9 digits
// Companies: 20000000X-29999999X
// People: 40000000X-79999999X
//
// The taxpayer identification number consists of nine digits, the first
// 8 digits are the taxpayer’s own number, and the last digit is the control
// number.  The control number is formed by a certain algorithm in the process
// of assigning a taxpayer identification number by a computer in the State
// Tax Committee of the Republic of Uzbekistan.
// - https://www.lex.uz/acts/499945
var vatUZ = regexp.MustCompile(`^[24-7][0-9]{8}$`)

func uzVAT(v string) bool {
	if !vatUZ.MatchString(v) {
		return false
	}

	// TODO if possible figure out what algorithm is used for the control digit

	return true
}
