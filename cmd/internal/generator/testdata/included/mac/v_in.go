package testdata

type Validator struct {
	F1 string  `is:"mac"`
	F2 *string `is:"mac"`
	F3 string  `is:"mac:6"`
	F4 *string `is:"mac:8"`
}
