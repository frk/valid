package testdata

type Validator struct {
	F1 string  `pre:"trimprefix:foo"`
	F2 *string `pre:"trimprefix:bar"`
}
