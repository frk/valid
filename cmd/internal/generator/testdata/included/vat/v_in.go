package testdata

type Validator struct {
	F1 string  `is:"vat"`
	F2 *string `is:"vat"`
	F3 string  `is:"vat:gb"`
	F4 *string `is:"vat:ru"`
}
