package in

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	var rxzip = regexp.MustCompile(`^[1-9][0-9]{2}[ ]?[0-9]{3}$`)
	var rxzipneg = regexp.MustCompile(`^(?:10|29|35|54|55|65|66|86|87|88|89)`)

	country.Add(country.Country{
		A2: "IN", A3: "IND", Num: "356",
		// References:
		// - https://en.wikipedia.org/wiki/Postal_Index_Number
		// - https://en.youbianku.com/India
		Zip: country.MatchStringFunc(func(v string) bool {
			return rxzip.MatchString(v) && !rxzipneg.MatchString(v)
		}),
		Phone: regexp.MustCompile(`^(?:\+?91|0)?[6789][0-9]{9}$`),
	})
}
