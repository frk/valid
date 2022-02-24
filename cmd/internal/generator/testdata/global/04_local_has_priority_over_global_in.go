package testdata

type T04Validator struct {
	F1 string   `is:"required"`
	F2 []string `is:"required,len::9"`
	ec errorConstructor
}

type errorConstructor struct{}

func (errorConstructor) Error(key string, val any, rule string, args ...any) error {
	// ...
	return nil
}
