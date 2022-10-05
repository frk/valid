package testdata

type Validator struct {
	F1 string `is:"eq:foo:bar:baz"`
	F2 *int   `is:"eq:64:128"`
	F3 any    `is:"eq:foo:0.8:true"`
}
