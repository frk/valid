package testdata

type Validator struct {
	F1 string  `is:"len:8"`
	F2 *string `is:"len:8:16"`
	F3 []byte  `is:"len:8:"`
	F4 *[]byte `is:"len::16"`
}
