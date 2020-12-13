package ng

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NG", A3: "NGA", Num: "566",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?234|0)?[789][0-9]{9}$`),
	})
}
