package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genSliceBlock(o *types.Obj, block blockType) {
	E := o.Type.Elem
	if !E.HasRules() {
		return
	}

	switch block {
	default:
		g.genSliceCode(o, E)
	case else_block:
		g.RL(`} else {`)
		g.genSliceCode(o, E)
		g.L(`}`)
	case sub_block:
		g.L(`{`)
		g.genSliceCode(o, E)
		g.L(`}`)
	}
}

func (g *generator) genSliceCode(o, E *types.Obj) {
	g.L(`for _, $0 := range $1 {`, E, o)
	g.genObjCode(E)
	g.L(`}`)
}
