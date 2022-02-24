package mypkg

func NewError(k string, v any, r string, o ...any) error {
	return nil
}

type ErrorList struct {
	// ...
}

func (*ErrorList) Error(k string, v any, r string, o ...any) {
	// ...
}

func (*ErrorList) Out() error {
	return nil
}
