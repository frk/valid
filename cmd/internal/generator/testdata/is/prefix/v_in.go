package testdata

type Validator struct {
	F1 string  `is:"prefix:foo"`
	F2 *string `is:"prefix:bar"`
	F3 string  `is:"prefix:foo:bar:baz"`
}
