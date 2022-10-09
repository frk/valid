package testdata

type Validator struct {
	F1 string  `is:"alpha"`
	F2 *string `is:"alpha:ja"`
}
