package testdata

type T45aValidator struct {
	F1 string `is:"required"`
}

func (v *T45aValidator) beforevalidate() error {
	return nil
}

func (v T45aValidator) AfterValidate() error {
	return nil
}

type T45bValidator struct {
	F1 string `is:"required"`
	ea errorAggregator
}

func (v *T45bValidator) beforevalidate() error {
	return nil
}

func (v T45bValidator) AfterValidate() error {
	return nil
}
