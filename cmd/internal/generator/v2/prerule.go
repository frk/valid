package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genPreRuleBlock(f *types.StructField, o *types.Obj, block blockType) {
	switch block {
	case sub_block:
		g.L(`{`)
		g.genPreRuleCode(f, o)
		g.L(`}`)

	default:
		g.genPreRuleCode(f, o)
	}
}

func (g *generator) genPreRuleCode(f *types.StructField, o *types.Obj) {
	g.P(`$x = `)
	for i := len(o.PreRules) - 1; i >= 0; i-- {
		fn := specs.GetFunc(o.PreRules[i].Spec)
		g.P(`$0(`, fn)
	}
	g.P(`$x`)
	for _, r := range o.PreRules {
		for _, a := range r.Args {
			g.P(`, $0`, a)
		}
		g.S(`)`)
	}
	g.L(``)
}
