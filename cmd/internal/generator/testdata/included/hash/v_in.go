package testdata

type Validator struct {
	F1 string  `is:"hash:md5"`
	F2 *string `is:"hash:sha512"`
}
