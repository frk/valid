package testdata

type Validator struct {
	F1 string  `is:"uuid"`
	F2 *string `is:"uuid"`
	F3 string  `is:"uuid:3"`
	F4 string  `is:"uuid:v4"`
	F5 string  `is:"uuid:v5"`
}
