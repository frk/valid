package cn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CN", A3: "CHN", Num: "156",
		Zip:   regexp.MustCompile(`^(?:0[1-7]|1[012356]|2[0-7]|3[0-6]|4[0-7]|5[1-7]|6[1-7]|7[1-5]|8[1345]|9[09])[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:(?:\+|00)86)?1(?:[3568][0-9]|4[579]|6[67]|7[01235678]|9[012356789])[0-9]{8}$`),
	})
}
