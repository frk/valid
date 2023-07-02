package types

import (
	"go/types"

	"github.com/frk/valid/cmd/internal/v2/source"
)

// Func is used to represent a function.
type Func struct {
	Name string
	Type *Type
}

func AnalyzeFunc(fun *types.Func, src *source.Source) *Func {
	a := &analyzer{src: src}
	t := a.analyzeType(fun.Type())
	if pkg := fun.Pkg(); pkg != nil {
		t.Pkg.Path = pkg.Path()
		t.Pkg.Name = pkg.Name()
	}

	f := new(Func)
	f.Name = fun.Name()
	f.Type = t
	return f
}
