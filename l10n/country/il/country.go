package il

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IL", A3: "ISR", Num: "376",
		Zip:   regexp.MustCompile(`^(?:[0-9]{5}|[0-9]{7})$`),
		Phone: regexp.MustCompile(`^(?:\+972|0)(?:[23489]|5[012345689]|77)[1-9][0-9]{6}$`),
	})
}
