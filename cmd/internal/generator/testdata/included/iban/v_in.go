package testdata

type Validator struct {
	F1 string  `is:"iban"`
	F2 *string `is:"iban"`
}
