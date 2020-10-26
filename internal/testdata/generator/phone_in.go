package testdata

type PhoneValidator struct {
	F1 string   `is:"phone"`
	F2 **string `is:"phone:us"`
	F3 **string `is:"required,phone:us:ca:jp"`
}
