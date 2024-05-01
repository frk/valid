package generator

import (
	"strconv"

	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (b *bb) prepArgs(n *rules.Node, r *rules.Rule) {
	args := make([]GO.ExprNode, len(r.Args))
	for i, a := range r.Args {
		args[i] = b.ruleArg(n, r, i, a)
	}
	b.g.argmap[r] = args
}

func (b *bb) ruleArg(n *rules.Node, r *rules.Rule, i int, a *rules.Arg) GO.ExprNode {
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
		if r.Spec.FType.IsVariadic {
			tt = tt.Elem
		}
	}

	if a.Type == rules.ARG_FIELD_ABS || a.Type == rules.ARG_FIELD_REL {
		return b.fieldArg(a, tt)
	}
	return b.constArg(n, r, a, tt)
}

func (b *bb) fieldArg(a *rules.Arg, t *gotype.Type) GO.ExprNode {
	x, leaf := b.fieldArgSelector(a)
	if t.PtrOf(leaf.Type) {
		return GO.UnaryExpr{Op: GO.UnaryAmp, X: x}
	}
	if leaf.Type.NeedsConversion(t) {
		T := GO.Ident{t.TypeString(nil)}
		return GO.CallExpr{Fun: T, Args: GO.ArgsList{List: x}}
	}
	return x
}

func (b *bb) fieldArgSelector(a *rules.Arg) (x GO.ExprNode, leaf *gotype.StructField) {
	x = b.g.recv // default base
	s := b.g.info.KeyMap[a.Value].Selector

	if len(b.elems) > 0 {
		// resolve base expression
		for i := len(b.elems) - 1; i >= 0; i-- {
			if b.elems[i].n.Type.IsStructOrStructPointer() {
				x = b.elems[i].x
				break
			}
		}

		// resolve elem distance
		var distance int
		for i := len(s) - 1; i >= 0; i-- {
			if s[i].Type.Is(gotype.K_ARRAY, gotype.K_SLICE, gotype.K_MAP) {
				break
			}
			distance += 1
		}

		if len(s) > distance {
			s = s[(len(s) - distance):]
		}
	}

	for _, f := range s {
		if f.IsEmbedded {
			continue
		}

		x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
		leaf = f
	}
	return x, leaf
}

func (b *bb) constArg(n *rules.Node, r *rules.Rule, a *rules.Arg, t *gotype.Type) (x GO.ExprNode) {
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
