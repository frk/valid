package testdata

type Validator struct {
	F1 float64 `is:"lt:3.14"`
	F2 *int    `is:"lt:314"`
	F3 *uint8  `is:"lt:31,required"`
}
