package testdata

type T34Validator struct {
	F1 string   `is:"phone"`
	F2 **string `is:"phone:us"`
	F3 **string `is:"required,phone:jp"`
}
