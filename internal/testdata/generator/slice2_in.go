package testdata

type Slice2Validator struct {
	F1 []string `is:"required,[]email"`
	F2 []string `is:"len:1:,[]email"`
}
