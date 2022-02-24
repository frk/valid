package mypkg

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

// enum + IsValid
type MyEnumIsValid string

const (
	MyEnumFoo MyEnumIsValid = "foo"
	MyEnumBar MyEnumIsValid = "bar"
	MyEnumBaz MyEnumIsValid = "baz"
)

func (MyEnumIsValid) IsValid() bool {
	// ...
	return false
}
