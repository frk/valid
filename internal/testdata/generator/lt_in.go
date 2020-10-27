package testdata

type LessThanValidator struct {
	F1 float64 `is:"lt:3.14"`
	F2 *uint   `is:"lt:123"`
	F3 **int32 `is:"required,lt:1"`
}
