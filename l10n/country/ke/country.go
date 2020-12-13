package ke

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KE", A3: "KEN", Num: "404",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?254|0)(?:7|1)[0-9]{8}$`),
	})
}
