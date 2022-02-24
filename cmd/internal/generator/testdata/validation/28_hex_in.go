package testdata

type T28Validator struct {
	F1 string   `is:"hex"`
	F2 **string `is:"hex"`
	F3 **string `is:"required,hex"`
}
