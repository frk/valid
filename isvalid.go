package isvalid

import (
	"encoding/base64"
	"net/mail"
	"regexp"
	"strings"
)

var rxInt = regexp.MustCompile(`^[+-]?[0-9]+$`)

// Int reports whether or not v represents a valid integer.
//
// isvalid:rule
//	{ "name": "int", "err": {"text": "string content must match an integer"} }
func Int(v string) bool {
	return rxInt.MatchString(v)
}

var rxUint = regexp.MustCompile(`^\+?[0-9]+$`)

// Uint reports whether or not v represents a valid unsigned integer.
//
// isvalid:rule
//	{ "name": "uint", "err": {"text": "string content must match an unsigned integer"} }
func Uint(v string) bool {
	return rxUint.MatchString(v)
}

var rxFloat = regexp.MustCompile(`^[+-]?(?:[0-9]*)?(?:\.[0-9]*)?(?:[eE][+-]?[0-9]+)?$`)

// Float reports whether or not v represents a valid float.
//
// isvalid:rule
//	{ "name": "float", "err": {"text": "string content must match a floating point number"} }
func Float(v string) bool {
	if v == "" || v == "." || v == "+" || v == "-" {
		return false
	}
	return rxFloat.MatchString(v)
}

// Bool reports whether or not v represents a valid boolean. The following
// are considered valid boolean values: "true", "false", "TRUE", and "FALSE".
//
// isvalid:rule
//	{ "name": "bool", "err": {"text": "string content must match a boolean value"} }
func Bool(v string) bool {
	if len(v) == 4 && (v == "true" || v == "TRUE") {
		return true
	}
	if len(v) == 5 && (v == "false" || v == "FALSE") {
		return true
	}
	return false
}

var rxBinary = regexp.MustCompile(`^(?:0[bB])?[0-1]+$`)

// Binary reports whether or not v represents a valid binary integer.
//
// isvalid:rule
//	{ "name": "binary", "err": {"text": "string content must match a binary number"} }
func Binary(v string) bool {
	return rxBinary.MatchString(v)
}

var rxOctal = regexp.MustCompile(`^(?:0[oO])?[0-7]+$`)

// Octal reports whether or not v represents a valid octal integer.
//
// isvalid:rule
//	{ "name": "octal", "err": {"text": "string content must match an octal number"} }
func Octal(v string) bool {
	return rxOctal.MatchString(v)
}

var rxHex = regexp.MustCompile(`^(?:0[xXhH])?[0-9A-Fa-f]+$`)

// Hex reports whether or not v is a valid hexadecimal string.
//
// isvalid:rule
//	{ "name": "hex", "err": {"text": "must be a valid hexadecimal string"} }
func Hex(v string) bool {
	return rxHex.MatchString(v)
}

var rxNumeric = regexp.MustCompile(`^[+-]?[0-9]*\.?[0-9]+$`)

// Numeric reports whether or not v is a valid numeric string.
//
// isvalid:rule
//	{ "name": "numeric", "err": {"text": "string content must match a numeric value"} }
func Numeric(v string) bool {
	return rxNumeric.MatchString(v)
}

var rxDigits = regexp.MustCompile(`^[0-9]+$`)

// Digits reports whether or not v is a string of digits.
//
// isvalid:rule
//	{ "name": "digits", "err": {"text": "must contain only digits"} }
func Digits(v string) bool {
	return rxDigits.MatchString(v)
}

var rxHexColor = regexp.MustCompile(`^#?(?i:[0-9A-F]{3}|[0-9A-F]{4}|[0-9A-F]{6}|[0-9A-F]{8})$`)

// HexColor reports whether or not v is a valid hexadecimal color code.
//
// isvalid:rule
//	{ "name": "hexcolor", "err": {"text": "must represent a valid hexadecimal color code"} }
func HexColor(v string) bool {
	return rxHexColor.MatchString(v)
}

var rxBase32 = regexp.MustCompile(`^[A-Z2-7]+=*$`)

// Base32 reports whether or not v is a valid base32 string.
//
// isvalid:rule
//	{ "name": "base32", "err": {"text": "must be a valid base32 string"} }
func Base32(v string) bool {
	return (len(v)%8 == 0) && rxBase32.MatchString(v)
}

var rxBase58 = regexp.MustCompile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$`)

// Base58 reports whether or not v is a valid base58 string.
//
// isvalid:rule
//	{ "name": "base58", "err": {"text": "must be a valid base58 string"} }
func Base58(v string) bool {
	return rxBase58.MatchString(v)
}

// Base64 reports whether or not v is a valid base64 string. NOTE: The standard
// "encoding/base64" package is used for validation. With urlsafe=false StdEncoding
// is used and with urlsafe=true RawURLEncoding is used.
//
// isvalid:rule
//	{
// 	     "name": "base64",
// 	     "arg_map": [{
//		"key": null,
//		"value": "false"
//	     }, {
//		"key": "url",
//		"value": "true"
//	     }],
//	     "err": {"text": "must be a valid base64 string"}
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

// LowerCase reports whether or not v is an all lower-case string.
//
// isvalid:rule
//	{ "name": "lower", "err": {"text": "must contain only lower-case characters"} }
func LowerCase(v string) bool {
	return v == strings.ToLower(v)
}

// UpperCase reports whether or not v is an all upper-case string.
//
// isvalid:rule
//	{ "name": "upper", "err": {"text": "must contain only upper-case characters"} }
func UpperCase(v string) bool {
	return v == strings.ToUpper(v)
}

// Email reports whether or not v is a valid email address. NOTE: Email uses
// "net/mail".ParseAddress to determine the validity of v.
//
// isvalid:rule
//	{ "name": "email", "err": {"text": "must be a valid email address"} }
func Email(v string) bool {
	_, err := mail.ParseAddress(v)
	return err == nil
}

// URL reports whether or not v is a valid Uniform Resource Locator.
//
// isvalid:rule
//	{ "name": "url", "err": {"text": "must be a valid URL"} }
func URL(v string) bool {
	return false
}

// URI reports whether or not v is a valid Uniform Resource Identifier.
//
// isvalid:rule
//	{ "name": "uri", "err": {"text": "must be a valid URI"} }
func URI(v string) bool {
	return false
}

// PAN reports whether or not v is a valid Primary Account Number.
//
// isvalid:rule
//	{ "name": "pan", "err": {"text": "must be a valid PAN"} }
func PAN(v string) bool {
	return false
}

// CVV reports whether or not v is a valid Card Verification Value.
//
// isvalid:rule
//	{ "name": "cvv", "err": {"text": "must be a valid CVV"} }
func CVV(v string) bool {
	return false
}

// SSN reports whether or not v is a valid Social Security Number.
//
// isvalid:rule
//	{ "name": "ssn", "err": {"text": "must be a valid SSN"} }
func SSN(v string) bool {
	return false
}

var rxEIN = regexp.MustCompile(`^\d{2}[- ]{0,1}\d{7}$`)

// EIN reports whether or not v is a valid Employer Identification Number.
//
// isvalid:rule
//	{ "name": "ein", "err": {"text": "must be a valid EIN"} }
func EIN(v string) bool {
	return rxEIN.MatchString(v)
}

var (
	rxAlpha = map[string]*regexp.Regexp{
		"be":  regexp.MustCompile(`^(?i:[АБВГДЕЁЖЗІЙКЛМНОПРСТУЎФХЦЧШЫЬЭЮЯ]+)$`),
		"bg":  regexp.MustCompile(`^(?i:[АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ]+)$`),
		"cnr": regexp.MustCompile(`^(?i:[A-ZČĆĐŠŚŽŹАБВГДЂЕЖЗЗ́ИЈКЛЉМНЊОПРСС́ТЋУФХЦЧЏШ]+)$`),
		"cs":  regexp.MustCompile(`^(?i:[A-ZÁČĎÉĚÍŇÓŘŠŤÚŮÝŽ]+)$`),
		"en":  regexp.MustCompile(`^(?i:[A-Z]+)$`),
		"mk":  regexp.MustCompile(`^(?i:[АБВГДЃЕЖЗЅИЈКЛЉМНЊОПРСТЌУФХЦЧЏШ]+)$`),
		"pl":  regexp.MustCompile(`^(?i:[AĄBCĆDEĘFGHIJKLŁMNŃOÓPQRSŚTUVWXYZŹŻ]+)$`),
		"ru":  regexp.MustCompile(`^(?i:[АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ]+)$`),
		// TODO "sh" is deprecated, use "hbs"
		"sh": regexp.MustCompile(`^(?i:[ABCČĆDDžĐEFGHIJKLLMNNOPRSŠTUVZŽ]+)$`),
		"sk": regexp.MustCompile(`^(?i:[A-ZÁÄČĎÉÍĹĽŇÓÔŔŠŤÚÝŽ]+)$`),
		"sl": regexp.MustCompile(`^(?i:[ABCČDEFGHIJKLMNOPRSŠTUVZŽ]+)$`),
		"sr": regexp.MustCompile(`^(?i:[АБВГДЂЕЖЗИЈКЛЉМНЊОПРСТЋУФХЦЧЏШ]+)$`),
		"uk": regexp.MustCompile(`^(?i:[АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ]+)$`),
		// "wen" sorbian languages
		"wen": regexp.MustCompile(`^(?i:[ABCČĆDEĚFGHIJKŁLMNŃOÓPRŘŔSŠŚTUWYZŽŹ]+)$`),
	}

	rxAlnum = map[string]*regexp.Regexp{
		"be":  regexp.MustCompile(`^(?i:[0-9АБВГДЕЁЖЗІЙКЛМНОПРСТУЎФХЦЧШЫЬЭЮЯ]+)$`),
		"bg":  regexp.MustCompile(`^(?i:[0-9АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ]+)$`),
		"cnr": regexp.MustCompile(`^(?i:[0-9A-ZČĆĐŠŚŽŹАБВГДЂЕЖЗЗ́ИЈКЛЉМНЊОПРСС́ТЋУФХЦЧЏШ]+)$`),
		"cs":  regexp.MustCompile(`^(?i:[0-9A-ZÁČĎÉĚÍŇÓŘŠŤÚŮÝŽ]+)$`),
		"en":  regexp.MustCompile(`^(?i:[0-9A-Z]+)$`),
		"mk":  regexp.MustCompile(`^(?i:[0-9АБВГДЃЕЖЗЅИЈКЛЉМНЊОПРСТЌУФХЦЧЏШ]+)$`),
		"pl":  regexp.MustCompile(`^(?i:[0-9AĄBCĆDEĘFGHIJKLŁMNŃOÓPQRSŚTUVWXYZŹŻ]+)$`),
		"ru":  regexp.MustCompile(`^(?i:[0-9АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ]+)$`),
		// TODO "sh" is deprecated, use "hbs"
		"sh": regexp.MustCompile(`^(?i:[0-9ABCČĆDDžĐEFGHIJKLLMNNOPRSŠTUVZŽ]+)$`),
		"sk": regexp.MustCompile(`^(?i:[0-9A-ZÁÄČĎÉÍĹĽŇÓÔŔŠŤÚÝŽ]+)$`),
		"sl": regexp.MustCompile(`^(?i:[0-9ABCČDEFGHIJKLMNOPRSŠTUVZŽ]+)$`),
		"sr": regexp.MustCompile(`^(?i:[0-9АБВГДЂЕЖЗИЈКЛЉМНЊОПРСТЋУФХЦЧЏШ]+)$`),
		"uk": regexp.MustCompile(`^(?i:[0-9АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ]+)$`),
		// "wen" sorbian languages
		"wen": regexp.MustCompile(`^(?i:[0-9ABCČĆDEĚFGHIJKŁLMNŃOÓPRŘŔSŠŚTUWYZŽŹ]+)$`),
	}
)

// Alpha reports whether or not v is a valid alphabetic string. The variadic langs
// argument can be used to pass in zero or more language tags that are then used
// to check v's characters against the alphabets of those languages. If no language
// tag is provided Alpha will by default run the check against the english alphabet.
// For Alpha to return true it is enough if v is alphabetic in just one language.
//
// isvalid:rule
//	{ "name": "alpha", "err": {"text": "must be an alphabetic string"} }
func Alpha(v string, langs ...string) bool {
	if len(langs) == 0 {
		if re, ok := rxAlpha["en"]; ok && re.MatchString(v) {
			return true
		}
	}
	for _, lang := range langs {
		if re, ok := rxAlpha[lang]; ok && re.MatchString(v) {
			return true
		}
	}
	return false
}

// Alnum reports whether or not v is a valid alphanumeric string. The variadic langs
// argument can be used to pass in zero or more language tags that are then used
// to check v's characters against the alphabets of those languages. If no language
// tag is provided Alnum will by default run the check against the english alphabet.
// For Alnum to return true it is enough if v is alphanumeric in just one language.
//
// isvalid:rule
//	{ "name": "alnum", "err": {"text": "must be an alphanumeric string"} }
func Alnum(v string, langs ...string) bool {
	if len(langs) == 0 {
		if re, ok := rxAlnum["en"]; ok && re.MatchString(v) {
			return true
		}
	}
	for _, lang := range langs {
		if re, ok := rxAlnum[lang]; ok && re.MatchString(v) {
			return true
		}
	}
	return false
}

// CIDR reports whether or not v is a valid Classless Inter-Domain Routing notation.
//
// isvalid:rule
//	{ "name": "cidr", "err": {"text": "must be a valid CIDR notation"} }
func CIDR(v string) bool {
	return false
}

// Phone reports whether or not v is a valid phone number.
//
// isvalid:rule
//	{ "name": "phone", "err": {"text": "must be a valid phone number"} }
func Phone(v string, cc ...string) bool {
	return false
}

// Zip reports whether or not v is a valid zip code.
//
// isvalid:rule
//	{ "name": "zip", "err": {"text": "must be a valid zip code"} }
func Zip(v string, cc ...string) bool {
	return false
}

// UUID reports whether or not v is a valid Universally Unique IDentifier.
//
// isvalid:rule
//	{ "name": "uuid", "arg_min": 0, "arg_max": 5, "err": {"text": "must be a valid UUID"} }
func UUID(v string, ver ...int) bool {
	return false
}

// IP reports whether or not v is a valid IP address.
//
// isvalid:rule
//	{ "name": "ip", "arg_min": 0, "arg_max": 2, "err": {"text": "must be a valid IP"} }
func IP(v string, ver ...int) bool {
	return false
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
//
//	0 // accepts both EUI-48 and EUI-64
//	6 // accepts EUI-48 format only
//	8 // accepts EUI-64 format only
//
// The allowed formatting of the identifiers is as follows:
//
//	// MAC - EUI-48 format
//	08:00:2b:01:02:03
//	08-00-2b-01-02-03
//	08002b010203
//
//	// MAC - EUI-64 format
//	08:00:2b:01:02:03:04:05
//	08-00-2b-01-02-03-04-05
//	08002b0102030405
//
// isvalid:rule
//	{
//		"name": "mac",
//		"arg_map": [{"key": null, "value": "0"}],
//		"err": {"text": "must be a valid MAC"}
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

// RFC reports whether or not v is a valid representation of the specified RFC standard.
//
// isvalid:rule
//	{ "name": "rfc", "err": {"text": "must be a valid RFC", "with_args": true} }
func RFC(v string, num int) bool {
	return false
}

// ISO reports whether or not v is a valid representation of the specified ISO standard.
//
// isvalid:rule
//	{ "name": "iso", "err": {"text": "must be a valid ISO", "with_args": true} }
func ISO(v string, num int) bool {
	return false
}

// Match reports whether or not the v contains any match of the regular expression re.
// NOTE: Match will panic if re has not been registered upfront with RegisterRegexp.
//
// isvalid:rule
//	{ "name": "re", "err": {"text": "must match the regular expression", "with_args": true} }
func Match(v string, re string) bool {
	return regexpCache.m[re].MatchString(v)
}
