package testdata

type Validator struct {
	F1 string  `is:"decimal"`
	F2 *string `is:"decimal"`
	F3 string  `is:"decimal:bg"`
	F4 *string `is:"decimal:hu"`
}
