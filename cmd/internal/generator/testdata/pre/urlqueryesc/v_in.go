package testdata

type Validator struct {
	F1 string  `pre:"urlqueryesc"`
	F2 *string `pre:"urlqueryesc"`
}
