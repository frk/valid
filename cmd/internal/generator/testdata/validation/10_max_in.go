package testdata

type T10Validator struct {
	F1 float64 `is:"max:3.14"`
	F2 *uint   `is:"max:123"`
	F3 **int32 `is:"required,max:1"`
}
