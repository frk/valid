package testdata

type BaseFieldsWithRulesAndRequiredValidator struct {
	// base field with rules & required
	F15a *string     `is:"required"`
	F15b **string    `is:"required"`
	F15c *****string `is:"required"`
	F16a **string    `is:"required,email"`
	F17a *string     `is:"required,hex,len:8:128"`
	F17b ***string   `is:"required,hex,len:8:128"`
	F18  *****string `is:"required,prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
}
