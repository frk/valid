package testdata

type T53Validator struct {
	F1 map[string]string      `is:"[]email"`
	F2 map[string]string      `is:"[email]"`
	F3 map[string]string      `is:"[email]phone"`
	F4 ***map[**string]***int `is:"[phone]rng:18:64"`
}
