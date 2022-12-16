package mypkg

type ConstType1 uint

const (
	CT1 ConstType1 = iota
	CT2
	_
	CT4
	ct5 // unexported
)
