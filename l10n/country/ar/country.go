package ar

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AR", A3: "ARG", Num: "032",
		Zip:   regexp.MustCompile(`^(?:[0-9]{4})|(?:[A-Z][0-9]{4}[A-Z]{3})$`),
		Phone: regexp.MustCompile(`^\+?549(?:11|[2368][0-9])[0-9]{8}$`),
		VAT:   regexp.MustCompile(`^[0-9]{11}$`),
	})
}
