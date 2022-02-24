package generator

import (
	"fmt"

	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"
)

var _ = fmt.Println

func nodesAST(n *rules.Node, b bb) {
	if !n.HasRules() {
		return
	}

	switch {
	case n.IsPtr() && n.IsOptional() && len(n.Base().PreRules) > 0:
		n = b.ptrOptionalStmt(n)
		b.preruleStmt(n)
		if n.NeedsTempVar() {
			b.tempVar()
		}
		nodesAST(n, b.new())

	case n.IsPtr() && n.IsRequired() && len(n.Base().PreRules) > 0:
		n = b.ptrRequiredStmt(n)
		b.preruleStmt(n)
		if n.NeedsTempVar() {
			b.tempVar()
		}
		nodesAST(n, b.new())

	case n.IsPtr() && n.IsNoGuard() && len(n.Base().PreRules) > 0:
		n = b.ptrNoGuardStmt(n)
		// TODO?
		nodesAST(n, b.new())

	case n.IsPtr() && n.IsOptional():
		n = b.ptrOptionalStmt(n)
		if n.NeedsTempVar() {
			b.tempVar()
		}
		if !n.IsRules.Empty() {
			b.isruleStmt(n)
		}
		if n.IsStruct() && n.PtrDepth() < 2 {
			b.pund()
		}
		nodesAST(n, b.new())

	case n.IsPtr() && n.IsRequired():
		n = b.ptrRequiredStmt(n)
		if n.NeedsTempVar() {
			b.tempVar()
		}
		if !n.IsRules.Empty() {
			b.isruleStmt(n)
		}
		if n.IsStruct() && n.PtrDepth() < 2 {
			b.pund()
		}
		nodesAST(n, b.new())

	case n.IsPtr() && n.IsNoGuard():
		n = b.ptrNoGuardStmt(n)
		if n.NeedsTempVar() {
			b.subBlock()
			b.tempVar()
		}
		if !n.IsRules.Empty() {
			b.isruleStmt(n)
		}
		if n.IsStruct() && n.PtrDepth() < 2 {
			b.pund()
		}
		nodesAST(n, b.new())

	case !n.IsPtr() && !n.PreRules.Empty() && n.IsRules.Empty():
		b.preruleStmt(n)
		nodesAST(n, b.new())

	case !n.IsPtr() && !n.PreRules.Empty() && !n.IsRules.Empty():
		b.preruleStmt(n)
		b.isruleStmt(n)
		nodesAST(n, b.new())

	case !n.IsPtr() && n.PreRules.Empty() && !n.IsRules.Empty():
		b.isruleStmt(n)
		nodesAST(n, b.new())

	case n.Type.Is(gotype.K_ARRAY, gotype.K_SLICE):
		rc := b.arrayForStmt(n)
		nodesAST(n.Elem, b.with(rc.Value))

	case n.Type.Is(gotype.K_MAP):
		rc := b.mapForStmt(n)
		nodesAST(n.Key, b.with(rc.Key))
		nodesAST(n.Elem, b.with(rc.Value))

	case n.Type.Is(gotype.K_STRUCT):
		for _, f := range n.Fields {
			nodesAST(f.Type, b.field(f))
		}
	default:
		panic("shouldn't reach")
	}
}
