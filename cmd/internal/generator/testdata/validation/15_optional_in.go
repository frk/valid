package testdata

type T15Validator struct {
	F1 string  `is:"email,optional"`
	F2 string  `is:"email,len:5:85,optional"`
	F3 *string `is:"email,optional"`
	F4 *string `is:"email,len:5:85,optional"`
}
