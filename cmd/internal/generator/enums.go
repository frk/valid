package generator

import (
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (g *gg) prepEnums(n *rules.Node, r *rules.Rule) {
	var pkgname string
	if g.info.Validator.Type.Pkg.Path != n.Type.Pkg.Path {
		pkg := g.pkg(n.Type.Pkg)
		pkgname = pkg.name
	}

	enums := g.info.EnumMap[n.Type]
	exprs := make([]GO.ExprNode, len(enums))
	for i, en := range enums {
		expr := GO.ExprNode(GO.Ident{en.Name})
		if len(pkgname) > 0 {
			expr = GO.QualifiedIdent{pkgname, en.Name}
		}
		exprs[i] = expr
	}

	g.enumap[r] = exprs
}
