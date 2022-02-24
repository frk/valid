package testdata

type T04Validator struct {
	F1 string    `is:"required,email" pre:"trim"`
	F2 string    `is:"required,len::256" pre:"repeat:2"`
	F3 *string   `is:"required,alnum" pre:"trim,title,quote"`
	F4 ***string `is:"required,upper" pre:"upper"`
	F5 **float64 `is:"required" pre:"round"`
}
