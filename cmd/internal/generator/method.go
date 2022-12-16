package generator

import (
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func methodAST(g *gg, info *rules.Info) (dec GO.MethodDecl) {
	g.info = info
	g.recv = GO.Ident{"v"}

	dec.Recv.Name = g.recv.(GO.Ident)
	dec.Recv.Type = GO.Ident{g.info.Validator.Type.Name}
	dec.Name.Name = "Validate"
	dec.Type.Results = GO.ParamList{{Type: ERROR}}

	if gotype.Globals.ErrorAggregator != nil && info.Validator.ErrorHandlerField == nil {
		newErrorAggregatorAST(info, g.block(&dec.Body))
	}
	if before := info.Validator.BeforeValidateMethod; before != nil {
		hookAST(before, g.block(&dec.Body))
	}

	nodesAST(info.RootNode, g.block(&dec.Body))
	exitAST(info, g.block(&dec.Body))

	return dec
}

func exitAST(info *rules.Info, b bb) {
	stmt := GO.ReturnStmt{Result: NIL}
	errh := info.Validator.ErrorHandlerField
	after := info.Validator.AfterValidateMethod

	hasAgg := false
	if (errh != nil && errh.IsAggregator) || (errh == nil && gotype.Globals.ErrorAggregator != nil) {
		hasAgg = true
	}

	switch {
	case hasAgg == true && after != nil:
		call := GO.CallExpr{}
		if h := info.Validator.ErrorHandlerField; h != nil && h.IsAggregator {
			x := GO.SelectorExpr{X: b.g.recv, Sel: GO.Ident{h.Name}}
			call.Fun = GO.SelectorExpr{X: x, Sel: GO.Ident{"Out"}}
		} else if h == nil && gotype.Globals.ErrorAggregator != nil {
			x := GO.SelectorExpr{X: GO.Ident{"ea"}, Sel: GO.Ident{"Out"}}
			call.Fun = x
		}

		ifs := new(GO.IfStmt)
		ifs.Init = GO.AssignStmt{Token: GO.AssignDefine, Lhs: ERR, Rhs: call}
		ifs.Cond = GO.BinaryExpr{Op: GO.BinaryNeq, X: ERR, Y: NIL}
		ifs.Body = GO.BlockStmt{[]GO.StmtNode{GO.ReturnStmt{ERR}}}

		b.add(ifs)
		hookAST(after, b)

	case hasAgg == true && after == nil:
		if h := info.Validator.ErrorHandlerField; h != nil && h.IsAggregator {
			x := GO.SelectorExpr{X: b.g.recv, Sel: GO.Ident{h.Name}}
			x = GO.SelectorExpr{X: x, Sel: GO.Ident{"Out"}}
			stmt = GO.ReturnStmt{Result: GO.CallExpr{Fun: x}}
		} else if h == nil && gotype.Globals.ErrorAggregator != nil {
			x := GO.SelectorExpr{X: GO.Ident{"ea"}, Sel: GO.Ident{"Out"}}
			stmt = GO.ReturnStmt{Result: GO.CallExpr{Fun: x}}
		}

	case hasAgg == false && after != nil:
		hookAST(after, b)

	}

	b.add(stmt)
}

func newErrorAggregatorAST(info *rules.Info, b bb) {
	agg := gotype.Globals.ErrorAggregator
	pkg := b.g.addImport(agg.Pkg)

	typ := GO.QualifiedIdent{pkg.name, agg.Name}
	call := GO.CallExpr{Fun: GO.Ident{"new"}}
	call.Args = GO.ArgsList{List: GO.ExprList{typ}}

	b.add(GO.AssignStmt{Token: GO.AssignDefine, Lhs: GO.Ident{"ea"}, Rhs: call})
	b.add(GO.NL{})
}
