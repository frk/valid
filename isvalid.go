// NOTE(mkopriva): DO NOT REMOVE THE "isvalid:rule" DIRECTIVES AND THE JSON THAT
// FOLLOWS AFTER THEM. The directives & json in the functions' comments are required
// by the internal/analysis code to resolve the properties of a rule's function.

package isvalid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/frk/isvalid/internal/cldr"
	"github.com/frk/isvalid/internal/tables"
	"github.com/frk/isvalid/l10n/country"
)

var _ = log.Println
var _ = fmt.Println

// convenience func for removing characters from a string
func rmchar(v string, f func(r rune) bool) string {
	return strings.Map(func(r rune) rune {
		if f(r) {
			return -1
		}
		return r
	}, v)

}

// convenience func for strings known to contain digits only
func atoi(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

var rxASCII = regexp.MustCompile(`^[[:ascii:]]*$`)

// ASCII reports whether or not v is an ASCII string.
//
//	isvalid:rule
//	{
//		"name": "ascii",
//		"err": { "text": "must contain only ASCII characters" }
//	}
func ASCII(v string) bool {
	return rxASCII.MatchString(v)
}

// Alpha reports whether or not v is a valid alphabetic string.
//
//	isvalid:rule
//	{
//		"name": "alpha",
//		"opts": [[ { "key": null, "value": "en" } ]],
//		"err": { "text": "must be an alphabetic string" }
//	}
func Alpha(v string, lang string) bool {
	lang = strings.ToLower(lang)
	if len(lang) == 0 {
		lang = "en"
	}

	if rx, ok := tables.Alpha[lang]; ok {
		return rx.MatchString(v)
	}
	return false
}

// Alnum reports whether or not v is a valid alphanumeric string.
//
//	isvalid:rule
//	{
//		"name": "alnum",
//		"opts": [[ { "key": null, "value": "en" } ]],
//		"err": { "text": "must be an alphanumeric string" }
//	}
func Alnum(v string, lang string) bool {
	lang = strings.ToLower(lang)
	if len(lang) == 0 {
		lang = "en"
	}

	if rx, ok := tables.Alnum[lang]; ok {
		return rx.MatchString(v)
	}
	return false
}

var rxBIC = regexp.MustCompile(`^[A-z]{4}[A-z]{2}\w{2}(\w{3})?$`)

// BIC reports whether or not v represents a valid Bank Identification Code or SWIFT code.
//
//	isvalid:rule
//	{
//		"name": "bic",
//		"err": { "text": "must be a valid BIC or SWIFT code" }
//	}
func BIC(v string) bool {
	return rxBIC.MatchString(v)
}

var rxBTC = regexp.MustCompile(`^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`)

// BTC reports whether or not v represents a valid BTC address.
//
//	isvalid:rule
//	{
//		"name": "btc",
//		"err": { "text": "must be a valid BTC address" }
//	}
func BTC(v string) bool {
	return rxBTC.MatchString(v)
}

var rxBase32 = regexp.MustCompile(`^[A-Z2-7]+=*$`)

// Base32 reports whether or not v is a valid base32 string.
//
//	isvalid:rule
//	{
//		"name": "base32",
//		"err": { "text": "must be a valid base32 string" }
//	}
func Base32(v string) bool {
	return (len(v)%8 == 0) && rxBase32.MatchString(v)
}

var rxBase58 = regexp.MustCompile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$`)

// Base58 reports whether or not v is a valid base58 string.
//
//	isvalid:rule
//	{
//		"name": "base58",
//		"err": { "text": "must be a valid base58 string" }
//	}
func Base58(v string) bool {
	return rxBase58.MatchString(v)
}

// Base64 reports whether or not v is a valid base64 string. NOTE: The standard
// "encoding/base64" package is used for validation. With urlsafe=false StdEncoding
// is used and with urlsafe=true RawURLEncoding is used.
//
//	isvalid:rule
//	{
//		"name": "base64",
//		"opts": [[
//			{ "key": null, "value": "false"	},
//			{ "key": "url", "value": "true" }
//		]],
//		"err": { "text": "must be a valid base64 string" }
// 	}
func Base64(v string, urlsafe bool) bool {
	if urlsafe {
		if i := strings.IndexAny(v, "\r\n"); i > -1 {
			return false
		}
		_, err := base64.RawURLEncoding.DecodeString(v)
		return err == nil
	}
	_, err := base64.StdEncoding.DecodeString(v)
	return err == nil
}

var rxBinary = regexp.MustCompile(`^(?:0[bB])?[0-1]+$`)

// Binary reports whether or not v represents a valid binary integer.
//
//	isvalid:rule
//	{
//		"name": "binary",
//		"err": { "text": "string content must match a binary number" }
//	}
func Binary(v string) bool {
	return rxBinary.MatchString(v)
}

// Bool reports whether or not v represents a valid boolean. The following
// are considered valid boolean values: "true", "false", "TRUE", and "FALSE".
//
//	isvalid:rule
//	{
//		"name": "bool",
//		"err": { "text": "string content must match a boolean value" }
//	}
func Bool(v string) bool {
	if len(v) == 4 && (v == "true" || v == "TRUE") {
		return true
	}
	if len(v) == 5 && (v == "false" || v == "FALSE") {
		return true
	}
	return false
}

// CIDR reports whether or not v is a valid Classless Inter-Domain Routing notation.
// NOTE: CIDR uses "net".ParseCIDR to determine the validity of v.
//
//	isvalid:rule
//	{
//		"name": "cidr",
//		"err": { "text": "must be a valid CIDR notation" }
//	}
func CIDR(v string) bool {
	_, _, err := net.ParseCIDR(v)
	return err == nil
}

var rxCVV = regexp.MustCompile(`^[0-9]{3,4}$`)

// CVV reports whether or not v is a valid Card Verification Value.
//
//	isvalid:rule
//	{
//		"name": "cvv",
//		"err": { "text": "must be a valid CVV" }
//	}
func CVV(v string) bool {
	return rxCVV.MatchString(v)
}

type CurrencyOpts struct {
	NeedSym bool
}

var CurrencyOptsDefault = CurrencyOpts{
	NeedSym: false,
}

// Currency reports whether or not v represents a valid Currency amount.
//
//	isvalid:rule
//	{
//		"name": "ccy",
//		"opts": [[
//			{ "key": null, "value": "usd" }
//		], [
//			{ "key": null, "value": "nil" }
//		]],
//		"err": { "text": "must be a valid currency amount" }
// 	}
func Currency(v string, code string, opts *CurrencyOpts) bool {
	if len(v) == 0 || len(code) != 3 {
		return false
	}
	code = strings.ToUpper(code)
	ccy, ok := tables.ISO4217[code]
	if !ok {
		return false
	}

	if opts == nil {
		opts = &CurrencyOptsDefault
	}

	// check for leading currency symbol
	if ccy.Left && len(ccy.Sym) > 0 {
		hasSym := false
		for _, sym := range ccy.Sym {
			if len(v) >= len(sym) && v[:len(sym)] == sym {
				v = v[len(sym):]
				hasSym = true
				break
			}
		}
		if !hasSym && opts.NeedSym {
			return false
		}
		if hasSym {
			for _, sep := range ccy.Sep {
				if len(v) >= len(sep) && v[:len(sep)] == sep {
					v = v[len(sep):]
					break
				}
			}
		}
	}

	// check for traling currency symbol
	if ccy.Right && len(ccy.Sym) > 0 {
		hasSym := false
		for _, sym := range ccy.Sym {
			if len(v) >= len(sym) && v[len(v)-len(sym):] == sym {
				v = v[:len(v)-len(sym)]
				hasSym = true
				break
			}
		}
		if !hasSym && opts.NeedSym {
			return false
		}
		if hasSym {
			for _, sep := range ccy.Sep {
				if len(v) >= len(sep) && v[len(v)-len(sep):] == sep {
					v = v[:len(v)-len(sep)]
					break
				}
			}
		}
	}

	// split into number groups
	ngs, last := []string{}, 0
	var thoSep rune

numLoop:
	for i, r := range v {
		if r >= '0' && r <= '9' {
			continue numLoop
		}
		for _, sep := range ccy.ThoSep {
			if r == sep {
				// mixed thousands separators; fail
				if thoSep > 0 && sep != thoSep {
					return false
				}

				ngs = append(ngs, v[last:i])
				last = i + utf8.RuneLen(r)
				thoSep = sep
				continue numLoop
			}
		}
		for _, sep := range ccy.DecSep {
			if r == sep {
				ngs = append(ngs, v[last:i])
				last = i + utf8.RuneLen(r)

				decimal := v[last:]
				// incorrect number of decimal digits; fail
				if len(decimal) != ccy.DecNum {
					return false
				}
				// make sure the rest is all digit
				for _, d := range decimal {
					if d < '0' || d > '9' {
						return false
					}
				}
				last = len(v) // no need to add decimal to ngs
				break numLoop
			}
		}

		// unacceptable rune; fail
		return false
	}
	if last < len(v)-1 {
		ngs = append(ngs, v[last:])
	}

	for i, ng := range ngs {
		if len(ng) == 0 {
			return false
		}
		if i > 0 && len(ng) != 3 {
			return false
		}
		if i == 0 {
			if len(ng) > 3 && thoSep != 0 {
				return false
			}
			if ng[0] == '0' && (len(ng) > 1 || thoSep != 0) {
				return false
			}
		}
	}
	return true
}

var rxDataURIMediaType = regexp.MustCompile(`^(?i)[a-z]+\/[a-z0-9\-\+]+$`)
var rxDataURIAttribute = regexp.MustCompile(`^(?i)^[a-z\-]+=[a-z0-9\-]+$`)
var rxDataURIData = regexp.MustCompile(`^(?i)[a-z0-9!\$&'\(\)\*\+,;=\-\._~:@\/\?%\s]*$`)

// DataURI reports whether or not v is a valid data URI.
//
//	isvalid:rule
//	{
//		"name": "datauri",
//		"err": { "text": "must be a valid data URI" }
//	}
func DataURI(v string) bool {
	vv := strings.Split(v, ",")
	if len(vv) < 2 {
		return false
	}

	attrs := strings.Split(strings.TrimSpace(vv[0]), ";")
	if attrs[0][:5] != "data:" {
		return false
	}

	mediaType := attrs[0][5:]
	if len(mediaType) > 0 && !rxDataURIMediaType.MatchString(mediaType) {
		return false
	}

	for i := 1; i < len(attrs); i++ {
		if i == len(attrs)-1 && strings.ToLower(attrs[i]) == "base64" {
			continue // ok
		}

		if !rxDataURIAttribute.MatchString(attrs[i]) {
			return false
		}
	}

	data := vv[1:]
	for i := 0; i < len(data); i++ {
		if !rxDataURIData.MatchString(data[i]) {
			return false
		}
	}
	return true
}

// Decimal reports whether or not v represents a valid decimal number.
//
//	isvalid:rule
//	{
//		"name": "decimal",
//		"opts": [[
//			{ "key": null, "value": "en" }
//		]]
//		"err": { "text": "string content must match a decimal number" }
//	}
func Decimal(v string, locale string) bool {
	if len(v) > 1 && (v[0] == '+' || v[0] == '-') {
		v = v[1:]
	}
	if len(v) == 0 {
		return false
	}

	loc, ok := cldr.Locale(locale)
	if !ok {
		return false
	}

	var last int
	var groups []string
	for i, r := range v {
		if r >= loc.DigitZero && r <= loc.DigitNine {
			continue
		}
		if loc.SepGroup > 0 && r == loc.SepGroup {
			groups = append(groups, v[last:i])
			last = i + utf8.RuneLen(r)
			continue
		}

		if loc.SepDecimal > 0 && r == loc.SepDecimal {
			groups = append(groups, v[last:i])
			last = i + utf8.RuneLen(r)

			// trailing decimal; fail
			if last == len(v) {
				return false
			}

			fraction := v[last:]

			// make sure the rest is all digit
			for _, d := range fraction {
				if d < loc.DigitZero || d > loc.DigitNine {
					return false
				}
			}

			last = len(v)
			break
		}

		// unacceptable rune; fail
		return false
	}
	if last < len(v) {
		groups = append(groups, v[last:])
	}

	if len(groups) == 1 && len(groups[0]) == 0 {
		return true
	}

	for i, g := range groups {
		if len(g) == 0 {
			return false
		}
		if i > 0 && len(g) != 3 {
			return false
		}
		if i == 0 {
			if len(g) > 3 && len(groups) > 1 {
				return false
			}
			if g[0] == '0' && len(groups) > 1 {
				return false
			}
		}
	}
	return true
}

var rxDigits = regexp.MustCompile(`^[0-9]+$`)

// Digits reports whether or not v is a string of digits.
//
//	isvalid:rule
//	{
//		"name": "digits",
//		"err": { "text": "must contain only digits" }
//	}
func Digits(v string) bool {
	return rxDigits.MatchString(v)
}

// EAN reports whether or not v is a valid European Article Number.
//
//	isvalid:rule
//	{
//		"name": "ean",
//		"err": { "text": "must be a valid EAN" }
//	}
func EAN(v string) bool {
	length := len(v)
	if length != 8 && length != 13 {
		return false
	}
	if !rxDigits.MatchString(v) {
		return false
	}

	// the accumulate checksum
	sum := 0
	for i, digit := range v[:length-1] {

		// the digit's weigth by position
		weight := 1
		if length == 8 && i%2 == 0 {
			weight = 3
		} else if length == 13 && i%2 != 0 {
			weight = 3
		}

		sum += atoi(string(digit)) * weight
	}

	// the calculated check digit
	check := 0
	if remainder := (10 - (sum % 10)); remainder < 10 {
		check = remainder
	}
	return check == atoi(string(v[length-1]))
}

// EIN reports whether or not v is a valid Employer Identification Number.
//
//	isvalid:rule
//	{
//		"name": "ein",
//		"err": { "text": "must be a valid EIN" }
//	}
func EIN(v string) bool {
	// TODO
	return false
}

var rxETH = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// ETH reports whether or not v is a valid ethereum address.
//
//	isvalid:rule
//	{
//		"name": "eth",
//		"err": { "text": "must be a valid ethereum address" }
//	}
func ETH(v string) bool {
	return rxETH.MatchString(v)
}

// Email reports whether or not v is a valid email address. NOTE: Email uses
// "net/mail".ParseAddress to determine the validity of v.
//
//	isvalid:rule
//	{
//		"name": "email",
//		"err": { "text": "must be a valid email address" }
//	}
func Email(v string) bool {
	_, err := mail.ParseAddress(v)
	return err == nil
}

var rxTLD = regexp.MustCompile(`^(?i:[a-z\x{00a1}-\x{ffff}]{2,}|xn[a-z0-9-]{2,})$`)
var rxTLDIllegal = regexp.MustCompile(`[\s\x{2002}-\x{200B}\x{202F}\x{205F}\x{3000}\x{FEFF}\x{DB40}\x{DC20}\x{00A9}\x{FFFD}]`)
var rxFQDNPart = regexp.MustCompile(`^[a-zA-Z\x{00a1}-\x{ffff}0-9-]+$`)
var rxFQDNPartIllegal = regexp.MustCompile(`[\x{ff01}-\x{ff5e}]`)

// FQDN reports whether or not v is a valid Fully Qualified Domain Name. NOTE: FQDN
// TLD is required, numeric TLDs or trailing dots are disallowed, and underscores
// are forbidden.
//
//	isvalid:rule
//	{
//		"name": "fqdn",
//		"err": { "text": "must be a valid FQDN" }
//	}
func FQDN(v string) bool {
	parts := strings.Split(v, ".")
	for _, part := range parts {
		if len(part) > 63 {
			return false
		}
	}

	// tld must be present, must match pattern, must not contain illegal chars, must not be all digits
	if len(parts) < 2 {
		return false
	}
	tld := parts[len(parts)-1]
	if !rxTLD.MatchString(tld) || rxTLDIllegal.MatchString(tld) || rxDigits.MatchString(tld) {
		return false
	}

	for _, part := range parts {
		if len(part) > 0 && (part[0] == '-' || part[len(part)-1] == '-') {
			return false
		}
		if !rxFQDNPart.MatchString(part) || rxFQDNPartIllegal.MatchString(part) {
			return false
		}
	}
	return true
}

var rxFloat = regexp.MustCompile(`^[+-]?(?:[0-9]*)?(?:\.[0-9]*)?(?:[eE][+-]?[0-9]+)?$`)

// Float reports whether or not v represents a valid float.
//
//	isvalid:rule
//	{
//		"name": "float",
//		"err": { "text": "string content must match a floating point number" }
//	}
func Float(v string) bool {
	if v == "" || v == "." || v == "+" || v == "-" {
		return false
	}
	return rxFloat.MatchString(v)
}

var rxHSLComma = regexp.MustCompile(`^(?i)(?:hsl)a?\(\s*(?:(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?))(?:deg|grad|rad|turn|\s*)(?:\s*,\s*(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?)%){2}\s*(?:,\s*(?:(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?)%?)\s*)?\)$`)
var rxHSLSpace = regexp.MustCompile(`^(?i)(?:hsl)a?\(\s*(?:(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?))(?:deg|grad|rad|turn|\s)(?:\s*(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?)%){2}\s*(?:\/\s*(?:(?:\+|\-)?(?:[0-9]+(?:\.[0-9]+)?(?:e(?:\+|\-)?[0-9]+)?|\.[0-9]+(?:e(?:\+|\-)?[0-9]+)?)%?)\s*)?\)$`)

// HSL reports whether or not v represents an HSL color value.
//
//	isvalid:rule
//	{
//		"name": "hsl",
//		"err": { "text": "must be a valid HSL color" }
//	}
func HSL(v string) bool {
	return rxHSLComma.MatchString(v) || rxHSLSpace.MatchString(v)
}

var rxHash = regexp.MustCompile(`^[0-9A-Fa-f]+$`)

// Hash reports whether or not v is a hash of the specified algorithm.
//
//	isvalid:rule
//	{
//		"name": "hash",
//		"err": { "text": "must be a valid hash" }
//	}
func Hash(v string, algo string) bool {
	if hlen := tables.HashAlgoLen[algo]; hlen > 0 && hlen == len(v) {
		return rxHash.MatchString(v)
	}
	return false
}

var rxHex = regexp.MustCompile(`^(?:0[xXhH])?[0-9A-Fa-f]+$`)

// Hex reports whether or not v is a valid hexadecimal string.
//
//	isvalid:rule
//	{
//		"name": "hex",
//		"err": { "text": "must be a valid hexadecimal string" }
//	}
func Hex(v string) bool {
	return rxHex.MatchString(v)
}

var rxHexColor = regexp.MustCompile(`^#?(?i:[0-9A-F]{3}|[0-9A-F]{4}|[0-9A-F]{6}|[0-9A-F]{8})$`)

// HexColor reports whether or not v is a valid hexadecimal color code.
//
//	isvalid:rule
//	{
//		"name": "hexcolor",
//		"err": { "text": "must represent a valid hexadecimal color code" }
//	}
func HexColor(v string) bool {
	return rxHexColor.MatchString(v)
}

// IBAN reports whether or not v is an International Bank Account Number.
//
//	isvalid:rule
//	{
//		"name": "iban",
//		"err": { "text": "must be a valid IBAN" }
//	}
func IBAN(v string) bool {
	v = rmchar(v, func(r rune) bool { return r == ' ' || r == '-' })
	if len(v) < 2 {
		return false
	}

	v = strings.ToUpper(v)
	if rx, ok := tables.IBANRegexp[v[:2]]; !ok || !rx.MatchString(v) {
		return false
	}

	// rearrange by moving the four initial characters to the end of the string
	v = v[4:] + v[:4]

	// convert to decimal int by replacing each letter in the string with two digits
	var D string
	for _, r := range v {
		if r >= 'A' && 'Z' >= r {
			D += strconv.Itoa(int(r - 55))
		} else {
			D += string(r)
		}
	}

	// The modulo algorithm for checking the IBAN is taken from:
	// https://en.wikipedia.org/wiki/International_Bank_Account_Number#Modulo_operation_on_IBAN
	//
	// 1. Starting from the leftmost digit of D, construct a number using
	//    the first 9 digits and call it N.
	// 2. Calculate N mod 97.
	// 3. Construct a new 9-digit N by concatenating above result (step 2)
	//    with the next 7 digits of D. If there are fewer than 7 digits remaining
	//    in D but at least one, then construct a new N, which will have less
	//    than 9 digits, from the above result (step 2) followed by the remaining
	//    digits of D.
	// 4. Repeat steps 2–3 until all the digits of D have been processed.
	//
	d, D := D[:9], D[9:]
	for len(D) > 0 {
		mod := strconv.Itoa(atoi(d) % 97)

		if len(D) >= 7 {
			d, D = (mod + D[:7]), D[7:]
		} else {
			d, D = (mod + D), ""
		}
	}

	return (atoi(d) % 97) == 1
}

// IC reports whether or not v is an Identity Card number.
//
//	isvalid:rule
//	{
//		"name": "ic",
//		"err": { "text": "must be a valid identity card number" }
//	}
func IC(v string) bool {
	// TODO
	return false
}

var rxIMEI = regexp.MustCompile(`^[0-9]{15}$`)
var rxIMEIHyphenated = regexp.MustCompile(`^\d{2}-\d{6}-\d{6}-\d{1}$`)

// IMEI reports whether or not v is an IMEI number.
//
//	isvalid:rule
//	{
//		"name": "imei",
//		"err": { "text": "must be a valid IMEI number" }
//	}
func IMEI(v string) bool {
	if !rxIMEI.MatchString(v) && !rxIMEIHyphenated.MatchString(v) {
		return false
	}

	v = rmchar(v, func(r rune) bool { return r == '-' })

	vlen := len(v)
	check := atoi(string(v[vlen-1]))

	sum, mul := 0, 2
	for i := vlen - 2; i >= 0; i-- {
		prod := atoi(v[i:i+1]) * mul
		if prod >= 10 {
			sum += (prod % 10) + 1
		} else {
			sum += prod
		}

		if mul == 2 {
			mul = 1
		} else {
			mul = 2
		}
	}

	return check == ((10 - (sum % 10)) % 10)
}

var rxIPv6Block = regexp.MustCompile(`^(?i)[0-9A-F]{1,4}$`)

// IP reports whether or not v is a valid IP address. The ver argument specifies
// the IP's expected version. The ver argument can be one of the following three
// values:
//	0 accepts both IPv4 and IPv6
//	4 accepts IPv4 only
//	6 accepts IPv6 only
//
// ---
//
//	isvalid:rule
//	{
//		"name": "ip",
//		"opts": [[
//			{ "key": null, "value": "0" },
//			{ "key": "v4", "value": "4" },
//			{ "key": "v6", "value": "6" }
//		]],
//		"err": { "text": "must be a valid IP" }
//	}
func IP(v string, ver int) bool {
	if ver == 0 {
		return IP(v, 4) || IP(v, 6)
	}

	if ver == 4 {
		numbers := strings.Split(v, ".")
		if len(numbers) != 4 {
			return false
		}

		for _, num := range numbers {
			if len(num) > 1 && num[0] == '0' {
				return false // fail if leading zero
			}

			i, err := strconv.Atoi(num)
			if err != nil {
				return false
			}

			if i < 0 || i > 255 {
				return false
			}
		}
		return true
	}

	if ver == 6 {
		parts := strings.Split(v, "%")
		// too many parts
		if len(parts) > 2 {
			return false
		}
		// zoneid must not be empty
		if len(parts) == 2 && len(parts[1]) == 0 {
			return false
		}

		address := parts[0]
		if address == "::" {
			return true
		}

		blocks := strings.Split(address, ":")
		hasIPv4 := false
		hasOmissionBlock := false
		wantBlockNum := 8
		for i := 0; i < len(blocks); i++ {
			// check omission block
			if blocks[i] == "" {
				// more than one omission block? fail
				if hasOmissionBlock {
					return false
				}

				// trailing or leading single ":"; fail
				if i+1 >= len(blocks) || (i == 0 && blocks[i+1] != "") {
					return false
				}

				// trailing or leading "::"; skip the next block
				if blocks[i+1] == "" && (i == 0 || i+2 == len(blocks)) {
					i += 1
				}

				hasOmissionBlock = true
				continue
			}

			// check last block for ipv4 mapping
			if len(blocks) == i+1 {
				// At least some OS accept the last 32 bits of an IPv6 address
				// (i.e. 2 of the blocks) in IPv4 notation, and RFC 3493 says
				// that '::ffff:a.b.c.d' is valid for IPv4-mapped IPv6 addresses,
				// and '::a.b.c.d' is deprecated, but also valid.
				if hasIPv4 = IP(blocks[i], 4); hasIPv4 {
					wantBlockNum = 7
					continue
				}
			}

			// check individual block
			if !rxIPv6Block.MatchString(blocks[i]) {
				return false
			}
		}

		if hasOmissionBlock {
			return len(blocks) >= 1 && len(blocks) < wantBlockNum
		}
		return len(blocks) == wantBlockNum
	}

	return false
}

// IPRange reports whether or not v is a valid IP range (IPv4 only).
//
//	isvalid:rule
//	{
//		"name": "iprange",
//		"err": { "text": "must be a valid IP range" }
//	}
func IPRange(v string) bool {
	parts := strings.Split(v, "/")
	if len(parts) != 2 {
		return false
	}

	subnet := parts[1]
	if len(subnet) < 1 || len(subnet) > 2 || !rxDigits.MatchString(subnet) {
		return false
	}

	if len(subnet) > 1 && subnet[0] == '0' {
		return false
	}

	num := atoi(subnet)
	return IP(parts[0], 4) && num >= 0 && num <= 32
}

var rxISBN10Maybe = regexp.MustCompile(`^(?:[0-9]{9}X|[0-9]{10})$`)

// ISBN reports whether or not v is a valid International Standard Book Number.
// The ver argument specifies the ISBN's expected version. The ver argument can
// be one of the following three values:
//	0 accepts both 10 and 13
//	10 accepts version 10 only
//	13 accepts version 13 only
//
// ---
//
//	isvalid:rule
//	{
//		"name": "isbn",
//		"opts": [[ { "key": null, "value": "0" } ]],
//		"err": { "text": "must be a valid ISBN" }
//	}
func ISBN(v string, ver int) bool {
	if ver == 0 {
		return ISBN(v, 10) || ISBN(v, 13)
	}

	v = rmchar(v, func(r rune) bool { return r == ' ' || r == '-' })

	if ver == 10 {
		if len(v) != 10 || !rxISBN10Maybe.MatchString(v) {
			return false
		}

		tmp, sum := 0, 0
		for i := 0; i < 9; i++ {
			tmp += atoi(string(v[i]))
			sum += tmp
		}

		if v[9] == 'X' {
			tmp += 10
		} else {
			tmp += atoi(string(v[9]))
		}
		sum += tmp

		return sum%11 == 0
	}

	if ver == 13 {
		if len(v) != 13 || !rxDigits.MatchString(v) {
			return false
		}

		sum := 0
		for i := 0; i < 12; i++ {
			if i%2 == 0 {
				sum += atoi(string(v[i]))
			} else {
				sum += 3 * atoi(string(v[i]))
			}
		}

		r := 10 - (sum % 10)
		d := atoi(string(v[12]))
		return (r == 10 && d == 0) || (r < 10 && r == d)
	}
	return false
}

var rxISIN = regexp.MustCompile(`^[A-Z]{2}[0-9A-Z]{9}[0-9]$`)

// ISIN reports whether or not v is a valid International Securities Identification Number.
//
//	isvalid:rule
//	{
//		"name": "isin",
//		"err": { "text": "must be a valid ISIN" }
//	}
func ISIN(v string) bool {
	if !rxISIN.MatchString(v) {
		return false
	}

	ints := make([]int, 0, len(v))
	for _, r := range v {
		if r >= 'A' && r <= 'Z' {
			i64, err := strconv.ParseInt(string(r), 36, 32)
			if err != nil {
				return false
			}
			if i64 >= 10 {
				d1 := int(i64 / 10)
				ints = append(ints, d1)
			}

			d2 := int(i64 % 10)
			ints = append(ints, d2)
		} else {
			ints = append(ints, atoi(string(r)))
		}
	}

	sum := 0
	double := true
	for i := len(ints) - 2; i >= 0; i-- {
		num := ints[i]
		if double {
			num *= 2
			if num >= 10 {
				sum += num + 1
			} else {
				sum += num
			}
		} else {
			sum += num
		}

		double = !double
	}

	return (10000-sum)%10 == ints[len(ints)-1]
}

// ISO639 reports whether or not v is a valid language code as defined by the
// ISO 639 set of standards. Currently only standards 639-1 & 639-2 are supported.
// The num argument specifies which of the supported standards should be tested.
// The num argument can be one of the following three values:
//	0 tests against both 639-1 & 639-2
//	1 tests against 639-1 only
//	2 tests against 639-2 only
//
// ---
//
//	isvalid:rule
//	{
//		"name": "iso369",
//		"opts": [[ { "key": null, "value": "0" } ]],
//		"err": { "text": "must be a valid ISO 639 value" }
//	}
func ISO639(v string, num int) bool {
	if num == 0 {
		if len(v) == 2 && ISO639(v, 1) {
			return true
		}
		if len(v) == 3 && ISO639(v, 2) {
			return true
		}
		return false
	}

	if num == 1 && len(v) == 2 {
		v = strings.ToLower(v)
		_, ok := tables.ISO639_1[v]
		return ok
	}

	if num == 2 && len(v) == 3 {
		v = strings.ToLower(v)
		_, ok := tables.ISO639_2[v]
		return ok
	}

	return false
}

// ISO31661A reports whether or not v is a valid country code as defined by the
// ISO 3166-1 Alpha standard. The num argument specifies which of the two alpha
// sets of the standard should be tested. The num argument can be one of the
// following three values:
//	0 tests against both Alpha-2 and Alpha-3
//	2 tests against Alpha-2 only
//	3 tests against Alpha-3 only
//
// ---
//
//	isvalid:rule
//	{
//		"name": "iso31661a",
//		"opts": [[ { "key": null, "value": "0" } ]],
//		"err": { "text": "must be a valid ISO 3166-1 Alpha value" }
//	}
func ISO31661A(v string, num int) bool {
	if num == 0 {
		if len(v) == 2 && ISO31661A(v, 2) {
			return true
		}
		if len(v) == 3 && ISO31661A(v, 3) {
			return true
		}
		return false
	}

	if num == 2 && len(v) == 2 {
		v = strings.ToUpper(v)
		_, ok := country.ISO31661A_2[v]
		return ok
	}

	if num == 3 && len(v) == 3 {
		v = strings.ToUpper(v)
		_, ok := country.ISO31661A_3[v]
		return ok
	}

	return false
}

// ISO4217 reports whether or not v is a valid currency code as defined by the ISO 4217 standard.
//
//	isvalid:rule
//	{
//		"name": "iso4217",
//		"err": { "text": "must be a valid ISO 4217 value" }
//	}
func ISO4217(v string) bool {
	if len(v) == 3 {
		v = strings.ToUpper(v)
		_, ok := tables.ISO4217[v]
		return ok
	}

	return false
}

var rxISRC = regexp.MustCompile(`^[A-Z]{2}[0-9A-Z]{3}\d{2}\d{5}$`)

// ISRC reports whether or not v is a valid International Standard Recording Code.
//
//	isvalid:rule
//	{
//		"name": "isrc",
//		"err": { "text": "must be a valid ISRC" }
//	}
func ISRC(v string) bool {
	return rxISRC.MatchString(v)
}

var rxISSN = regexp.MustCompile(`^(?i)[0-9]{4}-?[0-9]{3}[0-9X]$`)

// ISSN reports whether or not v is a valid International Standard Serial Number.
//
//	isvalid:rule
//	{
//		"name": "issn",
//		"err": { "text": "must be a valid ISSN" }
//	}
func ISSN(v string, requireHyphen, caseSensitive bool) bool {
	if requireHyphen && !strings.Contains(v, "-") {
		return false
	}
	if caseSensitive && strings.Contains(v, "x") {
		return false
	}
	if !rxISSN.MatchString(v) {
		return false
	}

	v = rmchar(v, func(r rune) bool { return r == '-' })
	v = strings.ToUpper(v)

	sum := 0
	for i := 0; i < len(v); i++ {
		if v[i] == 'X' {
			sum += 10 * (8 - i)
		} else {
			sum += atoi(string(v[i])) * (8 - i)
		}
	}
	return sum%11 == 0
}

// In reports whether or not v is in the provided list.
//
//	isvalid:rule
//	{
//		"name": "in",
//		"err": { "text": "must be in the list" }
//	}
func In(v interface{}, list ...interface{}) bool {
	for _, item := range list {
		if v == item {
			return true
		}
	}
	return false
}

var rxInt = regexp.MustCompile(`^[+-]?[0-9]+$`)

// Int reports whether or not v represents a valid integer.
//
//	isvalid:rule
//	{
//		"name": "int",
//		"err": { "text": "string content must match an integer" }
//	}
func Int(v string) bool {
	return rxInt.MatchString(v)
}

// JSON reports whether or not v is a valid JSON value. NOTE: the input is validated
// using json.Unmarshal which accepts primitive values as long as they are JSON values.
//
//	isvalid:rule
//	{
//		"name": "json",
//		"err": { "text": "must be a valid JSON" }
//	}
func JSON(v []byte) bool {
	var x interface{}
	return json.Unmarshal(v, &x) == nil
}

// JWT reports whether or not v is a valid JSON Web Token.
//
//	isvalid:rule
//	{
//		"name": "jwt",
//		"err": { "text": "must be a valid JWT" }
//	}
func JWT(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	for _, part := range parts {
		if !Base64(part, true) {
			return false
		}
	}
	return true
}

var rxLat = regexp.MustCompile(`^\(?[+-]?(?:90(?:\.0+)?|[1-8]?\d(?:\.\d+)?)$`)
var rxLong = regexp.MustCompile(`^\s?[+-]?(?:180(?:\.0+)?|1[0-7]\d(?:\.\d+)?|\d{1,2}(?:\.\d+)?)\)?$`)

var rxLatDMS = regexp.MustCompile(`^(?:(?:[1-8]?\d)\D+(?:[1-5]?\d|60)\D+(?:[1-5]?\d|60)(?:\.\d+)?|90\D+0\D+0)\D+[NSns]?$`)
var rxLongDMS = regexp.MustCompile(`^\s*(?:[1-7]?\d{1,2}\D+(?:[1-5]?\d|60)\D+(?:[1-5]?\d|60)(?:\.\d+)?|180\D+0\D+0)\D+[EWew]?$`)

// LatLong reports whether or not v is a valid latitude-longitude coordinate string.
//
//	isvalid:rule
//	{
//		"name": "latlong",
//		"opts": [[
//			{ "key": null, "value": "false" },
//			{ "key": "dms", "value": "true" }
//		]],
//		"err": { "text": "must be a valid latitude-longitude coordinate" }
//	}
func LatLong(v string, dms bool) bool {
	pair := strings.Split(v, ",")
	if len(pair) != 2 || len(pair[0]) == 0 || len(pair[1]) == 0 {
		return false
	}
	if (pair[0][0] == '(' && pair[1][len(pair[1])-1] != ')') ||
		(pair[0][0] != '(' && pair[1][len(pair[1])-1] == ')') {
		return false
	}
	if dms {
		return rxLatDMS.MatchString(pair[0]) && rxLongDMS.MatchString(pair[1])
	}
	return rxLat.MatchString(pair[0]) && rxLong.MatchString(pair[1])
}

var rxLocale = regexp.MustCompile(`^[A-z]{2,4}(?:[_-](?:[A-z]{4}|[\d]{3}))?(?:[_-](?:[A-z]{2}|[\d]{3}))?$`)

// Locale reports whether or not v is a valid locale.
//
//	isvalid:rule
//	{
//		"name": "locale",
//		"err": { "text": "must be a valid locale" }
//	}
func Locale(v string) bool {
	if v == "en_US_POSIX" || v == "ca_ES_VALENCIA" {
		return true
	}
	return rxLocale.MatchString(v)
}

// LowerCase reports whether or not v is an all lower-case string.
//
//	isvalid:rule
//	{
//		"name": "lower",
//		"err": { "text": "must contain only lower-case characters" }
//	}
func LowerCase(v string) bool {
	return v == strings.ToLower(v)
}

var rxMAC6 = regexp.MustCompile(`^[0-9a-fA-F]{12}$`)
var rxMAC6Colon = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
var rxMAC6Hyphen = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}-){5}[0-9a-fA-F]{2}$`)

var rxMAC8 = regexp.MustCompile(`^[0-9a-fA-F]{16}$`)
var rxMAC8Colon = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}:){7}[0-9a-fA-F]{2}$`)
var rxMAC8Hyphen = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}-){7}[0-9a-fA-F]{2}$`)

// MAC reports whether or not v is a valid MAC address. The space argument specifies
// the identifier's expected address space (in bytes). The space argument can be one
// of the following three values:
//	0 accepts both EUI-48 and EUI-64
//	6 accepts EUI-48 format only
//	8 accepts EUI-64 format only
//
// The allowed formatting of the identifiers is as follows:
//	// EUI-48 format
//	"08:00:2b:01:02:03"
//	"08-00-2b-01-02-03"
//	"08002b010203"
//
//	// EUI-64 format
//	"08:00:2b:01:02:03:04:05"
//	"08-00-2b-01-02-03-04-05"
//	"08002b0102030405"
//
// ---
//
//	isvalid:rule
//	{
//		"name": "mac",
//		"opts": [[ { "key": null, "value": "0" } ]],
//		"err": { "text": "must be a valid MAC" }
//	}
func MAC(v string, space int) bool {
	if space == 0 {
		return rxMAC6.MatchString(v) ||
			rxMAC6Colon.MatchString(v) ||
			rxMAC6Hyphen.MatchString(v) ||
			rxMAC8.MatchString(v) ||
			rxMAC8Colon.MatchString(v) ||
			rxMAC8Hyphen.MatchString(v)
	} else if space == 6 {
		return rxMAC6.MatchString(v) ||
			rxMAC6Colon.MatchString(v) ||
			rxMAC6Hyphen.MatchString(v)
	} else if space == 8 {
		return rxMAC8.MatchString(v) ||
			rxMAC8Colon.MatchString(v) ||
			rxMAC8Hyphen.MatchString(v)
	}
	return false
}

// MD5 reports whether or not v is a valid MD5 hash.
//
//	isvalid:rule
//	{
//		"name": "md5",
//		"err": { "text": "must be a valid MD5 hash" }
//	}
func MD5(v string) bool {
	return len(v) == 32 && rxHash.MatchString(v)
}

var rxMIMESimple = regexp.MustCompile(`^(?i)(?:application|audio|font|image|message|model|multipart|text|video)\/[a-zA-Z0-9\.\-\+]{1,100}$`)
var rxMIMEText = regexp.MustCompile(`^(?i)text\/[a-zA-Z0-9\.\-\+]{1,100};\s?charset=(?:"[a-zA-Z0-9\.\-\+\s]{0,70}"|[a-zA-Z0-9\.\-\+]{0,70})(?:\s?\([a-zA-Z0-9\.\-\+\s]{1,20}\))?$`)
var rxMIMEMultipart = regexp.MustCompile(`^(?i)multipart\/[a-zA-Z0-9\.\-\+]{1,100}(?:;\s?(?:boundary|charset)=(?:"[a-zA-Z0-9\.\-\+\s]{0,70}"|[a-zA-Z0-9\.\-\+]{0,70})(?:\s?\([a-zA-Z0-9\.\-\+\s]{1,20}\))?){0,2}$`)

// MIME reports whether or not v is of a valid MIME type format.
//
// NOTE: This function only checks is the string format follows the etablished
// rules by the according RFC specifications. This function supports 'charset'
// in textual media types (https://tools.ietf.org/html/rfc6657).
//
// This function does not check against all the media types listed by the IANA
// (https://www.iana.org/assignments/media-types/media-types.xhtml)
//
// More informations in the RFC specifications:
//	- https://tools.ietf.org/html/rfc2045
// 	- https://tools.ietf.org/html/rfc2046
// 	- https://tools.ietf.org/html/rfc7231#section-3.1.1.1
// 	- https://tools.ietf.org/html/rfc7231#section-3.1.1.5
//
// ---
//
//	isvalid:rule
//	{
//		"name": "mime",
//		"err": { "text": "must be a valid media type" }
//	}
func MIME(v string) bool {
	return rxMIMESimple.MatchString(v) || rxMIMEText.MatchString(v) || rxMIMEMultipart.MatchString(v)
}

var rxMagnetURI = regexp.MustCompile(`^(?i)magnet:\?xt=urn:[a-z0-9]+:[a-z0-9]{32,40}&dn=.+&tr=.+$`)

// MagnetURI reports whether or not v is a valid magned URI.
//
//	isvalid:rule
//	{
//		"name": "magneturi",
//		"err": { "text": "must be a valid magnet URI" }
//	}
func MagnetURI(v string) bool {
	return rxMagnetURI.MatchString(v)
}

// Match reports whether or not v contains any match of the *registered* regular
// expression re. NOTE: re MUST be registered upfront with RegisterRegexp otherwise
// Match will return false even if v matches the regular expression.
//
//	isvalid:rule
//	{
//		"name": "re",
//		"err": {
//			"text": "must match the regular expression",
//			"with_opts": true
//		}
//	}
func Match(v string, re string) bool {
	if rx, ok := rxcache.m[re]; ok {
		return rx.MatchString(v)
	}
	return false
}

// MongoId reports whether or not v is a valid hex-encoded representation of a MongoDB ObjectId.
//
//	isvalid:rule
//	{
//		"name": "mongoid",
//		"err": { "text": "must be a valid Mongo Object Id" }
//	}
func MongoId(v string) bool {
	return len(v) == 24 && rxHex.MatchString(v)
}

var rxNumeric = regexp.MustCompile(`^[+-]?[0-9]*\.?[0-9]+$`)

// Numeric reports whether or not v is a valid numeric string.
//
//	isvalid:rule
//	{
//		"name": "numeric",
//		"err": { "text": "string content must match a numeric value" }
//	}
func Numeric(v string) bool {
	return rxNumeric.MatchString(v)
}

var rxOctal = regexp.MustCompile(`^(?:0[oO])?[0-7]+$`)

// Octal reports whether or not v represents a valid octal integer.
//
//	isvalid:rule
//	{
//		"name": "octal",
//		"err": { "text": "string content must match an octal number" }
//	}
func Octal(v string) bool {
	return rxOctal.MatchString(v)
}

var rxPAN = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3,6})?|5[1-5][0-9]{14}|(?:222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}|6(?:011|5[0-9][0-9])[0-9]{12,15}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11}|6[27][0-9]{14})$`)

// PAN reports whether or not v is a valid Primary Account Number or Credit Card number.
//
//	isvalid:rule
//	{
//		"name": "pan",
//		"err": { "text": "must be a valid PAN" }
//	}
func PAN(v string) bool {
	v = rmchar(v, func(r rune) bool { return r == ' ' || r == '-' })
	if !rxPAN.MatchString(v) {
		return false
	}

	// luhn check
	var sum int
	var double bool
	for i := len(v) - 1; i >= 0; i-- {
		num := atoi(string(v[i]))

		if double {
			num *= 2
			if num > 9 {
				num = (num % 10) + 1
			}
		}
		double = !double

		sum += num
	}
	return sum%10 == 0
}

// PassportNumber reports whether or not v is a valid passport number.
//
//	isvalid:rule
//	{
//		"name": "passport",
//		"err": { "text": "must be a valid passport number" }
//	}
func PassportNumber(v string) bool {
	// TODO
	return false
}

// Phone reports whether or not v is a valid phone number in the country
// identified by the given country code cc.
//
//	isvalid:rule
//	{
//		"name": "phone",
//		"opts": [[ { "key": null, "value": "us" } ]],
//		"err": { "text": "must be a valid phone number" }
//	}
func Phone(v string, cc string) bool {
	if c, ok := country.Get(cc); ok && c.Phone != nil {
		return c.Phone.MatchString(v)
	}
	return false
}

// Port reports whether or not v is a valid port number.
//
//	isvalid:rule
//	{
//		"name": "port",
//		"err": { "text": "must be a valid port number" }
//	}
func Port(v string) bool {
	_, err := strconv.ParseUint(v, 10, 16)
	return err == nil
}

var rxRGB = regexp.MustCompile(`^rgb\((?:(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5]),){2}(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\)$`)
var rxRGBA = regexp.MustCompile(`^rgba\((?:(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5]),){3}(?:0?\.\d|1(?:\.0)?|0(?:\.0)?)\)$`)
var rxRGBPercent = regexp.MustCompile(`^rgb\((?:(?:[0-9]%|[1-9][0-9]%|100%),){2}(?:[0-9]%|[1-9][0-9]%|100%)\)`)
var rxRGBAPercent = regexp.MustCompile(`^rgba\((?:(?:[0-9]%|[1-9][0-9]%|100%),){3}(?:0?\.\d|1(?:\.0)?|0(?:\.0)?)\)`)

// RGB reports whether or not v is a valid RGB color.
//
//	isvalid:rule
//	{
//		"name": "rgb",
//		"err": { "text": "must be a valid RGB color" }
//	}
func RGB(v string) bool {
	if len(v) >= 4 && v[:4] == "rgba" {
		if strings.Contains(v, "%") {
			return rxRGBAPercent.MatchString(v)
		} else {
			return rxRGBA.MatchString(v)
		}
	} else {
		if strings.Contains(v, "%") {
			return rxRGBPercent.MatchString(v)
		} else {
			return rxRGB.MatchString(v)
		}
	}

	// MAYBE-TODO this seems like it could be handled by a simple parser rather
	// than the excessive regexp (https://en.wikipedia.org/wiki/Web_colors#CSS_colors)
	return false
}

// SSN reports whether or not v has a valid Social Security Number format.
//
//	isvalid:rule
//	{
//		"name": "ssn",
//		"err": { "text": "must be a valid SSN" }
//	}
func SSN(v string) bool {
	if len(v) == 11 {
		if v[3] != '-' || v[6] != '-' {
			return false
		}
		v = v[:3] + v[4:6] + v[7:]
	}
	if len(v) != 9 {
		return false
	}

	// must be digits only
	for _, r := range v {
		if r < '0' || r > '9' {
			return false
		}
	}

	area, group, serial := v[:3], v[3:5], v[5:]
	// no digit group can have all zeroes
	if area == "000" || group == "00" || serial == "0000" {
		return false
	}
	// area cannot be 666 or 900 and above
	if area == "666" || area[0] == '9' {
		return false
	}

	return true
}

var rxSemVer = regexp.MustCompile(`^(?i)(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)` +
	`(?:-(?:(?:0|[1-9]\d*|\d*[a-z-][0-9a-z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-z-][0-9a-z-]*))*))` +
	`?(?:\+(?:[0-9a-z-]+(?:\.[0-9a-z-]+)*))?$`)

// SemVer reports whether or not v is a valid Semantic Versioning number.
// For reference: https://semver.org/
//
//	isvalid:rule
//	{
//		"name": "semver",
//		"err": { "text": "must be a valid semver number" }
//	}
func SemVer(v string) bool {
	return rxSemVer.MatchString(v)
}

var rxSlug = regexp.MustCompile(`^(?:(?:[a-z0-9]+)(?:-[a-z0-9]+)?)+$`)

// Slug reports whether or not v is a valid slug.
//
//	isvalid:rule
//	{
//		"name": "slug",
//		"err": { "text": "must be a valid slug" }
//	}
func Slug(v string) bool {
	return rxSlug.MatchString(v)
}

type StrongPasswordOpts struct {
	MinLen     int
	MinLower   int
	MinUpper   int
	MinNumbers int
	MinSymbols int
}

var StrongPasswordOptsDefault = StrongPasswordOpts{
	MinLen:     8,
	MinLower:   1,
	MinUpper:   1,
	MinNumbers: 1,
	MinSymbols: 1,
}

// StrongPassword reports whether or not v is a strong password.
//
//	isvalid:rule
//	{
//		"name": "strongpass",
//		"opts": [[ { "key": null, "value": "nil" } ]],
//		"err": { "text": "must be a strong password" }
//	}
func StrongPassword(v string, opts *StrongPasswordOpts) bool {
	if opts == nil {
		opts = &StrongPasswordOptsDefault
	}

	if len(v) < opts.MinLen {
		return false
	}

	var lo, up, num, sym int
	for _, r := range v {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				up += 1
			} else {
				lo += 1
			}
		} else if unicode.IsNumber(r) {
			num += 1
		} else {
			sym += 1
		}
	}

	return lo >= opts.MinLower && up >= opts.MinUpper &&
		num >= opts.MinNumbers && sym >= opts.MinSymbols
}

// URL reports whether or not v is a valid Uniform Resource Locator.
//
//	isvalid:rule
//	{
//		"name": "url",
//		"err": { "text": "must be a valid URL" }
//	}
func URL(v string) bool {
	// TODO consider implementing a parser, maybe port https://github.com/servo/rust-url to Go
	return false
}

var rxUUIDv3 = regexp.MustCompile(`^(?i)[0-9A-F]{8}-[0-9A-F]{4}-3[0-9A-F]{3}-[0-9A-F]{4}-[0-9A-F]{12}$`)
var rxUUIDv4 = regexp.MustCompile(`^(?i)[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`)
var rxUUIDv5 = regexp.MustCompile(`^(?i)[0-9A-F]{8}-[0-9A-F]{4}-5[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`)

// UUID reports whether or not v is a valid Universally Unique IDentifier.
// NOTE: only versions 3, 4, and 5 are currently supported.
//
//	isvalid:rule
//	{
//		"name": "uuid",
//		"opts": [[ { "key": null, "value": "4" } ]],
//		"err": { "text": "must be a valid UUID" }
//	}
func UUID(v string, ver int) bool {
	// TODO add the rest of the versions?
	switch ver {
	case 3:
		return rxUUIDv3.MatchString(v)
	case 4:
		return rxUUIDv4.MatchString(v)
	case 5:
		return rxUUIDv5.MatchString(v)
	}
	return false
}

var rxUint = regexp.MustCompile(`^\+?[0-9]+$`)

// Uint reports whether or not v represents a valid unsigned integer.
//
//	isvalid:rule
//	{
//		"name": "uint",
//		"err": { "text": "string content must match an unsigned integer" }
//	}
func Uint(v string) bool {
	return rxUint.MatchString(v)
}

// UpperCase reports whether or not v is an all upper-case string.
//
//	isvalid:rule
//	{
//		"name": "upper",
//		"err": { "text": "must contain only upper-case characters" }
//	}
func UpperCase(v string) bool {
	return v == strings.ToUpper(v)
}

// VAT reports whether or not v is a valid Value Added Tax number.
//
//	isvalid:rule
//	{
//		"name": "vat",
//		"err": { "text": "must be a valid VAT number" }
//	}
func VAT(v string) bool {
	// TODO https://en.wikipedia.org/wiki/VAT_identification_number
	return false
}

// Zip reports whether or not v is a valid zip / postal code for the country
// identified by the given country code cc.
//
//	isvalid:rule
//	{
//		"name": "zip",
//		"opts": [[ { "key": null, "value": "us" } ]],
//		"err": { "text": "must be a valid zip code" }
//	}
func Zip(v string, cc string) bool {
	if c, ok := country.Get(cc); ok && c.Zip != nil {
		return c.Zip.MatchString(v)
	}
	return false
}
