package ly

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LY", A3: "LBY", Num: "434",
		Phone: regexp.MustCompile(`^(?:(?:\+?218)|0)?(?:9[1-6][0-9]{7}|[1-8][0-9]{7,9})$`),
	})
}
