package testdata

type Validator struct {
	F1 string  `is:"lower"`
	F2 *string `is:"lower"`
}
