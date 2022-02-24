package testdata

type T24Validator struct {
	F1 string   `is:"cidr"`
	F2 **string `is:"cidr"`
	F3 **string `is:"required,cidr"`
}
