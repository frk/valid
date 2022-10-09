package testdata

type Validator struct {
	F1 string  `is:"cidr"`
	F2 *string `is:"cidr"`
}
