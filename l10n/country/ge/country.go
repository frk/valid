package ge

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GE", A3: "GEO", Num: "268",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?995)?(?:5|79)[0-9]{7}$`),
	})
}
