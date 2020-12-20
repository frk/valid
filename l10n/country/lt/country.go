package lt

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LT", A3: "LTU", Num: "440",
		Zip:   regexp.MustCompile(`^LT\-[0-9]{5}$`),
		Phone: regexp.MustCompile(`^(?:\+370|8)[0-9]{8}$`),
		VAT:   regexp.MustCompile(`^LT[0-9]{9}(?:[0-9]{3})?$`),
	})
}
