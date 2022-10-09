package testdata

type Validator struct {
	F1 float64 `is:"lte:3.14"`
	F2 *int    `is:"lte:314"`
	F3 *uint8  `is:"lte:31,required"`
}
