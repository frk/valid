package testdata

type Validator struct {
	F1 *string `is:"email,omitnil"`
	F2 []int64 `is:"len:5,omitnil"`
}
