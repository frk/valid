package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genIsRuleExpr(o *types.Obj, r *rules.Rule) {
	switch r.Spec.Kind {
	case rules.REQUIRED:
		switch r.Name {
		case "notnil":
			g.genNotNilExpr(o, r)
		case "required":
			g.genRequiredExpr(o, r)
		}
	case rules.OPTIONAL:
		switch r.Name {
		case "omitnil":
			g.genOmitNilExpr(o, r)
		case "optional":
			g.genOptionalExpr(o, r)
		}
	case rules.COMPARABLE:
		switch r.Name {
		case "eq":
			g.genEqualExpr(o, r)
		case "ne":
			g.genNotEqualExpr(o, r)
		}
	case rules.ORDERED:
		switch r.Name {
		case "gt":
			g.genGreaterThanExpr(o, r)
		case "lt":
			g.genLessThanExpr(o, r)
		case "gte", "min":
			g.genGreaterThanOrEqualExpr(o, r)
		case "lte", "max":
			g.genLessThanOrEqualExpr(o, r)
		}
	case rules.LENGTH:
		switch r.Name {
		case "len":
			g.genLenExpr(o, r)
		case "runecount":
			g.genRuneCountExpr(o, r)
		}
	case rules.RANGE:
		g.genRangeExpr(o, r)
	case rules.ENUM:
		g.genEnumExpr(o, r)
	case rules.FUNCTION:
		g.genFuncCallExpr(o, r)
	case rules.METHOD:
		g.genMethodExpr(o, r)
	}
}

func (g *generator) genOmitNilExpr(o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.MAP, types.SLICE, types.INTERFACE:
		g.P(`$0 != nil`, o)
	case types.PTR, types.FUNC, types.CHAN:
		g.P(`$0 != nil`, o)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genOptionalExpr(o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.STRING:
		g.P(`$0 != ""`, o)
	case types.MAP, types.SLICE:
		g.P(`len($0) > 0`, o)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.P(`$0 > 0`, o)
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.P(`$0 > 0`, o)
	case types.FLOAT32, types.FLOAT64:
		g.P(`$0 > 0.0`, o)
	case types.BOOL:
		g.P(`$0 == true`, o)
	case types.INTERFACE:
		g.P(`$0 != nil`, o)
	case types.PTR:
		g.P(`$0 != nil`, o)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genNotNilExpr(o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.SLICE, types.MAP, types.INTERFACE:
		g.P(`$0 == nil`, o)
	case types.PTR, types.CHAN, types.FUNC:
		g.P(`$0 == nil`, o)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genRequiredExpr(o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.STRING:
		g.P(`$0 == ""`, o)
	case types.SLICE, types.MAP:
		g.P(`len($0) == 0`, o)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.P(`$0 == 0`, o)
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.P(`$0 == 0`, o)
	case types.FLOAT32, types.FLOAT64:
		g.P(`$0 == 0.0`, o)
	case types.BOOL:
		g.P(`$0 == false`, o)
	case types.INTERFACE:
		g.P(`$0 == nil`, o)
	case types.PTR:
		g.P(`$0 == nil`, o)
	case types.STRUCT:
		if o.Type.Pkg == g.file.pkg {
			g.P(`$0 == ($1{})`, o, o.Type.Name)
		} else {
			pkg := g.file.addImport(o.Type.Pkg)
			g.P(`$0 == ($1.$2{})`, o, pkg.name, o.Type.Name)
		}
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genEqualExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 != $1 `, o, r.Args[0])
	for _, a := range r.Args[1:] {
		g.P(`&& $0 != $1 `, o, a)
	}
}

func (g *generator) genNotEqualExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 == $1 `, o, r.Args[0])
	for _, a := range r.Args[1:] {
		g.P(`|| $0 == $1 `, o, a)
	}
}

func (g *generator) genGreaterThanExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 <= $1 `, o, r.Args[0])
}

func (g *generator) genLessThanExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 >= $1 `, o, r.Args[0])
}

func (g *generator) genGreaterThanOrEqualExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 < $1 `, o, r.Args[0])
}

func (g *generator) genLessThanOrEqualExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 > $1 `, o, r.Args[0])
}

func (g *generator) genLenExpr(o *types.Obj, r *rules.Rule) {
	switch {
	case len(r.Args) == 1:
		g.P(`len($0) != $1 `, o, r.Args[0])
	case len(r.Args) == 2 && r.Args[1].Value == "":
		g.P(`len($0) < $1 `, o, r.Args[0])
	case len(r.Args) == 2 && r.Args[0].Value == "":
		g.P(`len($0) > $1 `, o, r.Args[1])
	case len(r.Args) == 2:
		g.P(`len($0) < $1 || len($0) > $2 `, o, r.Args[0], r.Args[1])
	default:
		panic("shouldn't happen")
	}
}

func (g *generator) genRuneCountExpr(o *types.Obj, r *rules.Rule) {
	var fn *types.Func
	switch {
	case o.Type.Kind == types.STRING:
		p := types.Pkg{Path: "unicode/utf8", Name: "utf8"}
		fn = &types.Func{Name: "RuneCountInString", Type: &types.Type{Pkg: p}}
	case o.Type.Kind == types.SLICE && o.Type.Elem.Type.IsByte:
		p := types.Pkg{Path: "unicode/utf8", Name: "utf8"}
		fn = &types.Func{Name: "RuneCount", Type: &types.Type{Pkg: p}}
	default:
		panic("shouldn't happen")
	}

	switch {
	case len(r.Args) == 1:
		g.P(`$0($1) != $2 `, fn, o, r.Args[0])
	case len(r.Args) == 2 && r.Args[1].Value == "":
		g.P(`$0($1) < $2 `, fn, o, r.Args[0])
	case len(r.Args) == 2 && r.Args[0].Value == "":
		g.P(`$0($1) > $2 `, fn, o, r.Args[1])
	case len(r.Args) == 2:
		g.P(`$0($1) < $2 || $0($1) > $3 `, fn, o, r.Args[0], r.Args[1])
	default:
		panic("shouldn't happen")
	}
}

func (g *generator) genRangeExpr(o *types.Obj, r *rules.Rule) {
	g.P(`$0 < $1 || $0 > $2 `, o, r.Args[0], r.Args[1])
}

func (g *generator) genEnumExpr(o *types.Obj, r *rules.Rule) {
	enums := g.info.EnumMap[o.Type]
	g.P(`$0 != $1 `, o, enums[0])
	for _, e := range enums[1:] {
		g.P(`&& $0 != $1 `, o, e)
	}
}

// builds an expression that checks the value using the rule spec's function.
func (g *generator) genFuncCallExpr(o *types.Obj, r *rules.Rule) {
	fn := specs.GetFunc(r.Spec)

	// If this is the included regexp rule, then add
	// a registry call statement for the init function.
	if r.Spec.Name == "re" && fn.Type.IsIncluded() {
		g.file.addImport(fn.Type.Pkg)
		g.file.addRegExp(r.Args[0])
	}

	// If the object's type isn't identical to the target
	// check if it needs to be converted or not.
	x := func() { g.P("$0", o) }
	tt := fn.Type.In[0].Type
	if fn.Type.IsVariadic && len(fn.Type.In) == 1 {
		tt = tt.Elem.Type
	}
	if !o.Type.IsAssignableTo(tt) { // need conversion?
		x = func() { g.P("$0($1)", tt, o) }
	}

	switch {
	case fn.Type.CanError():
		g.P(`ok, err := $0($1`, fn, x)
		for _, a := range r.Args {
			g.P(`, $0`, a)
		}
		g.L(`); err != nil {`)
		g.L(`return err`)
		g.P(`} else if !ok `)

	default:
		switch {
		case len(r.Args) == 0:
			g.P(`!$0($1) `, fn, x)
		case len(r.Args) > 0 && r.Spec.JoinOp == 0:
			g.P(`!$0($1`, fn, x)
			for _, a := range r.Args {
				g.P(`, $0`, a)
			}
			g.P(`) `)
		case len(r.Args) > 0 && r.Spec.JoinOp > 0:
			if r.Spec.JoinOp == rules.JOIN_NOT {
				g.P(`$0($1, $2) `, fn, x, r.Args[0])
			} else {
				g.P(`!$0($1, $2) `, fn, x, r.Args[0])
			}
			for _, a := range r.Args[1:] {
				switch r.Spec.JoinOp {
				case rules.JOIN_NOT:
					g.P(`|| $0($1, $2) `, fn, x, a)
				case rules.JOIN_AND:
					g.P(`|| !$0($1, $2) `, fn, x, a)
				case rules.JOIN_OR:
					g.P(`&& !$0($1, $2) `, fn, x, a)
				}
			}
		}
	}
}

func (g *generator) genMethodExpr(o *types.Obj, r *rules.Rule) {
	n := g.nptr(o)
	if n > 1 || (n == 1 && o.Type.Kind == types.INTERFACE) {
		g.P(`!($0).$1()`, o, r.Spec.Func.Name)
	} else {
		g.P(`!$0.$1()`, o, r.Spec.Func.Name)
	}
	// TODO?
	// if n.PtrDepth() == 1 && n.Type.Kind != xtypes.K_INTERFACE {
	// 	x = b.vals[len(b.vals)-1]
	// }
}

////////////////////////////////////////////////////////////////////////////////
// helper

func (g *generator) isMultiExprRule(o *types.Obj, r *rules.Rule) bool {
	switch r.Spec.Kind {
	case rules.COMPARABLE:
		return len(r.Args) > 1
	case rules.LENGTH:
		return len(r.Args) == 2 && r.Args[0].Value != "" && r.Args[1].Value != ""
	case rules.RANGE:
		return true
	case rules.ENUM:
		return len(g.info.EnumMap[o.Type]) > 1
	case rules.FUNCTION:
		return r.Spec.JoinOp > 0 && len(r.Args) > 1
	}
	return false
}

func (g *generator) nptr(o *types.Obj) (depth int) {
	for {
		p, ok := g.info.PtrMap[o]
		if !ok || o == nil {
			break
		}

		o = p
		depth += 1
	}
	return depth
}
