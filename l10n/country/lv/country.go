package lv

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LV", A3: "LVA", Num: "428",
		Zip: regexp.MustCompile(`^LV\-[0-9]{4}$`),
	})
}
