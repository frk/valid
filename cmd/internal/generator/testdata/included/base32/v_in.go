package testdata

type Validator struct {
	F1 string  `is:"base32"`
	F2 *string `is:"base32"`
}
