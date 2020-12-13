package bd

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BD", A3: "BGD", Num: "050",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?880|0)1[13456789][0-9]{8}$`),
	})
}
