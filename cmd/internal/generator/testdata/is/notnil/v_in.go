package testdata

type Validator struct {
	F1 *string           `is:"notnil"`
	F2 []string          `is:"notnil"`
	F3 map[string]string `is:"notnil"`
	F4 interface{}       `is:"notnil"`

	F5 func()      `is:"notnil"`
	F6 chan string `is:"notnil"`
	F7 **[]string  `is:"notnil"`
	F8 ***string   `is:"notnil"`
}
