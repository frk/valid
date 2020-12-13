package sz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SZ", A3: "SWZ", Num: "748",
		Zip: regexp.MustCompile(`^[HMSL][0-9]{3}$`),
	})
}
