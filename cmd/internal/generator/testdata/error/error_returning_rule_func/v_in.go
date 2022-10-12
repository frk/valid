package testdata

type Validator struct {
	F1 string  `is:"rerr1"`
	F2 *string `is:"rerr1"`

	F3 string  `is:"rerr1,required"`
	F4 *string `is:"rerr1,required"`

	F5 string  `is:"rerr1,email,rerr2::13,required"`
	F6 *string `is:"email,rerr1,rerr2,required"`

	F7 []string  `is:"[]rerr2"`
	F8 []*string `is:"[]rerr2,rerr1,required"`
}
