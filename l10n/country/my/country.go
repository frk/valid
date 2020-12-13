package my

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MY", A3: "MYS", Num: "458",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?6?01){1}(?:(?:[0145]{1}(?:-| )?[0-9]{7,8})|(?:[236789]{1}(?:-| )?[0-9]{7}))$`),
	})
}
