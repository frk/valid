package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genMapBlock(o *types.Obj, block blockType) {
	K, V := o.Type.Key, o.Type.Elem
	if !K.HasRules() && !V.HasRules() {
		return
	}

	switch block {
	default:
		g.genMapCode(o, K, V)
	case else_block:
		g.RL(`} else {`)
		g.genMapCode(o, K, V)
		g.L(`}`)
	case sub_block:
		g.L(`{`)
		g.genMapCode(o, K, V)
		g.L(`}`)
	}
}

func (g *generator) genMapCode(o, K, V *types.Obj) {
	var (
		key string
		val string
	)

	if V.HasRules() {
		val = g.vars[V]
	}

	if K.HasRules() || (V.Type.Kind != types.PTR && len(V.PreRules) > 0) {
		key = g.vars[K]
	} else if val != "" {
		key = "_"
	}

	if val == "" {
		g.L(`for $0 := range $1 {`, key, o)
		g.genObjCode(K)
		g.L(`}`)
	} else {
		g.L(`for $0, $1 := range $2 {`, key, val, o)
		g.genObjCode(K)
		g.genObjCode(V)
		g.L(`}`)
	}
}
