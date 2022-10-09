package testdata

type Validator struct {
	F1 string  `is:"isbn"`
	F2 *string `is:"isbn"`
	F3 string  `is:"isbn:10"`
	F4 *string `is:"isbn:13"`
}
