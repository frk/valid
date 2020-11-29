package testdata

type UpperCaseValidator struct {
	F1 string   `is:"upper"`
	F2 **string `is:"upper"`
	F3 **string `is:"required,upper"`
}
