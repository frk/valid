package testdata

type T11Validator struct {
	F1 float64 `is:"rng:3.14:42"`
	F2 *uint   `is:"rng:8:256"`
	F3 **int32 `is:"required,rng:1:2"`
}
