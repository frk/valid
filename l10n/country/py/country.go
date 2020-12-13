package py

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PY", A3: "PRY", Num: "600",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?595|0)9[9876][0-9]{7}$`),
	})
}
