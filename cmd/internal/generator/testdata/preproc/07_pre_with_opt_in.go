package testdata

type T07Validator struct {
	F1 string  `pre:"pre_with_opt"`
	F2 *string `pre:"trim,pre_with_opt"`
	F3 string  `pre:"pre_with_opt,trim"`
	F4 *string `pre:"pre_with_opt,trim"`

	G1 string  `pre:"pre_with_opt2"`
	G2 *string `pre:"trim,pre_with_opt2"`
	G3 string  `pre:"pre_with_opt2,trim"`
	G4 *string `pre:"pre_with_opt2,trim"`
}
