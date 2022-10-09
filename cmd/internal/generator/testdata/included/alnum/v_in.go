package testdata

type Validator struct {
	F1 string  `is:"alnum"`
	F2 *string `is:"alnum:bg"`
}
