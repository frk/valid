package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genSliceBlock(f *types.StructField, t *types.Type, block blockType) {
	if !t.Elem.HasRules() {
		return
	}

	o := t.Elem
	switch block {
	default:
		g.genSliceCode(f, o)
	case else_block:
		g.RL(`} else {`)
		g.genSliceCode(f, o)
		g.L(`}`)
	case sub_block:
		g.L(`{`)
		g.genSliceCode(f, o)
		g.L(`}`)
	}
}

func (g *generator) genSliceCode(f *types.StructField, o *types.Obj) {
	x := g.vars["x"]
	g.L(`for _, e := range $x {`)
	g.vars["x"] = "e"
	g.genObjCode(f, o)
	g.L(`}`)
	g.vars["x"] = x
}
