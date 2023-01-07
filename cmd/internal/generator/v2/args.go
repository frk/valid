package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genArg(a *rules.Arg) {
	i := g.info.ArgIndexMap[a]
	r := g.info.ArgRuleMap[a]
	o := g.info.RuleObjMap[r]
	tt := o.Type // target type

	switch r.Spec.Kind {
	case rules.LENGTH:
		tt = &types.Type{Kind: types.INT}
	case rules.FUNCTION, rules.PREPROC:
		fn := specs.GetFunc(r.Spec)
		if in := fn.Type.In; len(in)-1 <= i {
			tt = in[len(in)-1].Type
		} else {
			tt = in[i+1].Type
		}
	}

	if a.IsFieldRef() {
		g.genFieldArg(a, tt)
	} else {
		g.genConstArg(a, tt, r)
	}
}

func (g *generator) genFieldArg(a *rules.Arg, tt *types.Type) {
	ff := g.info.FRefChainMap[a]

	x := "v"
	for _, f := range ff {
		x = x + "." + f.Name
	}

	switch o := ff.Last().Obj; true {
	case o.Type.IsAssignableTo(tt) == false: // need conversion?
		g.P("$0($1)", o.Type, x)
	case tt.IsPtrOf(o.Type): // need address?
		g.P("&$0", x)
	default:
		g.S(x)
	}
}

func (g *generator) genConstArg(a *rules.Arg, tt *types.Type, r *rules.Rule) {
	switch {
	case tt.IsEmptyIface() && a.Type != rules.ARG_STRING:
		g.S(a.Value)
	case a.Type == rules.ARG_BOOL:
		g.S(a.Value)
	case a.Type == rules.ARG_INT:
		g.S(a.Value)
	case a.Type == rules.ARG_FLOAT:
		g.S(a.Value)
	case a.Type == rules.ARG_STRING && r.Spec.UseRawString:
		g.F("`%s`", a.Value)
	case a.Type == rules.ARG_STRING:
		g.F(`"%s"`, a.Value)
	case a.Type == rules.ARG_UNKNOWN:
		g.genZeroArg(tt)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genZeroArg(tt *types.Type) {
	switch tt.Kind {
	case types.STRING:
		g.S(`""`)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.S("0")
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.S("0")
	case types.FLOAT32, types.FLOAT64:
		g.S("0.0")
	case types.BOOL:
		g.S("false")
	case types.PTR, types.SLICE, types.MAP, types.INTERFACE:
		g.S("nil")
	default:
		panic("shouldn't reach")
	}
}
