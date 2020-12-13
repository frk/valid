package ie

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IE", A3: "IRL", Num: "372",
		// References:
		// - https://stackoverflow.com/questions/33391412/validation-for-irish-eircode
		// - https://www.eircode.ie/docs/default-source/Common/prepareyourbusinessforeircode-edition3published.pdf
		// - https://en.wikipedia.org/wiki/Postal_addresses_in_the_Republic_of_Ireland
		Zip:   regexp.MustCompile(`^(?:[AC-FHKNPRTV-Y][0-9]{2}|D6W)[ -]?[0-9AC-FHKNPRTV-Y]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+?353|0)8[356789][0-9]{7}$`),
	})
}
