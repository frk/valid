package testdata

type RFCValidator struct {
	F1 string   `is:"rfc:1234"`
	F2 **string `is:"rfc:4321"`
	F3 **string `is:"required,rfc:6"`
}
