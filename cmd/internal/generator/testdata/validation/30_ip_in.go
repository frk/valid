package testdata

type T30Validator struct {
	F1 string   `is:"ip"`
	F2 **string `is:"ip:v4"`
	F3 **string `is:"required,ip:6"`
}
