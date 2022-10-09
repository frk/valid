package testdata

type Validator struct {
	F1 string  `pre:"trimsuffix:foo"`
	F2 *string `pre:"trimsuffix:bar"`
}
