package testdata

type T35Validator struct {
	F1 string   `is:"ssn"`
	F2 **string `is:"ssn"`
	F3 **string `is:"required,ssn"`
}
