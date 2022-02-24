package testdata

type T37Validator struct {
	F1 string   `is:"upper"`
	F2 **string `is:"upper"`
	F3 **string `is:"required,upper"`
}
