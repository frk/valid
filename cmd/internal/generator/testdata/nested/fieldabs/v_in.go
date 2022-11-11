package testdata

type Validator struct {
	F1 int
	F2 int `is:"gt:&F1"`

	S1 struct {
		S2 struct {
			F1 int
			F2 int `is:"gt:&S1.S2.F1"`
		}
		F2 int `is:"lt:&S1.S2.F1"`
	}
}
