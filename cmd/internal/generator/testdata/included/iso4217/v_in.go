package testdata

type Validator struct {
	F1 string  `is:"iso4217"`
	F2 *string `is:"iso4217"`
}
