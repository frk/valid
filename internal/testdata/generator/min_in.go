package testdata

type MinValidator struct {
	F1 float64  `is:"min:3.14"`
	F2 *int     `is:"min:123"`
	F3 **uint32 `is:"required,min:1"`
}
