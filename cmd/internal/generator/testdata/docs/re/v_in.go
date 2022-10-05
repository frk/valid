package testdata

type Validator struct {
	F1 string  `is:"re:foo"`
	F2 string  `is:"re:\"\\w+\""`
	F3 *string `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
}
