package testdata

type ISOValidator struct {
	F1 string   `is:"iso:1234"`
	F2 **string `is:"iso:4321"`
	F3 **string `is:"required,iso:6"`
}
