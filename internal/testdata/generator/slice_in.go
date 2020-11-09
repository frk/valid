package testdata

type SliceValidator struct {
	F1 []string       `is:"[]email"`
	F2 ***[]string    `is:"[]email"`
	F3 *[]*string     `is:"[]email,[]required"`
	F4 map[string]int `is:"[email]rng:18:64"`
}
