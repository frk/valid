package testdata

import (
	"go/parser"
)

type K1 uint
type K2 uint
type K3 string

const (
	k11 K1 = iota + 1
	k12
	k13

	k21 K2 = iota + 1
	k22
	k23

	k3foo K3 = "foo"
	k3bar K3 = "bar"

	k14 K1 = 4

	k3baz = K3("bar")
)

const k1111 = K1(111)

const CustomMode = parser.ParseComments | parser.Trace | parser.DeclarationErrors
