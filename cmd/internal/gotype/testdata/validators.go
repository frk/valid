package testdata

type Test1Validator struct {
	// ...
}

type Test2Validator struct {
	err errorConstructor
}

type Test3Validator struct {
	err errorAggregator
}

type Test4Validator struct {
	context string
}

type Test5Validator struct {
	// ...
}

func (Test5Validator) beforevalidate() error { return nil }
func (*Test5Validator) AfterValidate() error { return nil }

type errorConstructor struct{}

func (errorConstructor) Error(key string, val interface{}, rule string, args ...interface{}) error {
	// ...
	return nil
}

type errorAggregator struct{}

func (errorAggregator) Error(key string, val interface{}, rule string, args ...interface{}) {
	// ...
}

func (errorAggregator) Out() error {
	// ...
	return nil
}
