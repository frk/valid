package testdata

type SuffixValidator struct {
	F1 string   `is:"suffix:foo"`
	F2 *string  `is:"suffix:bar"`
	F3 **string `is:"required,suffix:foo:bar:baz"`
}
