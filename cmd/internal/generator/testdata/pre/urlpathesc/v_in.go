package testdata

type Validator struct {
	F1 string  `pre:"urlpathesc"`
	F2 *string `pre:"urlpathesc"`
}
