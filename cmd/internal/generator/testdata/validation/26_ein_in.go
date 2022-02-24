package testdata

type T26Validator struct {
	F1 string   `is:"ein"`
	F2 **string `is:"ein"`
	F3 **string `is:"required,ein"`
}
