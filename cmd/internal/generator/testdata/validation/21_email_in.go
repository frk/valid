package testdata

type T21Validator struct {
	F1 string   `is:"email"`
	F2 **string `is:"email"`
	F3 **string `is:"required,email"`
}
