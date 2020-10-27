package testdata

type GreaterThanValidator struct {
	F1 float64  `is:"gt:3.14"`
	F2 *int     `is:"gt:123"`
	F3 **uint32 `is:"required,gt:1"`
}
