package testdata

type Validator struct {
	F1 string  `is:"phone"`
	F2 *string `is:"phone"`
	F3 string  `is:"phone:gb"`
	F4 *string `is:"phone:jp"`
}
