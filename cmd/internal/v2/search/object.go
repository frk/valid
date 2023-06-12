package search

import (
	"go/ast"
	"go/token"
	"go/types"
)

// FindObject returns a top-level declared object that matches
// the given pkgpath and name. The returned object will either
// be a top-level declared type or a top-level declared function.
func FindObject(pkgpath, name string, a *AST) (obj types.Object, err error) {
	pkg, err := findpkg(pkgpath, name, a)
	if err != nil {
		return nil, err
	}

	for _, syn := range pkg.Syntax {
		for _, d := range syn.Decls {
			switch d := d.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, s := range d.Specs {
					s, ok := s.(*ast.TypeSpec)
					if !ok || s.Name.Name != name {
						continue
					}
					def, ok := pkg.TypesInfo.Defs[s.Name]
					if !ok {
						continue
					}
					if tn, ok := def.(*types.TypeName); ok {
						return tn, nil
					}
				}
			case *ast.FuncDecl:
				if d.Recv != nil || d.Name.Name != name {
					continue
				}
				def, ok := pkg.TypesInfo.Defs[d.Name]
				if !ok {
					continue
				}
				if f, ok := def.(*types.Func); ok {
					return f, nil
				}
			}
		}
	}

	return nil, &Error{C: ERR_OBJECT_NOTFOUND, pkg: pkgpath, name: name}
}
