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

func (g *generator) gen_slice_code(o *types.Obj) {
	e := o.Elem
	pre := e.PreRules
	e.PreRules = nil

	more_rules := e.HasRules()
	needs_index := (len(pre) > 0 && e.Type.Kind != types.PTR)

	switch {
	case needs_index && !more_rules:
	case needs_index && more_rules:
		g.L(`for i, $0 := range $1 {`, e, o)
		g.P(`	${0:g}`, e)
		g.L(`}`)

	case more_rules:
		g.L(`for _, $0 := range $1 {`, e, o)
		g.P(`	${0:g}`, e)
		g.L(`}`)
	}
}
