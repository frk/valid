package testdata

type Validator struct {
	F1 *string   `is:"email,noguard"`
	F2 **[]int64 `is:"len:5,noguard"`
}
