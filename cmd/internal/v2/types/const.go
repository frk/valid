package types

import (
	"sync"

	"github.com/frk/valid/cmd/internal/v2/source"
)

// Const represents the identifier of a declared constant.
type Const struct {
	// The package where the constant is declared.
	Pkg Pkg
	// The constant's type.
	Type *Type
	// Name of the constant.
	Name string
}

// FindConsts is a helper method that finds and
// returns all declared constants for a given type.
func FindConsts(t *Type, src *source.Source) (consts []Const) {
	ident := t.Pkg.Path + "." + t.Name

	constCache.RLock()
	consts, ok := constCache.m[ident]
	constCache.RUnlock()

	if ok { // already done?
		return consts
	} else {
		consts = nil
	}

	for _, c := range src.FindConsts(t.Pkg.Path, t.Name) {
		name := c.Name()
		// blank, skip
		if name == "_" {
			continue
		}

		pkg := c.Pkg()
		p := Pkg{Path: pkg.Path(), Name: pkg.Name()}
		k := Const{Pkg: p, Type: t, Name: name}
		consts = append(consts, k)
	}

	constCache.Lock()
	constCache.m[ident] = consts
	constCache.Unlock()

	return consts
}

////////////////////////////////////////////////////////////////////////////////
// cache

var constCache = struct {
	sync.RWMutex
	// m maps package-path qualified type names to
	// a slice of constants declared with that type.
	m map[string][]Const
}{m: make(map[string][]Const)}
