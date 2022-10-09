package testdata

type Validator struct {
	F1 string  `is:"iso639"`
	F2 *string `is:"iso639"`
	F3 string  `is:"iso639:1"`
	F4 *string `is:"iso639:2"`
}
