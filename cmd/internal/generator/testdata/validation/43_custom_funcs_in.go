package testdata

type T43Validator struct {
	F1 string   `is:"myrule"`
	F2 *string  `is:"myrule2"`
	F3 **string `is:"myrule2:foo:bar:baz"`
	F4 int64    `is:"myrule3:123:32.54:foo:true"`
}
