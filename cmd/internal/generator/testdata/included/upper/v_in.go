package testdata

type Validator struct {
	F1 string  `is:"upper"`
	F2 *string `is:"upper"`
}
