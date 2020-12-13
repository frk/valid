package al

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AL", A3: "ALB", Num: "008", Zip: country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+355|0)6[789][0-9]{6}$`),
	})
}
