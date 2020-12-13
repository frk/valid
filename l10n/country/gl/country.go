package gl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GL", A3: "GRL", Num: "304",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?299)?(?:[ ]?[0-9]{2}){3}$`),
	})
}
