package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genOptPtrCode(o *types.Obj) {
	if true {
		// base has only one rule; and that rule is a basic rule (returns only bool / does not return error)
		//
		g.L(`if ${0[&&]} {`, p_rules(o))
		// _ = g.L(`	$1`, p_end(o))
		// _ = g.L(`}`)
		g.genObjBlock(p_end(o), sub_block)
	} else {
		g.P(`if `)
		o = g.genOptPtrExpr(o)
		g.genObjBlock(o, sub_block)
	}
}

func (g *generator) genOptPtrExpr(o *types.Obj) *types.Obj {
	g.P(`$0 `, o.IsRules[0])
	for o.Type.Kind == types.PTR {
		o = o.Type.Elem
		if o.Has(rules.OPTIONAL) {
			g.P(`&& $0 `, o.IsRules[0])
		}
	}
	// "optional" already done; remove
	if o.Has(rules.OPTIONAL) {
		o.IsRules = o.IsRules[1:]
	}

	if len(o.PreRules) == 0 && len(o.IsRules) == 1 {
		if r := o.IsRules[0]; !specs.CanReturnError(r.Spec) {
			g.P(` && `)
			if g.isMultiExprRule(o, r) {
				g.P(`(`)
				g.genIsRuleExpr(o, r)
				g.P(`)`)
			} else {
				g.genIsRuleExpr(o, r)
			}
			g.L(` {
				return $0
			}`, g.ErrExpr(o, r))

			o.IsRules = o.IsRules[1:]
		}
	}
	return o
}
