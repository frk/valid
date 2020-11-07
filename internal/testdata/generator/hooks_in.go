package testdata

type HooksValidator struct {
	F1 string `is:"required"`
}

func (v *HooksValidator) beforevalidate() error {
	return nil
}

func (v HooksValidator) AfterValidate() error {
	return nil
}

type Hooks2Validator struct {
	F1 string `is:"required"`
	ea errorAggregator
}

func (v *Hooks2Validator) beforevalidate() error {
	return nil
}

func (v Hooks2Validator) AfterValidate() error {
	return nil
}
