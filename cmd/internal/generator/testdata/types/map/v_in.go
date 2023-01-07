package testdata

type Validator struct {
	F1 map[string]string      `is:"[]email"`
	F2 map[string]string      `is:"[email]"`
	F3 map[string]string      `is:"[email]phone"`
	F4 ***map[**string]***int `is:"[phone]rng:18:64"`

	F5 []map[*map[string]string][]int `is:"[][[email]phone:ca]len::10,[]rng:-54:256"`
}
