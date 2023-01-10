package testdata

type Validator struct {
	F1 string    `is:"email" pre:"trim"`
	F2 string    `is:"len::256" pre:"repeat:2"`
	F3 *string   `is:"alnum" pre:"trim,title,quote"`
	F4 ***string `is:"upper" pre:"upper"`
	F5 **float64 `pre:"round"`
}
