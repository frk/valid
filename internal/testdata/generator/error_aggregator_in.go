package testdata

type errorAggregator struct{}

func (errorAggregator) Error(key string, val interface{}, rule string, args ...interface{}) {
	// ...
}

func (errorAggregator) Out() error {
	// ...
	return nil
}

type ErrorAggregatorValidator struct {
	F1 string `is:"required,eq:foo"`
	ea errorAggregator
}
