package testdata

type Validator struct {
	F1 string  `is:"jwt"`
	F2 *string `is:"jwt"`
}
