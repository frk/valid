package testdata

type ReferencesValidator struct {
	Min, Max int

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
}
