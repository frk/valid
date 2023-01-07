package generate

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genReqCode(f *types.StructField, o *types.Obj) (x *types.Obj) {
	x = o
	g.P(`if $0 `, x.IsRules[0])
	for x.Type.Kind == types.PTR {
		x = x.Type.Elem
		g.deref()

		if x.Has(rules.REQUIRED) {
			g.P(`|| $0 `, x.IsRules[0])
		}
	}

	// "required" already done; remove
	if x.Has(rules.REQUIRED) {
		x.IsRules = x.IsRules[1:]
	}

	g.L(`{`)
	g.genError(f, o, o.IsRules[0])
	if len(x.IsRules) > 0 {
		g.genIsRuleBlock(f, x, elif_block)
		x.IsRules = nil // already done; remove
	} else {
		g.L(`}`)
	}

	return x
}
