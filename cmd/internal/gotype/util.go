package gotype

import (
	"github.com/frk/valid/cmd/internal/search"
)

func MustGetType(pkgpath, name string, a *search.AST) *Type {
	obj, err := search.FindObject(pkgpath, name, a)
	if err != nil {
		panic(err)
	}

	an := NewAnalyzer(obj.Pkg())
	return an.Object(obj)
}
