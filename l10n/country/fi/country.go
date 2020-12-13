package fi

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FI", A3: "FIN", Num: "246",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?358|0)[ ]?(?:4(?:0|1|2|4|5|6)?|50)[ ]?(?:[0-9][ ]?){4,8}[0-9]$`),
	})
}
