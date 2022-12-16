package generator

import (
	"strconv"

	"github.com/frk/valid/cmd/internal/rules"
	"github.com/frk/valid/cmd/internal/xtypes"

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
		tt = &xtypes.Type{Kind: xtypes.K_INT}
	case rules.FUNCTION, rules.PREPROC:
		if in := r.Spec.FType.In; len(in)-1 <= i {
			tt = in[len(in)-1].Type
		} else {
			tt = in[i+1].Type
		}
	}

	if a.Type == rules.ARG_FIELD_ABS || a.Type == rules.ARG_FIELD_REL {
		return g.fieldArg(a, tt)
	}
	return g.constArg(n, r, a, tt)
}

func (g *gg) fieldArg(a *rules.Arg, t *xtypes.Type) GO.ExprNode {
	f := g.info.KeyMap[a.Value]
	x := g.recv

	var last *xtypes.StructField
	for _, f := range f.Selector {
		x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
		last = f
	}

	if t.PtrOf(last.Object.Type) {
		return GO.UnaryExpr{Op: GO.UnaryAmp, X: x}
	}
	if last.Object.Type.NeedsConversion(t) {
		T := GO.Ident{t.TypeString(nil)}
		return GO.CallExpr{Fun: T, Args: GO.ArgsList{List: x}}
	}
	return x
}

func (g *gg) constArg(n *rules.Node, r *rules.Rule, a *rules.Arg, t *xtypes.Type) (x GO.ExprNode) {
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

func zeroValue(t *xtypes.Type) (x GO.ExprNode) {
	switch t.Kind {
	case xtypes.K_STRING:
		return GO.StringLit("")
	case xtypes.K_INT, xtypes.K_INT8, xtypes.K_INT16, xtypes.K_INT32, xtypes.K_INT64:
		return GO.IntLit(0)
	case xtypes.K_UINT, xtypes.K_UINT8, xtypes.K_UINT16, xtypes.K_UINT32, xtypes.K_UINT64:
		return GO.IntLit(0)
	case xtypes.K_FLOAT32, xtypes.K_FLOAT64:
		return GO.ValueLit("0.0")
	case xtypes.K_BOOL:
		return GO.ValueLit("false")
	case xtypes.K_PTR, xtypes.K_SLICE, xtypes.K_MAP, xtypes.K_INTERFACE:
		return NIL
	}

	panic("shouldn't reach")
	return nil
}
