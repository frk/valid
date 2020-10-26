package testdata

type URLValidator struct {
	F1 string   `is:"url"`
	F2 **string `is:"url"`
	F3 **string `is:"required,url"`
}
