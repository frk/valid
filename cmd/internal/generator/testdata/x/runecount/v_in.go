package testdata

type Validator struct {
	F1 string  `is:"runecount:8"`
	F2 *string `is:"runecount:8:16"`
	F3 []byte  `is:"runecount:8:"`
	F4 *[]byte `is:"runecount::16"`
}
