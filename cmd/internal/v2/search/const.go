package search

import (
	"go/ast"
	"go/token"
	"go/types"
)

// FindConstantsByType scans the given AST looking for all declared constants
// of the type identified by pkgpath and name. On success the result will be
// a slice of go/types.Const instances that represent those constants.
//
// FindConstantsByType is exepcted to be invoked *after* Search and the AST argument is
// expected to be the same as the one given to Search for caching the packages it loads.
func FindConstantsByType(pkgpath, name string, a *AST) (consts []*types.Const) {
	for _, pkg := range a.pkgs {
		if pkg.PkgPath != pkgpath {
			if _, ok := pkg.Imports[pkgpath]; !ok {
				// If pkg is not the target package, and it also
				// does not import the target package, go to next
				continue
			}
		}

		for _, syn := range pkg.Syntax {
			for _, dec := range syn.Decls {
				gd, ok := dec.(*ast.GenDecl)
				if !ok || gd.Tok != token.CONST {
					continue
				}

				for _, spec := range gd.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					for _, id := range valueSpec.Names {
						obj, ok := pkg.TypesInfo.Defs[id]
						if !ok {
							continue
						}

						if c, ok := obj.(*types.Const); ok {
							named, ok := c.Type().(*types.Named)
							if !ok {
								continue
							}

							tn := named.Obj()
							if tn.Name() != name || tn.Pkg().Path() != pkgpath {
								continue
							}

							consts = append(consts, c)
						}
					}
				}
			}

		}
	}
	return consts
}
