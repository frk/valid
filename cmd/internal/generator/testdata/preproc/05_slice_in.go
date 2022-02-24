package testdata

type T05Validator struct {
	F1 []string      `pre:"[]trim"`
	F2 []string      `is:"[]email" pre:"[]trim"`
	F3 [][]*string   `pre:"[][]trim"`
	F4 [][][]*string `is:"len::10,[]len:3:5,[]len:4:,[]email" pre:"[][][]trim"`
}
