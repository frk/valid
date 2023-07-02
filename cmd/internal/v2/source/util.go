package source

import (
	"go/types"
)

func (f Func) ParamsLen() int {
	s := f.Type().(*types.Signature)
	return s.Params().Len()
}

func (f Func) IsVariadic() bool {
	s := f.Type().(*types.Signature)
	return s.Variadic()
}
