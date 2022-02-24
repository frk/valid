package testdata

type T50Validator struct {
	// base field with rules
	F1 string `is:"email"`
	F2 string `is:"hex,len:8:128"`
	F3 string `is:"prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
}
