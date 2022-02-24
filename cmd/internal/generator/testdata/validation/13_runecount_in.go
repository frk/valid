package testdata

type T13Validator struct {
	F1 string   `is:"runecount:10"`
	F2 *string  `is:"runecount:8:256"`
	F3 **[]byte `is:"required,runecount:1:2"`
	F4 []byte   `is:"runecount:4:"`
	F5 string   `is:"runecount::15"`
}
