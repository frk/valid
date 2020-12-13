package mt

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MT", A3: "MLT", Num: "470",
		Zip:   regexp.MustCompile(`^(?i)[a-z]{3}\s{0,1}[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+?356|0)?(?:99|79|77|21|27|22|25)[0-9]{6}$`),
	})
}
