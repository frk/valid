package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genPreRuleCode(o *types.Obj) {
	g.P(`$0 = `, o)
	for i := len(o.PreRules) - 1; i >= 0; i-- {
		fn := specs.GetFunc(o.PreRules[i].Spec)
		g.P(`$0(`, fn)
	}
	g.P(`$0`, o)
	for _, r := range o.PreRules {
		for _, a := range r.Args {
			g.P(`, $0`, a)
		}
		g.S(`)`)
	}
	g.L(``)

	o.PreRules = nil
	g.genObjBlock(o, current_block)
}
