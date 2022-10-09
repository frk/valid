package testdata

type Validator struct {
	F1 string  `is:"ascii"`
	F2 *string `is:"ascii"`
}
