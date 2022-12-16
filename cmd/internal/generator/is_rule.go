package generator

import (
	"fmt"

	"github.com/frk/valid/cmd/internal/rules"
	"github.com/frk/valid/cmd/internal/xtypes"

	GO "github.com/frk/ast/golang"
)

var _ = fmt.Println

func (b *bb) isruleStmt(n *rules.Node) {
	switch {
	case b.optIfs == nil && n.IsOptional():
		b.optIfs = &GO.IfStmt{}
		b.optIfs.Cond = b.condExpr(n, n.IsRules[0])
		b.add(b.optIfs)
		b.use(&b.optIfs.Body)
		n.IsRules = n.IsRules[1:]

	case b.reqIfs == nil && n.IsRequired():
		// TODO
	}

	rr := n.IsRules
	n.IsRules = nil
	basicOnly := true

	ifs := &GO.IfStmt{}
	cur := ifs
	for _, r := range rr {
		if cur.Cond != nil {
			elif := &GO.IfStmt{}
			cur.Else = elif
			cur = elif
		}

		if r.IsBasic() {
			cur.Cond = b.condExpr(n, r)
			b.err(n, r, &cur.Body)
		} else {
			basicOnly = false

			// if ok, err := <call_expr>; err != nil {
			//	return err
			// } else if !ok {
			//	return <custom_error>
			// }
			cur.Init = GO.AssignStmt{
				Token: GO.AssignDefine,
				Lhs:   GO.IdentList{OK, ERR},
				Rhs:   b.functionCallExpr(n, r),
			}
			cur.Cond = GO.BinaryExpr{Op: GO.BinaryNeq, X: ERR, Y: NIL}
			cur.Body = GO.BlockStmt{[]GO.StmtNode{GO.ReturnStmt{ERR}}}

			elif := &GO.IfStmt{}
			cur.Else = elif
			cur = elif
			cur.Cond = GO.UnaryExpr{Op: GO.UnaryNot, X: OK}
			b.err(n, r, &cur.Body)
		}
	}

	switch {
	case b.optIfs != nil && len(rr) < 2 && basicOnly:
		cond := ifs.Cond
		if wantsParens(rr[0]) {
			cond = GO.ParenExpr{cond}
		}

		b.optIfs.Cond = binAnd(b.optIfs.Cond, cond)
		b.optIfs.Body.Add(ifs.Body.List...)
		b.optIfs.Else = ifs.Else

	case b.reqIfs != nil:
		b.reqIfs.Else = ifs
	default:
		b.add(ifs)
	}

	if n.HasSubRules() {
		els := &GO.BlockStmt{}
		cur.Else = els
		b.use(els)
	}
}

func (b *bb) ptrOptionalStmt(n *rules.Node) (base *rules.Node) {
	ifs := &GO.IfStmt{}
	for base = n; base.Type.Kind == xtypes.K_PTR; {
		if !base.IsOptional() {
			panic("shouldn't happen")
		}

		if ifs.Cond != nil {
			ifs.Cond = binAnd(ifs.Cond, b.condExpr(base, base.IsRules[0]))
		} else {
			ifs.Cond = b.condExpr(base, base.IsRules[0])
		}

		base = base.Elem
		b.pind()
	}

	b.add(ifs)
	b.use(&ifs.Body)
	b.optIfs = ifs

	// extend; if base is also optional
	if base.PreRules.Empty() && base.IsOptional() {
		ifs.Cond = binAnd(ifs.Cond, b.condExpr(base, base.IsRules[0]))
		base.IsRules = base.IsRules[1:]
	}

	return base
}

func (b *bb) ptrNoGuardStmt(n *rules.Node) (base *rules.Node) {
	for base = n; base.Type.Kind == xtypes.K_PTR; {
		if !base.IsNoGuard() {
			panic("shouldn't happen")
		}
		base = base.Elem
		b.pind()
	}

	return base
}

func (b *bb) ptrRequiredStmt(n *rules.Node) (base *rules.Node) {
	r0 := n.IsRules[0]

	ifs := &GO.IfStmt{}
	for base = n; base.Type.Kind == xtypes.K_PTR; {
		if !base.IsRequired() {
			panic("shouldn't happen")
		}
		if ifs.Cond != nil {
			ifs.Cond = binOr(ifs.Cond, b.condExpr(base, base.IsRules[0]))
		} else {
			ifs.Cond = b.condExpr(base, base.IsRules[0])
		}
		base = base.Elem
		b.pind()
	}

	b.add(ifs)
	b.err(base, r0, &ifs.Body)
	b.reqIfs = ifs

	// extend; if base is also required
	if base.PreRules.Empty() && base.IsRequired() {
		ifs.Cond = binOr(ifs.Cond, b.condExpr(base, base.IsRules[0]))
		base.IsRules = base.IsRules[1:]
	}

	if len(base.PreRules) > 0 || base.NeedsTempVar() || base.HasSubRules() {
		els := &GO.BlockStmt{}
		ifs.Else = els
		b.use(els)
		b.reqIfs = nil
	}

	return base
}

////////////////////////////////////////////////////////////////////////////////
// convenience

func binAnd(x, y GO.ExprNode) GO.BinaryExpr {
	return GO.BinaryExpr{Op: GO.BinaryLAnd, X: x, Y: y}
}

func binOr(x, y GO.ExprNode) GO.BinaryExpr {
	return GO.BinaryExpr{Op: GO.BinaryLOr, X: x, Y: y}
}

func wantsParens(r *rules.Rule) bool {
	return r.Spec.Kind == rules.ENUM || len(r.Args) > 1 &&
		(r.Spec.Kind != rules.FUNCTION || r.Spec.JoinOp > 0)
}
