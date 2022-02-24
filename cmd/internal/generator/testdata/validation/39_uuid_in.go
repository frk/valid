package testdata

type T39Validator struct {
	F1 string   `is:"uuid"`
	F2 **string `is:"uuid:v5"`
	F3 **string `is:"required,uuid:3"`
}
