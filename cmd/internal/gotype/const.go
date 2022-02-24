package gotype

import (
	"sync"

	"github.com/frk/valid/cmd/internal/search"
)

// Const represents the identifier of a declared constant.
type Const struct {
	// The constant's package.
	Pkg Pkg
	// Name of the constant.
	Name string
}

// Consts is a helper method that finds and
// returns all declared constants for a given type.
func (a *Analyzer) Consts(t *Type, ast *search.AST) (consts []Const) {
	ident := t.Pkg.Path + "." + t.Name

	constCache.RLock()
	consts, ok := constCache.m[ident]
	constCache.RUnlock()

	if ok { // already done?
		return consts
	} else {
		consts = nil
	}

	for _, c := range search.FindConstantsByType(t.Pkg.Path, t.Name, ast) {
		name := c.Name()
		// blank, skip
		if name == "_" {
			continue
		}

		pkg := c.Pkg()
		// imported but not exported, skip
		if a.pkg.Path != pkg.Path() && !c.Exported() {
			continue
		}

		consts = append(consts, Const{
			Name: name,
			Pkg: Pkg{
				Path: pkg.Path(),
				Name: pkg.Name(),
			},
		})
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
