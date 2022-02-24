package testdata

type T48Validator struct {
	// base field with rules & notnil
	F11a *string     `is:"notnil"`
	F11b **string    `is:"notnil"`
	F11c *****string `is:"notnil"`
	F12a **string    `is:"notnil,email"`
	F13a *string     `is:"notnil,hex,len:8:128"`
	F13b ***string   `is:"notnil,hex,len:8:128"`
	F14  *****string `is:"notnil,prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
}
