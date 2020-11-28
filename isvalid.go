package isvalid

import (
	"net/mail"
	"regexp"
)

var rxInt = regexp.MustCompile(`^[+-]?[0-9]+$`)

// Int reports whether or not v represents a valid integer.
func Int(v string) bool {
	return rxInt.MatchString(v)
}

var rxUint = regexp.MustCompile(`^\+?[0-9]+$`)

// Uint reports whether or not v represents a valid unsigned integer.
func Uint(v string) bool {
	return rxUint.MatchString(v)
}

var rxFloat = regexp.MustCompile(`^[+-]?(?:[0-9]*)?(?:\.[0-9]*)?(?:[eE][+-]?[0-9]+)?$`)

// Float reports whether or not v represents a valid float.
func Float(v string) bool {
	if v == "" || v == "." || v == "+" || v == "-" {
		return false
	}
	return rxFloat.MatchString(v)
}

var rxBinary = regexp.MustCompile(`^(?:0[bB])?[0-1]+$`)

// Binary reports whether or not v represents a valid binary integer.
func Binary(v string) bool {
	return rxBinary.MatchString(v)
}

var rxOctal = regexp.MustCompile(`^(?:0[oO])?[0-7]+$`)

// Octal reports whether or not v represents a valid octal integer.
func Octal(v string) bool {
	return rxOctal.MatchString(v)
}

var rxHex = regexp.MustCompile(`^(?:0[xXhH])?[0-9A-Fa-f]+$`)

// Hex reports whether or not v is a valid hexadecimal string.
func Hex(v string) bool {
	return rxHex.MatchString(v)
}

var rxNumeric = regexp.MustCompile(`^[+-]?[0-9]*\.?[0-9]+$`)

// Numeric reports whether or not v is a valid numeric string.
func Numeric(v string) bool {
	return rxNumeric.MatchString(v)
}

var rxDigits = regexp.MustCompile(`^[0-9]+$`)

// Digits reports whether or not v is a string of digits.
func Digits(v string) bool {
	return rxDigits.MatchString(v)
}

var rxHexColor = regexp.MustCompile(`^#?(?i:[0-9A-F]{3}|[0-9A-F]{4}|[0-9A-F]{6}|[0-9A-F]{8})$`)

// HexColor reports whether or not v is a valid hexadecimal color code.
func HexColor(v string) bool {
	return rxHexColor.MatchString(v)
}

var rxBase32 = regexp.MustCompile(`^[A-Z2-7]+=*$`)

// Base32 reports whether or not v is a valid base32 string.
func Base32(v string) bool {
	return (len(v)%8 == 0) && rxBase32.MatchString(v)
}

// Email reports whether or not v is a valid email address.
//
// NOTE: Email uses net/mail.ParseAddress to determine the validity of v.
func Email(v string) bool {
	_, err := mail.ParseAddress(v)
	return err == nil
}

func URL(v string) bool {
	return false
}

func URI(v string) bool {
	return false
}

func PAN(v string) bool {
	return false
}

func CVV(v string) bool {
	return false
}

func SSN(v string) bool {
	return false
}

func EIN(v string) bool {
	return false
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

func CIDR(v string) bool {
	return false
}

func Phone(v string, cc ...string) bool {
	return false
}

func Zip(v string, cc ...string) bool {
	return false
}

func UUID(v string, ver ...int) bool {
	return false
}

// IP reports whether or not v is a valid IP address.
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

func RFC(v string, num int) bool {
	return false
}

func ISO(v string, num int) bool {
	return false
}

// Match reports whether or not the v contains any match of the regular expression re.
//
// NOTE: Match will panic if re has not been registered upfront with RegisterRegexp.
func Match(v string, re string) bool {
	return regexpCache.m[re].MatchString(v)
}
