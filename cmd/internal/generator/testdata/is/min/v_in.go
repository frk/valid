package testdata

type Validator struct {
	F1 float64 `is:"min:3.14"`
	F2 *int    `is:"min:314"`
	F3 *uint8  `is:"min:31,required"`
}
