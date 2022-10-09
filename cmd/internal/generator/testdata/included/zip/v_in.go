package testdata

type Validator struct {
	F1 string  `is:"zip"`
	F2 *string `is:"zip"`
	F3 string  `is:"zip:gb"`
	F4 *string `is:"zip:ru"`
}
