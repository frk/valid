package testdata

type Validator struct {
	F1 string  `is:"email"`
	F2 *string `is:"email"`
}
