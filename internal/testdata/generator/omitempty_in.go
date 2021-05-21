package testdata

type OmitEmptyValidator struct {
	F1 string  `is:"email,omitempty"`
	F2 string  `is:"email,len:5:85,omitempty"`
	F3 *string `is:"email,omitempty"`
	F4 *string `is:"email,len:5:85,omitempty"`
}
