package om

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "OM", A3: "OMN", Num: "512",
		Zip:   country.RxZip3Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+|00)968)?(?:9[1-9])[0-9]{6}$`),
	})
}
