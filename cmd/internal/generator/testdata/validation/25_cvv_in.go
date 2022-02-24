package testdata

type T25Validator struct {
	F1 string   `is:"cvv"`
	F2 **string `is:"cvv"`
	F3 **string `is:"required,cvv"`
}
