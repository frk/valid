package testdata

type PANValidator struct {
	F1 string   `is:"pan"`
	F2 **string `is:"pan"`
	F3 **string `is:"required,pan"`
}
