package testdata

type T27Validator struct {
	F1 string  `is:"fqdn"`
	F2 string  `is:"fqdn"`
	F3 *string `is:"fqdn"`
	F4 *string `is:"fqdn"`
	F5 *string `is:"fqdn"`
}
