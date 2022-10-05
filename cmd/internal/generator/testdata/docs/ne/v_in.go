package testdata

type Validator struct {
	F1 string `is:"ne:foo:bar"`
	F2 *int   `is:"ne:64:128"`
	F3 any    `is:"ne:foo:0.8:true"`
}
