package testdata

type GreaterThanOrEqualValidator struct {
	F1 float64  `is:"gte:3.14"`
	F2 *int     `is:"gte:123"`
	F3 **uint32 `is:"required,gte:1"`
}
