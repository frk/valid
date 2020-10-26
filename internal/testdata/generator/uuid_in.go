package testdata

type UUIDValidator struct {
	F1 string   `is:"uuid"`
	F2 **string `is:"uuid:v4"`
	F3 **string `is:"required,uuid:3:v1:v5"`
}
