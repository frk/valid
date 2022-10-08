package generator

import (
	"fmt"

	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

var _ = fmt.Println

func (b *bb) condExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	if len(r.Args) > 0 {
		b.g.prepArgs(n, r)
	}
	if r.Spec.Kind == rules.ENUM {
		b.g.prepEnums(n, r)
	}

	switch r.Spec.Kind {
	case rules.OPTIONAL:
		return b.optionalCondExpr(n, r)
	case rules.REQUIRED:
		return b.requiredCondExpr(n, r)
	case rules.COMPARABLE:
		return b.comparableCondExpr(n, r)
	case rules.ORDERED:
		return b.orderedCondExpr(n, r)
	case rules.LENGTH:
		return b.lengthCondExpr(n, r)
	case rules.RANGE:
		return b.rangeCondExpr(n, r)
	case rules.ENUM:
		return b.enumCondExpr(n, r)
	case rules.FUNCTION:
		return b.functionCondExpr(n, r)
	case rules.METHOD:
		return b.methodCondExpr(n, r)
	}

	panic("shouldn't reach")
	return nil
}

// builds an expression that checks whether the value is empty or not.
func (b *bb) optionalCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	switch n.Type.Kind {
	case gotype.K_STRING:
		return GO.BinaryExpr{Op: GO.BinaryNeq, X: b.val, Y: GO.ValueLit(`""`)}
	case gotype.K_MAP, gotype.K_SLICE:
		return GO.BinaryExpr{Op: GO.BinaryGtr, X: GO.CallLenExpr{b.val}, Y: GO.IntLit(0)}
	case gotype.K_INT, gotype.K_INT8, gotype.K_INT16, gotype.K_INT32, gotype.K_INT64:
		return GO.BinaryExpr{Op: GO.BinaryGtr, X: b.val, Y: GO.IntLit(0)}
	case gotype.K_UINT, gotype.K_UINT8, gotype.K_UINT16, gotype.K_UINT32, gotype.K_UINT64:
		return GO.BinaryExpr{Op: GO.BinaryGtr, X: b.val, Y: GO.IntLit(0)}
	case gotype.K_FLOAT32, gotype.K_FLOAT64:
		return GO.BinaryExpr{Op: GO.BinaryGtr, X: b.val, Y: GO.ValueLit("0.0")}
	case gotype.K_BOOL:
		return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.ValueLit("true")}
	case gotype.K_INTERFACE:
		return GO.BinaryExpr{Op: GO.BinaryNeq, X: b.val, Y: NIL}
	case gotype.K_PTR:
		return GO.BinaryExpr{Op: GO.BinaryNeq, X: b.val, Y: NIL}
	}

	panic("shouldn't reach")
	return nil
}

// builds an expression that checks the value against the "zero" value.
func (b *bb) requiredCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	switch r.Name {
	case "notnil":
		switch n.Type.Kind {
		case gotype.K_MAP, gotype.K_SLICE, gotype.K_INTERFACE:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: NIL}
		case gotype.K_PTR:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: NIL}
		}
	case "required":
		switch n.Type.Kind {
		case gotype.K_STRING:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.ValueLit(`""`)}
		case gotype.K_MAP, gotype.K_SLICE:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: GO.CallLenExpr{b.val}, Y: GO.IntLit(0)}
		case gotype.K_INT, gotype.K_INT8, gotype.K_INT16, gotype.K_INT32, gotype.K_INT64:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.IntLit(0)}
		case gotype.K_UINT, gotype.K_UINT8, gotype.K_UINT16, gotype.K_UINT32, gotype.K_UINT64:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.IntLit(0)}
		case gotype.K_FLOAT32, gotype.K_FLOAT64:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.ValueLit("0.0")}
		case gotype.K_BOOL:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.ValueLit("false")}
		case gotype.K_INTERFACE:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: NIL}
		case gotype.K_PTR:
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: NIL}
		case gotype.K_STRUCT:
			lit := GO.StructLit{Compact: true}
			if n.Type.Pkg != b.g.pkg {
				pkg := b.g.addImport(n.Type.Pkg)
				lit.Type = GO.QualifiedIdent{pkg.name, n.Type.Name}
			} else {
				lit.Type = GO.Ident{n.Type.Name}
			}
			return GO.BinaryExpr{Op: GO.BinaryEql, X: b.val, Y: GO.ParenExpr{lit}}
		}
	}

	panic("shouldn't reach")
	return nil
}

// builds an expression that compares the value against the rule's arguments.
func (b *bb) comparableCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	var binOp, logOp GO.BinaryOp
	switch r.Name {
	case "eq":
		binOp = GO.BinaryNeq
		logOp = GO.BinaryLAnd
	case "ne":
		binOp = GO.BinaryEql
		logOp = GO.BinaryLOr
	}

	args := b.g.argmap[r]
	cond := GO.BinaryExpr{Op: binOp, X: b.val, Y: args[0]}
	for _, a := range args[1:] {
		y := GO.BinaryExpr{Op: binOp, X: b.val, Y: a}
		cond = GO.BinaryExpr{Op: logOp, X: cond, Y: y}
	}
	return cond
}

// builds an expression that compares the value against the rule's argument.
func (b *bb) orderedCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	var binOp GO.BinaryOp
	switch r.Name {
	case "gt":
		binOp = GO.BinaryLeq
	case "lt":
		binOp = GO.BinaryGeq
	case "gte":
		binOp = GO.BinaryLss
	case "lte":
		binOp = GO.BinaryGtr
	case "min":
		binOp = GO.BinaryLss
	case "max":
		binOp = GO.BinaryGtr
	}

	args := b.g.argmap[r]
	return GO.BinaryExpr{Op: binOp, X: b.val, Y: args[0]}
}

// used for adding pkgimport
var utf8 = gotype.Pkg{Path: "unicode/utf8", Name: "utf8"}

// builds an expression that checks the value's length.
func (b *bb) lengthCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	var cx GO.ExprNode
	switch r.Name {
	case "len":
		cx = GO.CallLenExpr{b.val}
	case "runecount":
		pkg := b.g.addImport(utf8)
		expr := GO.CallExpr{Args: GO.ArgsList{List: b.val}}
		if n.Type.Kind == gotype.K_STRING {
			expr.Fun = GO.QualifiedIdent{pkg.name, "RuneCountInString"}
		} else if n.Type.Kind == gotype.K_SLICE && n.Type.Elem.IsByte {
			expr.Fun = GO.QualifiedIdent{pkg.name, "RuneCount"}
		} else {
			panic("shouldn't reach")
		}
		cx = expr
	}

	switch args := b.g.argmap[r]; {
	case len(args) == 1:
		return GO.BinaryExpr{Op: GO.BinaryNeq, X: cx, Y: args[0]}
	case len(args) == 2 && r.Args[1].Value == "":
		return GO.BinaryExpr{Op: GO.BinaryLss, X: cx, Y: args[0]}
	case len(args) == 2 && r.Args[0].Value == "":
		return GO.BinaryExpr{Op: GO.BinaryGtr, X: cx, Y: args[1]}
	case len(args) == 2:
		lss := GO.BinaryExpr{Op: GO.BinaryLss, X: cx, Y: args[0]}
		gtr := GO.BinaryExpr{Op: GO.BinaryGtr, X: cx, Y: args[1]}
		return GO.BinaryExpr{Op: GO.BinaryLOr, X: lss, Y: gtr}
	}

	panic("shouldn't reach")
	return nil
}

// builds an expression that checks the value's numeric range.
func (b *bb) rangeCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	args := b.g.argmap[r]
	lss := GO.BinaryExpr{Op: GO.BinaryLss, X: b.val, Y: args[0]}
	gtr := GO.BinaryExpr{Op: GO.BinaryGtr, X: b.val, Y: args[1]}
	return GO.BinaryExpr{Op: GO.BinaryLOr, X: lss, Y: gtr}
}

// builds an expression that checks the value against a set of enums.
func (b *bb) enumCondExpr(n *rules.Node, r *rules.Rule) (x GO.ExprNode) {
	enums := b.g.enumap[r]
	x = GO.BinaryExpr{Op: GO.BinaryNeq, X: b.val, Y: enums[0]}
	for _, en := range enums[1:] {
		cond := GO.BinaryExpr{Op: GO.BinaryNeq, X: b.val, Y: en}
		x = GO.BinaryExpr{Op: GO.BinaryLAnd, X: x, Y: cond}
	}
	return x
}

// builds an expression that checks the value using the rule spec's function.
func (b *bb) functionCondExpr(n *rules.Node, r *rules.Rule) (x GO.ExprNode) {
	pkg := b.g.addImport(r.Spec.FType.Pkg)
	args := b.g.argmap[r]

	// If this is the included regexp rule, then add
	// a registry call statement for the init function.
	if r.Spec.Name == "re" && r.Spec.FType.IsIncluded() {
		b.g.init = append(b.g.init, r)
	}

	if r.Spec.JoinOp == 0 {
		cx := GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, r.Spec.FName}}
		cx.Args.List = append(GO.ExprList{b.val}, args...)
		return GO.UnaryExpr{Op: GO.UnaryNot, X: cx}
	}

	if r.Spec.JoinOp > 0 {
		for _, a := range args {
			cx := GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, r.Spec.FName}}
			cx.Args.List = GO.ExprList{b.val, a}

			switch {
			// x || x...
			case r.Spec.JoinOp == rules.JOIN_NOT && x != nil:
				x = GO.BinaryExpr{Op: GO.BinaryLOr, X: x, Y: cx}
			case r.Spec.JoinOp == rules.JOIN_NOT && x == nil:
				x = cx

			// !x || !x...
			case r.Spec.JoinOp == rules.JOIN_AND && x != nil:
				x = GO.BinaryExpr{Op: GO.BinaryLOr, X: x, Y: GO.UnaryExpr{Op: GO.UnaryNot, X: cx}}
			case r.Spec.JoinOp == rules.JOIN_AND && x == nil:
				x = GO.UnaryExpr{Op: GO.UnaryNot, X: cx}

			// !x && !x...
			case r.Spec.JoinOp == rules.JOIN_OR && x != nil:
				x = GO.BinaryExpr{Op: GO.BinaryLAnd, X: x, Y: GO.UnaryExpr{Op: GO.UnaryNot, X: cx}}
			case r.Spec.JoinOp == rules.JOIN_OR && x == nil:
				x = GO.UnaryExpr{Op: GO.UnaryNot, X: cx}
			}
		}
	}
	return x
}

// builds an expression that checks the value by invoking the designated method.
func (b *bb) methodCondExpr(n *rules.Node, r *rules.Rule) GO.ExprNode {
	x := b.val
	if n.PtrDepth() > 1 || (n.PtrDepth() == 1 && n.Type.Kind == gotype.K_INTERFACE) {
		x = GO.ParenExpr{x}
	}

	if n.PtrDepth() == 1 && n.Type.Kind != gotype.K_INTERFACE {
		x = b.vals[len(b.vals)-1]
	}

	// TODO if b.val's type is *some_type then we should either put parens
	// around, or, better yet, we should remove the indirection because calling
	// the method directly is shorthand for (*x).M()...
	call := GO.CallExpr{Fun: GO.SelectorExpr{X: x, Sel: GO.Ident{r.Spec.FName}}}
	return GO.UnaryExpr{Op: GO.UnaryNot, X: call}
}
