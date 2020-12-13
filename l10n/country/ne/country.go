package ne

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NE", A3: "NER", Num: "562",
		Zip: country.RxZip4Digits,
	})
}
