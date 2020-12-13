package si

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SI", A3: "SVN", Num: "705",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+386[ ]?|0)(?:(?:[0-9]{1}[ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2})|(?:[0-9]{2}(?:[ ]?[0-9]{3}){2}))$`),
	})
}
