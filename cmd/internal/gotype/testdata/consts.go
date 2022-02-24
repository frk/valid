package testdata

type ConstType1 uint

const (
	CT1 ConstType1 = iota
	CT2
	_
	CT4
	ct5 // unexported (but local to the TestValidator)
)

type ConstType2 string

const (
	ConstFoo  ConstType2 = "foo"
	ConstBar  ConstType2 = "bar"
	const_baz ConstType2 = "baz" // unexported (but local to the TestValidator)
)

type constType3 bool

const (
	kYES constType3 = true
	kNO  constType3 = false
)
