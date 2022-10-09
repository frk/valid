package testdata

type Validator struct {
	F1 float64 `is:"rng:3.14:42"`
	F2 *int    `is:"rng:-8:256"`
	F3 *uint8  `is:"rng:1:2,required"`
}
