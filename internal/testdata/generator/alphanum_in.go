package testdata

type AlphanumValidator struct {
	F1 string   `is:"alphanum"`
	F2 **string `is:"alphanum"`
	F3 **string `is:"required,alphanum"`
}
