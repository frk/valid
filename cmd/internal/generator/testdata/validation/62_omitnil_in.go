package testdata

type T62Validator struct {
	F1 *string `is:"email,omitnil"`
	F2 *int    `is:"min:8,max:128,omitnil"`
}
