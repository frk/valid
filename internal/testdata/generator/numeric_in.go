package testdata

type NumericValidator struct {
	F1 string   `is:"numeric"`
	F2 **string `is:"numeric"`
	F3 **string `is:"required,numeric"`
}
