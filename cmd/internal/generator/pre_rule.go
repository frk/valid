package generator

import (
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (b *bb) preruleStmt(n *rules.Node) {
	rhs := b.preprocCall(n)
	n.PreRules = nil

	lhs := b.val
	if b.idx != nil {
		lhs = b.idx
	}

	b.add(GO.AssignStmt{
		Token: GO.Assign,
		Lhs:   lhs,
		Rhs:   rhs,
	})

	if b.idx != nil && len(n.IsRules) > 0 {
		b.add(GO.AssignStmt{
			Token: GO.Assign,
			Lhs:   b.val,
			Rhs:   lhs,
		})
	}
}

func (b *bb) preprocCall(n *rules.Node) GO.CallExpr {
	for _, r := range n.PreRules {
		b.g.prepArgs(n, r)
	}

	r := n.PreRules[0]
	pkg := b.g.addImport(r.Spec.FType.Pkg)
	args := append(GO.ExprList{b.val}, b.g.argmap[r]...)

	call := GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, r.Spec.FName}}
	call.Args.List = args
	for _, r := range n.PreRules[1:] {
		pkg := b.g.addImport(r.Spec.FType.Pkg)
		args := append(GO.ExprList{call}, b.g.argmap[r]...)

		call = GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, r.Spec.FName}}
		call.Args.List = args
	}
	return call
}
