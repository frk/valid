package generator

import (
	"strconv"

	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (g *gg) prepArgs(n *rules.Node, r *rules.Rule) {
	args := make([]GO.ExprNode, len(r.Args))
	for i, a := range r.Args {
		args[i] = g.ruleArg(n, r, i, a)
	}
	g.argmap[r] = args
}

func (g *gg) ruleArg(n *rules.Node, r *rules.Rule, i int, a *rules.Arg) GO.ExprNode {
	tt := n.Type // target type

	switch r.Spec.Kind {
	case rules.LENGTH:
		tt = &gotype.Type{Kind: gotype.K_INT}
	case rules.FUNCTION, rules.PREPROC:
		if in := r.Spec.FType.In; len(in)-1 <= i {
			tt = in[len(in)-1].Type
		} else {
			tt = in[i+1].Type
		}
	}

	if a.Type == rules.ARG_FIELD {
		return g.fieldArg(a, tt)
	}
	return g.constArg(n, r, a, tt)
}

func (g *gg) fieldArg(a *rules.Arg, t *gotype.Type) GO.ExprNode {
	f := g.info.KeyMap[a.Value]
	x := g.recv

	var last *gotype.StructField
	for _, f := range f.Selector {
		x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
		last = f
	}

	if t.PtrOf(last.Type) {
		return GO.UnaryExpr{Op: GO.UnaryAmp, X: x}
	}
	if last.Type.NeedsConversion(t) {
		T := GO.Ident{t.TypeString(nil)}
		return GO.CallExpr{Fun: T, Args: GO.ArgsList{List: x}}
	}
	return x
}

func (g *gg) constArg(n *rules.Node, r *rules.Rule, a *rules.Arg, t *gotype.Type) (x GO.ExprNode) {
	if t.IsEmptyInterface() && a.Type != rules.ARG_STRING {
		return GO.ValueLit(a.Value)
	}
	if a.Type == rules.ARG_STRING && r.Spec.UseRawString {
		return GO.RawStringLit(a.Value)
	}

	switch a.Type {
	case rules.ARG_STRING:
		return GO.ValueLit(strconv.Quote(a.Value))
	case rules.ARG_BOOL, rules.ARG_INT, rules.ARG_FLOAT:
		return GO.ValueLit(a.Value)
	case rules.ARG_UNKNOWN:
		return zeroValue(t)
	}

	panic("shouldn't reach")
	return nil
}

func zeroValue(t *gotype.Type) (x GO.ExprNode) {
	switch t.Kind {
	case gotype.K_STRING:
		return GO.StringLit("")
	case gotype.K_INT, gotype.K_INT8, gotype.K_INT16, gotype.K_INT32, gotype.K_INT64:
		return GO.IntLit(0)
	case gotype.K_UINT, gotype.K_UINT8, gotype.K_UINT16, gotype.K_UINT32, gotype.K_UINT64:
		return GO.IntLit(0)
	case gotype.K_FLOAT32, gotype.K_FLOAT64:
		return GO.ValueLit("0.0")
	case gotype.K_BOOL:
		return GO.ValueLit("false")
	case gotype.K_PTR, gotype.K_SLICE, gotype.K_MAP, gotype.K_INTERFACE:
		return NIL
	}

	panic("shouldn't reach")
	return nil
}
