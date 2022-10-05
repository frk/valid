package testdata

type Validator struct {
	F1 float64 `is:"max:3.14"`
	F2 *int    `is:"max:314"`
	F3 *uint8  `is:"max:31,required"`
}
