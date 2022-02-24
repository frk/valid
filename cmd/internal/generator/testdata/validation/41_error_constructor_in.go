package testdata

type errorConstructor struct{}

func (errorConstructor) Error(key string, val interface{}, rule string, args ...interface{}) error {
	// ...
	return nil
}

type T41Validator struct {
	F1 string `is:"required,eq:foo"`
	ec errorConstructor
}
