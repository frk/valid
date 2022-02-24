package testdata

type T12Validator struct {
	F1 string   `is:"len:10"`
	F2 *string  `is:"len:8:256"`
	F3 **[]byte `is:"required,len:1:2"`
	F4 []int    `is:"len:4:"`
	F5 []string `is:"len::15"`
}
