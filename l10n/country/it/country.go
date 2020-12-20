package it

import (
	"regexp"

	"github.com/frk/isvalid/internal/algo"
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 11 digits (the first 7 digits is a progressive number, the following 3
	// means the province of residence, the last digit is a check number
	// - The check digit is calculated using Luhn's Algorithm.)
	rxvat := regexp.MustCompile(`^IT[0-9]{11}$`)

	country.Add(country.Country{
		A2: "IT", A3: "ITA", Num: "380",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?39)?[ ]?3[0-9]{2}[ ]?[0-9]{6,7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			return algo.Luhn(v[2:])
		}),
	})
}
