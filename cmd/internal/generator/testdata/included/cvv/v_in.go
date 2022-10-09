package testdata

type Validator struct {
	F1 string  `is:"cvv"`
	F2 *string `is:"cvv"`
}
