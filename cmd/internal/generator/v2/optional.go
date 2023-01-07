package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genOptCode(f *types.StructField, o *types.Obj) (x *types.Obj) {
	x = o
	g.P(`if $0 `, x.IsRules[0])
	for x.Type.Kind == types.PTR {
		x = x.Type.Elem
		g.deref()

		if x.Has(rules.OPTIONAL) {
			g.P(`&& $0 `, x.IsRules[0])
		}
	}

	// "optional" already done; remove
	if x.Has(rules.OPTIONAL) {
		x.IsRules = x.IsRules[1:]
	}

	switch {
	case len(x.IsRules) == 1 && !specs.CanReturnError(x.IsRules[0].Spec):
		g.P(` && `)
		if g.isMultiExprRule(x, x.IsRules[0]) {
			g.P(`(`)
			g.genIsRuleExpr(f, x, x.IsRules[0])
			g.P(`)`)
		} else {
			g.genIsRuleExpr(f, x, x.IsRules[0])
		}
		g.L(` {`)
		g.genError(f, x, x.IsRules[0])
		g.L(`}`)

	case len(x.IsRules) == 1 && specs.CanReturnError(x.IsRules[0].Spec):
		g.L(`{`)
		g.genIsRuleBlock(f, x, if_block)
		g.L(`}`)

	case len(x.IsRules) > 1:
		g.L(`{`)
		g.genIsRuleBlock(f, x, if_block)
		g.L(`}`)
	}

	x.IsRules = nil // :F
	return x
}
