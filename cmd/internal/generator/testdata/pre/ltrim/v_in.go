package testdata

type Validator struct {
	F1 string  `pre:"ltrim:abc"`
	F2 *string `pre:"ltrim:xyz"`
}
