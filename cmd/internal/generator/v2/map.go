package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genMapBlock(f *types.StructField, t *types.Type, block blockType) {
	K, V := t.Key, t.Elem
	if !K.HasRules() && !V.HasRules() {
		return
	}

	switch block {
	default:
		g.genMapCode(f, K, V)
	case else_block:
		g.RL(`} else {`)
		g.genMapCode(f, K, V)
		g.L(`}`)
	case sub_block:
		g.L(`{`)
		g.genMapCode(f, K, V)
		g.L(`}`)
	}
}

func (g *generator) genMapCode(f *types.StructField, K, V *types.Obj) {
	x := g.vars["x"]

	var v string
	if V.HasRules() {
		v = "e"
	}

	var k string
	if K.HasRules() || (V.Type.Kind != types.PTR && len(V.PreRules) > 0) {
		k = "k"
	} else if v != "" {
		k = "_"
	}

	if x == k {
		k += "2"
	}
	if x == v {
		v += "2"
	}

	if v == "" {
		g.L(`for $0 := range $x {`, k)
		g.vars["x"] = k
		g.genObjCode(f, K)
		g.L(`}`)
	} else {
		g.L(`for $0, $1 := range $x {`, k, v)
		g.vars["x"] = k
		g.genObjCode(f, K)
		g.vars["x"] = v
		g.genObjCode(f, V)
		g.L(`}`)
	}

	g.vars["x"] = x
}
