package testdata

type Validator struct {
	F1 string  `is:"digits"`
	F2 *string `is:"digits"`
}
