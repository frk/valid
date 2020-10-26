package testdata

type URIValidator struct {
	F1 string   `is:"uri"`
	F2 **string `is:"uri"`
	F3 **string `is:"required,uri"`
}
