package testdata

////////////////////////////////////////////////////////////////////////////////
// ok
////////////////////////////////////////////////////////////////////////////////

func CustomErrorConstructor(k string, v any, r string, o ...any) error {
	return nil
}

type CustomErrorAggregator struct{}

func (CustomErrorAggregator) Error(k string, v any, r string, o ...any) {
	//...
}

func (CustomErrorAggregator) Out() error {
	//...
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// not-ok
////////////////////////////////////////////////////////////////////////////////

type NotAFuncObject struct{}

func NotANamedTypeObject() {}

func ErrorConstructorWithBadSignature(k string, v any, r string, o any) error {
	return nil
}

type ErrorAggregatorWithBadImpl struct{}

func (ErrorAggregatorWithBadImpl) Error(k string, v any, r CustomString, o ...any) {
	//...
}

func (ErrorAggregatorWithBadImpl) Out() error {
	//...
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type CustomString string
