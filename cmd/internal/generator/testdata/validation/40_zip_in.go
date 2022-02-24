package testdata

type T40Validator struct {
	F1 string   `is:"zip"`
	F2 **string `is:"zip:us"`
	F3 **string `is:"required,zip:jp"`
}
