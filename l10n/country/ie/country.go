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
		// 'IE'+7 digits and one letter, optionally followed by a 'W' for married women, e.g. IE1234567T or IE1234567TW
		// or 'IE'+7 digits and two letters, e.g. IE1234567FA (since January 2013)
		// or 'IE'+one digit, one letter/"+"/"*", 5 digits and one letter (old style, currently being phased out)
		VAT: regexp.MustCompile(`^IE(?:[0-9]{7}[A-Z]{1,2}|[0-9][A-Z+*][0-9]{5}[A-Z])$`),
	})
}
