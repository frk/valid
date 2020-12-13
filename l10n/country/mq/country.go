package mq

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MQ", A3: "MTQ", Num: "474",
		Zip:   regexp.MustCompile(`^972(?:[0-8][0-9]|90)$`),
		Phone: regexp.MustCompile(`^(?:\+?596|0|00596)[67][0-9]{8}$`),
	})
}
