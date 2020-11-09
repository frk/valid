package mypkg

// rules

func MyRule(v string) bool {
	// ...
	return false
}

func MyRule2(v ...string) bool {
	// ...
	return false
}

func MyRule3(v int64, i int, f float64, s string, b bool) bool {
	// ...
	return false
}

// error handlers

type MyErrorConstructor struct{}

func (MyErrorConstructor) Error(key string, val interface{}, rule string, args ...interface{}) error {
	// ...
	return nil
}

type MyErrorAggregator struct{}

func (MyErrorAggregator) Error(key string, val interface{}, rule string, args ...interface{}) {
	// ...
}

func (MyErrorAggregator) Out() error {
	// ...
	return nil
}

// "isValider"

type MyString string

func (MyString) IsValid() bool {
	// ...
	return false
}

type MyInt int

func (*MyInt) IsValid() bool {
	// ...
	return false
}

// enum
type MyKind uint

const (
	_ MyKind = 1 << iota
	MyFoo
	MyBar
	_
	MyBaz
)
