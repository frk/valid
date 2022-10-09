package testdata

type Validator struct {
	F1 string  `pre:"rtrim:abc"`
	F2 *string `pre:"rtrim:xyz"`
}
