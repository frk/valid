package testdata

type T06Validator struct {
	F1 map[string]string    `pre:"[]trim"`
	F2 map[string]string    `is:"[]email" pre:"[]trim"`
	F3 map[string][]*string `pre:"[][]trim"`
	F4 []map[string]string  `is:"[]len:4:,[]email" pre:"[][]trim"`
}
