package id

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {

	country.Add(country.Country{
		A2: "ID", A3: "IDN", Num: "360",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?62|0)8(?:1[123456789]|2[1238]|3[1238]|5[12356789]|7[78]|9[56789]|8[123456789])[ ?|0-9]{5,11}$`),
		// 15 digit number (ex. 02.271.824.1-413.000)
		VAT: regexp.MustCompile(`^(?:[0-9]{15})|(?:[0-9]{2}.[0-9]{3}.[0-9]{3}.[0-9][\-–][0-9]{3}.[0-9]{3})$`),
	})
}
