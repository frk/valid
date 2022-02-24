package testdata

type T01Validator struct {
	F1 string            `is:"required"`
	F2 []string          `is:"required"`
	F3 map[string]string `is:"required"`
	F4 int               `is:"required"`
	F5 uint              `is:"required"`
	F6 float64           `is:"required"`
	F7 bool              `is:"required"`
	F8 *struct {
		// ...
	} `is:"required"`
	F9 interface {
		// ...
	} `is:"required"`
}
