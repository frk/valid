package testdata

type T44Validator struct {
	F1 string   `is:"re:foo"`
	F2 *string  `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
	F3 **string `is:"required,re:\"\\w+\""`
}
