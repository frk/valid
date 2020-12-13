package sl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SL", A3: "SLE", Num: "694",
		Phone: regexp.MustCompile(`^(?:0|94|\+94)?(?:7(?:0|1|2|5|6|7|8)(?: |-)?[0-9])[0-9]{6}$`),
	})
}
