package testdata

type Validator struct {
	F1 []string  `is:"len::12"`
	F2 []string  `is:"[]email"`
	F3 []string  `is:"[]alnum,len:11:256"`
	F4 []*string `is:"[]notnil,alnum"`

	F5 []string `is:"required,[]email"`
	F6 []string `is:"len:7:,[]email"`

	F7 ***[]string `is:"[]email"`
	F8 *[]*string  `is:"[]email,required"`
}
