package testdata

import (
	"github.com/frk/valid"
)

type Validator struct {
	F1 string  `is:"ccy"`
	F2 *string `is:"ccy"`
	F3 string  `is:"ccy:gbp"`
	F4 *string `is:"ccy:eur:&ccyOpts"`
	F5 string  `is:"ccy::&ccyOpts"`

	ccyOpts *valid.CurrencyOpts
}
