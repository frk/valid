package testdata

type Validator struct {
	F1 string `is:"in:foo:bar:baz"`
	F2 *int   `is:"in:10:20:30:40"`
}
