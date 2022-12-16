package generator

import (
	"strings"

	"github.com/frk/valid/cmd/internal/rules"
	"github.com/frk/valid/cmd/internal/xtypes"

	GO "github.com/frk/ast/golang"
)

func fileAST(pkg xtypes.Pkg, infos []*rules.Info) GO.File {
	g := &gg{pkg: pkg}
	g.argmap = make(map[*rules.Rule][]GO.ExprNode)
	g.enumap = make(map[*rules.Rule][]GO.ExprNode)

	f := GO.File{
		Preamble: FILE_PREAMBLE,
		PkgName:  pkg.Name,
	}
	for _, info := range infos {
		f.Decls = append(f.Decls, methodAST(g, info))
	}
	if len(g.init) > 0 {
		f.Decls = append([]GO.TopLevelDeclNode{initAST(g)}, f.Decls...)
	}
	if len(g.imports) > 0 {
		f.Imports = importsAST(g)
	}
	return f
}

func initAST(g *gg) GO.FuncDecl {
	init := GO.FuncDecl{Name: GO.Ident{"init"}}
	for _, r := range g.init {
		pkg := g.addImport(r.Spec.FType.Pkg)
		args := g.argmap[r]
		call := GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, "RegisterRegexp"}}
		call.Args.List = args[0]
		init.Body.Add(GO.ExprStmt{call})
	}

	return init
}

// importsAST produces an import declaration for packages that need to be imported.
func importsAST(g *gg) []GO.ImportDeclNode {
	var specs []GO.ImportSpec
	for _, pkg := range g.imports {
		spec := GO.ImportSpec{Path: GO.StringLit(pkg.path)}
		if pkg.local {
			spec.Name.Name = pkg.name
		}
		specs = append(specs, spec)
	}

	// Split the imports into 3 groups separated by a new line.
	// - the 1st group will contain imports from the standard library
	// - the 3rd group will contain imports from github.com/frk/valid...
	// - and the 2nd group will contain the rest of the imports.
	var sp1, sp2, sp3 []GO.ImportSpec
	for _, s := range specs {
		if strings.HasPrefix(string(s.Path), `github.com/frk/valid`) {
			sp3 = append(sp3, s)
		} else if i := strings.IndexByte(string(s.Path), '.'); i >= 0 {
			sp2 = append(sp2, s)
		} else {
			sp1 = append(sp1, s)
		}
	}

	specs = nil
	if len(sp1) > 0 {
		specs = append(specs, sp1...)
	}
	if len(sp2) > 0 {
		sp2[0].Doc = GO.NL{}
		specs = append(specs, sp2...)
	}
	if len(sp3) > 0 {
		sp3[0].Doc = GO.NL{}
		specs = append(specs, sp3...)
	}

	return []GO.ImportDeclNode{&GO.ImportDecl{Specs: specs}}
}
