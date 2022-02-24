package testdata

type T47Validator struct {
	// base field with rules & nil guard
	F7a *string     `is:"email"`
	F7b **string    `is:"email"`
	F8a *string     `is:"hex,len:8:128"`
	F8b **string    `is:"hex,len:8:128"`
	F10 *****string `is:"prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
}
