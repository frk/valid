package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) gen_arg_list_expr(args []*rules.Arg, t_any bool) {
	switch t_any {
	case false:
		g.P(`${0}`, args[0])
		for _, a := range args[1:] {
			g.P(`, ${0}`, a)
		}

	case true:
		g.P(`${0:any}`, args[0])
		for _, a := range args[1:] {
			g.P(`, ${0:any}`, a)
		}
	}
}

func (g *generator) gen_arg_expr(a *rules.Arg, t_any bool) {
	var tt *types.Type // target type
	r := g.info.ArgRuleMap[a]

	switch t_any {
	case false:
		i := g.info.ArgIndexMap[a]
		o := g.info.RuleObjMap[r]
		tt = o.Type

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

	case true:
		tt = &types.Type{Kind: types.INTERFACE}
	}

	if a.IsFieldRef() {
		g.gen_field_arg(a, tt)
	} else {
		g.gen_const_arg(a, tt, r)
	}
}

func (g *generator) gen_field_arg(a *rules.Arg, tt *types.Type) {
	f := g.info.FRefMap[a]
	switch o := f.Obj; true {
	case o.Type.IsAssignableTo(tt) == false: // need conversion?
		g.P("${0}(${1})", tt, o)
	case tt.IsPtrOf(o.Type): // need address?
		g.P("&${0}", o)
	default:
		g.P("${0}", o)
	}
}

func (g *generator) gen_const_arg(a *rules.Arg, tt *types.Type, r *rules.Rule) {
	switch {
	case tt.IsEmptyIface() && a.Type == rules.ARG_UNKNOWN:
		g.S(`""`)
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
