package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genStructBlock(o *types.Obj, block blockType) {
	if !o.Type.HasRules() {
		return
	}

	switch block {
	case current_block:
		g.genStructCode(o)

	case sub_block:
		g.L(`{`)
		g.genStructCode(o)
		g.L(`}`)

	case else_block:
		g.RL(`} else {`)
		g.genStructCode(o)
		g.L(`}`)
	}
}

func (g *generator) genStructCode(o *types.Obj) {
	if n := g.nptr(o); n > 1 {
		p := g.info.PtrMap[o]
		g.L(`$0 := *$1`, o, p)
	}

	for _, f := range o.Type.Fields {
		if !f.Obj.HasRules() {
			continue
		}

		g.genObjBlock(f.Obj, current_block)
	}
}
