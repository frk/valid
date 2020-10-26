package testdata

type HexValidator struct {
	F1 string   `is:"hex"`
	F2 **string `is:"hex"`
	F3 **string `is:"required,hex"`
}
