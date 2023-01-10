package generate

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
)

func (g *generator) genRulesExpr(rr []*rules.Rule, op exprOp) {
	switch rr[0].Spec.Kind {
	case rules.PREPROC:
		// TODO... if rules preproc...

	default:
		if op.has(unary_not) {
			g.P(`!${0} `, rr[0])
		} else {
			g.P(`${0} `, rr[0])
		}
		for _, r := range rr[1:] {
			switch {
			case op.has(unary_not, bool_and):
				g.P(`&& !${0}`, r)
			case op.has(unary_not, bool_or):
				g.P(`|| !${0}`, r)
			case op.has(bool_and):
				g.P(`&& ${0}`, r)
			case op.has(bool_or):
				g.P(`|| ${0}`, r)
			}
		}
	}
}

func (g *generator) genRuleExpr(r *rules.Rule) {
	// TODO
	o := g.info.RuleObjMap[r]
	g.genIsRuleExpr(o, r)
}
