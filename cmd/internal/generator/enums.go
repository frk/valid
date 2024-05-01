package generator

import (
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (g *gg) prepEnums(n *rules.Node, r *rules.Rule) {
	pkg := g.addImport(n.Type.Pkg)
	enums := g.info.EnumMap[n.Type]
	exprs := make([]GO.ExprNode, len(enums))
	for i, en := range enums {
		exprs[i] = pkgQualIdent(pkg, en.Name)
	}

	g.enumap[r] = exprs
}
