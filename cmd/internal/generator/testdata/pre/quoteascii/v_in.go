package testdata

type Validator struct {
	F1 string  `pre:"quoteascii"`
	F2 *string `pre:"quoteascii"`
}
