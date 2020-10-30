package testdata

type ReferencesValidator struct {
	Min, Max int

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
}

type References2Validator struct {
	Min, Max int
	ec       errorConstructor

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
}

type References3Validator struct {
	Min, Max int
	ea       errorAggregator

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
}
