package types

import (
	"go/types"

	"github.com/frk/valid/cmd/internal/search"
)

func MustGetType(pkgpath, name string, a *search.AST) *Type {
	o, err := search.FindObject(pkgpath, name, a)
	if err != nil {
		panic(err)
	}
	tn, ok := o.(*types.TypeName)
	if !ok {
		panic("object is not type name")
	}
	named, ok := tn.Type().(*types.Named)
	if !ok {
		panic("type name is not named")
	}

	return Analyze(named, a)
}
