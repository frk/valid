package testdata

type Validator struct {
	F1 string  `is:"iprange"`
	F2 *string `is:"iprange"`
}
