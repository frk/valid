package testdata

type PrefixValidator struct {
	F1 string   `is:"prefix:foo"`
	F2 *string  `is:"prefix:bar"`
	F3 **string `is:"required,prefix:foo:bar:baz"`
}
