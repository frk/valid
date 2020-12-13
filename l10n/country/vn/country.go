package vn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VN", A3: "VNM", Num: "704",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?84|0)(?:3[2-9]|5[2689]|7[0|6-9]|8[1-6|89]|9[0-9])(?:[0-9]{7})$`),
	})
}
