package testdata

type Validator struct {
	F1 float64 `is:"gte:3.14"`
	F2 *int    `is:"gte:314"`
	F3 *uint8  `is:"gte:31,required"`
}
