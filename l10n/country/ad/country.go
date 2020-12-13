package ad

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AD", A3: "AND", Num: "020",
		// NOTE(mkopriva): For the "AD5XX" post codes any digit between 0-9 is
		// allowed in the place of the Xs because I can't find anything substantial
		// to better handle the following: "PO Boxes in Andorra la Vella have separate
		// postcodes allocated to each group of 50 boxes - e.g., boxes 1001 to 1050
		// have a code of AD551, 1051 to 1100 a code of AD552 etc."
		// (from: https://en.wikipedia.org/wiki/Postal_codes_in_Andorra)
		Zip:   regexp.MustCompile(`^AD(?:[1-7]00|5[0-9]{2})$`),
		Phone: regexp.MustCompile(`^(?:\+376)?[346][0-9]{5}$`),
	})
}
