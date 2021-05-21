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

func MyBadRule1() bool {
	// ...
	return false
}

func MyBadRule2(v string) int {
	// ...
	return 0
}

func MyBadRule3(v int64, i int, f float64, s string, b bool) (bool, error) {
	// ...
	return false, nil
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

// enum without a (useful) const value
type MyNoConst uint

const (
	_ MyNoConst = 1 << iota
	_
	_
	myNoConst
)

type MyNonBasicType struct {
	// ...
}

type MyNonBasicType2 []string

// enum + isValider
type MyEnumValider string

const (
	MyEnumValiderFoo MyEnumValider = "foo"
	MyEnumValiderBar MyEnumValider = "bar"
	MyEnumValiderBaz MyEnumValider = "baz"
)

func (MyEnumValider) IsValid() bool {
	// ...
	return false
}
