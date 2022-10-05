package testdata

type Validator struct {
	F1 float64 `is:"gt:3.14"`
	F2 *int    `is:"gt:314"`
	F3 *uint8  `is:"gt:31,required"`
}
