package testdata

type Validator struct {
	F1 string `is:"email,optional"`
	F2 *int64 `is:"eq:42,optional"`
}
