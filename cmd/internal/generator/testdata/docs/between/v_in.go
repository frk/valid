package testdata

type Validator struct {
	F1 float64 `is:"between:3.14:42"`
	F2 *int    `is:"between:-8:256"`
	F3 *uint8  `is:"between:1:2,required"`
}
