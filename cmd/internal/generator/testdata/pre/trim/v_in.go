package testdata

type Validator struct {
	F1 string  `pre:"trim"`
	F2 *string `pre:"trim"`
	F3 string  `pre:"ltrim:abc,rtrim:xyz"`
	F4 struct {
		A []struct {
			X *string `pre:"trim,ltrim:abc,rtrim:xyz"`
		}
	}
}
