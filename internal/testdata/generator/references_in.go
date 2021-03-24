package testdata

type ReferencesValidator struct {
	Min, Max  int
	SomeValue string

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
	F3 string `is:"phone:&SomeValue"`
	F4 string `is:"contains:&SomeValue:bar:baz"`
}

type References2Validator struct {
	Min, Max  int
	SomeValue string
	ec        errorConstructor

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
	F3 string `is:"phone:&SomeValue"`
	F4 string `is:"contains:foo:bar:&SomeValue"`
}

type References3Validator struct {
	Min, Max  int
	SomeValue string
	ea        errorAggregator

	F1 string `is:"len::&Max"`
	F2 *int   `is:"rng:&Min:&Max"`
	F3 string `is:"phone:&SomeValue"`
	F4 string `is:"contains:foo:&SomeValue:baz"`
}
