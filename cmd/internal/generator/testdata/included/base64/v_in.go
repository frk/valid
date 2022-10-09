package testdata

type Validator struct {
	F1 string  `is:"base64"`
	F2 *string `is:"base64"`
	F3 string  `is:"base64:url"`
	F4 *string `is:"base64:true"`
}
