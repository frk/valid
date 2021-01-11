package ve

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {

	country.Add(country.Country{
		A2: "VE", A3: "VEN", Num: "862",
		Zip: regexp.MustCompile(`^[0-9]{4}(?:-[A-Z])?$`),
		// First digit must be (J, G, V, E), one dash (-), next 9 (nine) numbers
		// like J-305959918, in some cases can be written like J-30595991-8
		VAT: regexp.MustCompile(`^[JGVE]-[0-9]{8}-?[0-9]$`),
	})
}
