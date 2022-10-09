package testdata

type Validator struct {
	F1 string  `is:"ssn"`
	F2 *string `is:"ssn"`
}
