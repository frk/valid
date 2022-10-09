package testdata

type Validator struct {
	F1 float64  `pre:"round"`
	F2 *float64 `pre:"round"`
}
