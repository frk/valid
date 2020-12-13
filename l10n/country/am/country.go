package am

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AM", A3: "ARM", Num: "051", Zip: country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?374|0)(?:(?:10|[9|7][0-9])[0-9]{6}|[2-4][0-9]{7})$`),
	})
}
