package testdata

type T52Validator struct {
	F1 []string    `is:"[]email"`
	F2 ***[]string `is:"[]email"`
	F3 *[]*string  `is:"[]email,required"`
}
