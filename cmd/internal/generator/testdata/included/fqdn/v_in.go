package testdata

type Validator struct {
	F1 string  `is:"fqdn"`
	F2 *string `is:"fqdn"`
}
