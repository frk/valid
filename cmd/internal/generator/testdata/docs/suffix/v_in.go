package testdata

type Validator struct {
	F1 string  `is:"suffix:foo"`
	F2 *string `is:"suffix:bar"`
	F3 string  `is:"suffix:foo:bar:baz"`
}
