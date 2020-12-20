package al

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AL", A3: "ALB", Num: "008", Zip: country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+355|0)6[789][0-9]{6}$`),
		// 10 characters, the first position following the prefix
		// is "J" or "K" or "L", and the last character is a letter
		// – e.g. K99999999L or L99999999G
		VAT: regexp.MustCompile(`^[JKL][0-9]{8}[A-Z]$`),
	})
}
