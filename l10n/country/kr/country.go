package kr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KR", A3: "KOR", Num: "410",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+?82)[ \-]?)?0?1(?:[0|1|6|7|8|9]{1})[ \-]?[0-9]{3,4}[ \-]?[0-9]{4}$`),
	})
}
