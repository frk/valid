package np

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NP", A3: "NPL", Num: "524",
		Zip:   regexp.MustCompile(`^(?:10|21|22|32|33|34|44|45|56|57)[0-9]{3}$|^(?:977)$`),
		Phone: regexp.MustCompile(`^(?:\+?977)?9[78][0-9]{8}$`),
	})
}
