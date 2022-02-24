package testdata

type T31Validator struct {
	F1 string   `is:"mac"`
	F2 **string `is:"mac:6"`
	F3 **string `is:"required,mac:8"`
}
