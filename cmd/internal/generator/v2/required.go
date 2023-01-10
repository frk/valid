package generate

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genReqPtrCode(o *types.Obj) {
	p := o // retain the pointer object

	g.P(`if `)
	o = g.genReqPtrExpr(o)
	g.L(` {
		return $0`, g.ErrExpr(p, p.IsRules[0]))

	switch {
	case len(o.IsRules) > 0:
		g.genIsRuleBlock(o, elif_block)

	case len(o.IsRules) == 0:
		g.L(`}`)
	}
}

func (g *generator) genReqPtrExpr(o *types.Obj) *types.Obj {
	g.P(`$0 `, o.IsRules[0])
	for o.Type.Kind == types.PTR {
		o = o.Type.Elem
		if o.Has(rules.REQUIRED) {
			g.P(`|| $0 `, o.IsRules[0])
		}
	}
	// "required" already done; remove
	if o.Has(rules.REQUIRED) {
		o.IsRules = o.IsRules[1:]
	}
	return o
}
