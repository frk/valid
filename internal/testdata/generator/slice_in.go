package testdata

type SliceValidator struct {
	F1 []string       `is:"[]email"`
	F2 ***[]string    `is:"[]email"`
	F3 *[]*string     `is:"[]email,required"`
	F4 map[string]int `is:"[email]rng:18:64"`

	F5 []map[*map[string]string][]int `is:"[][[email]phone:us:ca]len::10,[]rng:-54:256"`
}
