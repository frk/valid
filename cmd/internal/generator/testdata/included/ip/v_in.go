package testdata

type Validator struct {
	F1 string  `is:"ip"`
	F2 *string `is:"ip"`
	F3 string  `is:"ip:4"`
	F4 *string `is:"ip:6"`
}
