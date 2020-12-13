package gb

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GB", A3: "GBR", Num: "826",
		Zip:   regexp.MustCompile(`^(?i)(?:gir\s?0aa|[a-z]{1,2}[0-9][0-9a-z]?\s?(?:[0-9][a-z]{2})?)$`),
		Phone: regexp.MustCompile(`^(?:\+?44|0)7[0-9]{9}$`),
	})
}
