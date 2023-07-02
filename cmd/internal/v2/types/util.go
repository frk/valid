package types

import (
	"go/types"

	"github.com/frk/valid/cmd/internal/v2/source"
)

func MustGetType(pkgpath, name string, src *source.Source) *Type {
	o, err := src.FindObject(pkgpath, name)
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

	return Analyze(named, src)
}
