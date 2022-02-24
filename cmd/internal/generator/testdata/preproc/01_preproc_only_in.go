package testdata

type T01Validator struct {
	F1 string    `pre:"trim"`
	F2 string    `pre:"repeat:2"`
	F3 *string   `pre:"trim,repeat:1"`
	F4 string    `pre:"replace:a:b:3"`
	F5 string    `pre:"replace:foo:bar"`
	F6 string    `pre:"trim,title,quote"`
	F7 **float64 `pre:"round"`
}
