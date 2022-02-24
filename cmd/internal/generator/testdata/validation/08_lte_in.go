package testdata

type T08Validator struct {
	F1 float64 `is:"lte:3.14"`
	F2 *uint   `is:"lte:123"`
	F3 **int32 `is:"required,lte:1"`
}
