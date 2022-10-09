package testdata

type Validator struct {
	F1 string  `is:"iso31661a"`
	F2 *string `is:"iso31661a"`
	F3 string  `is:"iso31661a:2"`
	F4 *string `is:"iso31661a:3"`
}
