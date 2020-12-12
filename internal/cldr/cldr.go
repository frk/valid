package cldr

import (
	"strings"
)

type LocaleInfo struct {
	// The locale's language tag.
	Lang string
	// The decimal separator of the locale.
	SepDecimal rune
	// The group (thousands) separator of the locale.
	SepGroup rune
	// The rune for the locale's digit 0.
	DigitZero rune
	// The rune for the locale's digit 9.
	DigitNine rune
}

func Locale(loc string) (LocaleInfo, bool) {
	li, ok := localemap[loc]
	if !ok {
		for {
			if i := strings.LastIndexByte(loc, '_'); i < 0 {
				break
			} else {
				loc = loc[:i]
				if li, ok := localemap[loc]; ok {
					return li, ok
				}
			}
		}
	}
	return li, ok
}

var localemap map[string]LocaleInfo

func init() {
	localemap = make(map[string]LocaleInfo)
	for _, li := range localeslice {
		localemap[li.Lang] = li
	}
}
