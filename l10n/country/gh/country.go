package gh

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GH", A3: "GHA", Num: "288",
		Zip:   regexp.MustCompile(`^[A-Z][A-Z0-9]-[0-9]{4}-[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+233|0)(?:20|50|24|54|27|57|26|56|23|28)[0-9]{7}$`),
	})
}
