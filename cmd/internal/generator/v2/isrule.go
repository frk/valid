package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

type blockType uint

const (
	no_block blockType = iota
	if_block
	else_block
	elif_block
	sub_block
)

func (g *generator) genIsRuleBlock(f *types.StructField, o *types.Obj, block blockType) {
	switch block {
	case if_block:
		g.P(`if `)
	case elif_block:
		g.P(`} else if `)
	}

	if o.Has(rules.OPTIONAL) && len(o.IsRules) > 1 {
		g.genCondExpr(f, o, o.IsRules[0])
		g.P(` && `)
		o.IsRules = o.IsRules[1:] //:F
	}

	r := o.IsRules[0]
	g.genCondExpr(f, o, r)
	g.L(` {`)
	g.genError(f, o, r)
	for _, r := range o.IsRules[1:] {
		g.P(`} else if `)
		g.genCondExpr(f, o, r)
		g.L(` {`)
		g.genError(f, o, r)
	}
	g.L(`}`)
}

func (g *generator) genCondExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.genIsRuleExpr(f, o, r)
}

func (g *generator) genIsRuleExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch r.Spec.Kind {
	case rules.REQUIRED:
		switch r.Name {
		case "notnil":
			g.genNotNilExpr(f, o, r)
		case "required":
			g.genRequiredExpr(f, o, r)
		}
	case rules.OPTIONAL:
		switch r.Name {
		case "omitnil":
			g.genOmitNilExpr(f, o, r)
		case "optional":
			g.genOptionalExpr(f, o, r)
		}
	case rules.COMPARABLE:
		switch r.Name {
		case "eq":
			g.genEqualExpr(f, o, r)
		case "ne":
			g.genNotEqualExpr(f, o, r)
		}
	case rules.ORDERED:
		switch r.Name {
		case "gt":
			g.genGreaterThanExpr(f, o, r)
		case "lt":
			g.genLessThanExpr(f, o, r)
		case "gte", "min":
			g.genGreaterThanOrEqualExpr(f, o, r)
		case "lte", "max":
			g.genLessThanOrEqualExpr(f, o, r)
		}
	case rules.LENGTH:
		switch r.Name {
		case "len":
			g.genLenExpr(f, o, r)
		case "runecount":
			g.genRuneCountExpr(f, o, r)
		}
	case rules.RANGE:
		g.genRangeExpr(f, o, r)
	case rules.ENUM:
		g.genEnumExpr(f, o, r)
	case rules.FUNCTION:
		g.genFuncCallExpr(f, o, r)
	case rules.METHOD:
		g.genMethodExpr(f, o, r)
	}
}

func (g *generator) genOmitNilExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.MAP, types.SLICE, types.INTERFACE:
		g.P(`$x != nil`)
	case types.PTR, types.FUNC, types.CHAN:
		g.P(`$x != nil`)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genOptionalExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.STRING:
		g.P(`$x != ""`)
	case types.MAP, types.SLICE:
		g.P(`len($x) > 0`)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.P(`$x > 0`)
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.P(`$x > 0`)
	case types.FLOAT32, types.FLOAT64:
		g.P(`$x > 0.0`)
	case types.BOOL:
		g.P(`$x == true`)
	case types.INTERFACE:
		g.P(`$x != nil`)
	case types.PTR:
		g.P(`$x != nil`)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genNotNilExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.SLICE, types.MAP, types.INTERFACE:
		g.P(`$x == nil`)
	case types.PTR, types.CHAN, types.FUNC:
		g.P(`$x == nil`)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genRequiredExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch o.Type.Kind {
	case types.STRING:
		g.P(`$x == ""`)
	case types.SLICE, types.MAP:
		g.P(`len($x) == 0`)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.P(`$x == 0`)
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.P(`$x == 0`)
	case types.FLOAT32, types.FLOAT64:
		g.P(`$x == 0.0`)
	case types.BOOL:
		g.P(`$x == false`)
	case types.INTERFACE:
		g.P(`$x == nil`)
	case types.PTR:
		g.P(`$x == nil`)
	case types.STRUCT:
		if o.Type.Pkg == g.file.pkg {
			g.P(`$x == ($0{})`, o.Type.Name)
		} else {
			pkg := g.file.addImport(o.Type.Pkg)
			g.P(`$x == ($0.$1{})`, pkg.name, o.Type.Name)
		}
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genEqualExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x != $0 `, r.Args[0])
	for _, a := range r.Args[1:] {
		g.P(`&& $x != $0 `, a)
	}
}

func (g *generator) genNotEqualExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x == $0 `, r.Args[0])
	for _, a := range r.Args[1:] {
		g.P(`|| $x == $0 `, a)
	}
}

func (g *generator) genGreaterThanExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x <= $0 `, r.Args[0])
}

func (g *generator) genLessThanExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x >= $0 `, r.Args[0])
}

func (g *generator) genGreaterThanOrEqualExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x < $0 `, r.Args[0])
}

func (g *generator) genLessThanOrEqualExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x > $0 `, r.Args[0])
}

func (g *generator) genLenExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	switch {
	case len(r.Args) == 1:
		g.P(`len($x) != $0 `, r.Args[0])
	case len(r.Args) == 2 && r.Args[1].Value == "":
		g.P(`len($x) < $0 `, r.Args[0])
	case len(r.Args) == 2 && r.Args[0].Value == "":
		g.P(`len($x) > $0 `, r.Args[1])
	case len(r.Args) == 2:
		g.P(`len($x) < $0 || len($x) > $1 `, r.Args[0], r.Args[1])
	default:
		panic("shouldn't happen")
	}
}

func (g *generator) genRuneCountExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	pn := g.file.addImport(types.Pkg{Path: "unicode/utf8", Name: "utf8"}).name
	fn := ""

	switch {
	case o.Type.Kind == types.STRING:
		fn = "RuneCountInString"
	case o.Type.Kind == types.SLICE && o.Type.Elem.Type.IsByte:
		fn = "RuneCount"
	default:
		panic("shouldn't happen")
	}

	switch {
	case len(r.Args) == 1:
		g.P(`$0.$1($x) != $2 `, pn, fn, r.Args[0])
	case len(r.Args) == 2 && r.Args[1].Value == "":
		g.P(`$0.$1($x) < $2 `, pn, fn, r.Args[0])
	case len(r.Args) == 2 && r.Args[0].Value == "":
		g.P(`$0.$1($x) > $2 `, pn, fn, r.Args[1])
	case len(r.Args) == 2:
		g.P(`$0.$1($x) < $2 || $0.$1($x) > $3 `, pn, fn, r.Args[0], r.Args[1])
	default:
		panic("shouldn't happen")
	}
}

func (g *generator) genRangeExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	g.P(`$x < $0 || $x > $1 `, r.Args[0], r.Args[1])
}

func (g *generator) genEnumExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	enums := g.info.EnumMap[o.Type]
	g.P(`$x != $0 `, enums[0])
	for _, e := range enums[1:] {
		g.P(`&& $x != $0 `, e)
	}
}

// builds an expression that checks the value using the rule spec's function.
func (g *generator) genFuncCallExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	fn := specs.GetFunc(r.Spec)

	// If this is the included regexp rule, then add
	// a registry call statement for the init function.
	if r.Spec.Name == "re" && fn.Type.IsIncluded() {
		g.file.addImport(fn.Type.Pkg)
		g.file.addRegExp(r.Args[0])
	}

	// If the object's type isn't identical to the target
	// check if it needs to be converted or not.
	x := func() { g.P("$x") }
	tt := fn.Type.In[0].Type
	if fn.Type.IsVariadic && len(fn.Type.In) == 1 {
		tt = tt.Elem.Type
	}
	if !o.Type.IsAssignableTo(tt) { // need conversion?
		x = func() { g.P("$0($x)", tt) }
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

func (g *generator) genMethodExpr(f *types.StructField, o *types.Obj, r *rules.Rule) {
	if g.nptr > 1 || (g.nptr == 1 && o.Type.Kind == types.INTERFACE) {
		g.P(`!($x).$0()`, r.Spec.Func.Name)
	} else {
		g.P(`!$x.$0()`, r.Spec.Func.Name)
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
