package eg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "EG", A3: "EGY", Num: "818",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+?20)|0)?1[0125][0-9]{8}$`),
	})
}
