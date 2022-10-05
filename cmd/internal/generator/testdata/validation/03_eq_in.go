package testdata

type T03Validator struct {
	F1 string   `is:"eq:foo"`
	F2 *float32 `is:"eq:3.14"`
}
