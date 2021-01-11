package ni

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NI", A3: "NIC", Num: "558",
		Zip: country.RxZip5Digits,
		// 3 digits, 1 dash, 6 digits, 1 dash, 4 digits followed by 1 letter
		VAT: regexp.MustCompile(`^[0-9]{3}-[0-9]{6}-[0-9]{4}[A-Z]$`),
	})
}
