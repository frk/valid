package pa

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PA", A3: "PAN", Num: "591",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?507)[0-9]{7,8}$`),
	})
}
