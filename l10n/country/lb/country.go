package lb

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LB", A3: "LBN", Num: "422",
		Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:[0-9]{4} [0-9]{4})$`),
		Phone: regexp.MustCompile(`^(?:\+?961)?(?:(?:3|81)[0-9]{6}|7[0-9]{7})$`),
	})
}
