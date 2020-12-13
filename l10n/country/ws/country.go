package ws

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "WS", A3: "WSM", Num: "882",
		Zip: regexp.MustCompile(`^WS[0-9]{4}$`),
	})
}
