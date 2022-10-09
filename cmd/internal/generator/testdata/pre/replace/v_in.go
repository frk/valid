package testdata

type Validator struct {
	F1 string  `pre:"replace:a:b:3"`
	F2 *string `pre:"replace:foo:bar"`
}
