package testdata

import (
	"go/ast"
)

type mytype string

const (
	foo mytype = "foo"
	bar mytype = "bar"
	baz mytype = "baz"
)

type Validator struct {
	F1 mytype      `is:"enum"`
	F2 ast.ObjKind `is:"enum"`
}
