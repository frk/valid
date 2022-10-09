package testdata

type Validator struct {
	F1 float64  `pre:"ceil"`
	F2 *float64 `pre:"ceil"`
}
