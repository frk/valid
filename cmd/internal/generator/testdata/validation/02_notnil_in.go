package testdata

type T02Validator struct {
	F1 *string           `is:"notnil"`
	F2 []string          `is:"notnil"`
	F3 map[string]string `is:"notnil"`
	F4 interface{}       `is:"notnil"`
}
