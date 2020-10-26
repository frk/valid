package testdata

type HexColorValidator struct {
	F1 string   `is:"hexcolor"`
	F2 **string `is:"hexcolor"`
	F3 **string `is:"required,hexcolor"`
}
