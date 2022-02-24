package testdata

type T01Validator struct {
	F1 string   `is:"required"`
	F2 []string `is:"required,len::9"`
}
