package generate

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) gen_obj_field_key(o *types.Obj) {
	f := g.info.ObjFieldMap[o]
	k := g.info.FKeyMap[f]
	g.Q(k)
}

type blockType uint

const (
	current_block blockType = iota
	sub_block
	if_block
	else_block
	elif_block
)

func (g *generator) genObjBlock(o *types.Obj, block blockType) {
	if !o.HasRules() {
		return
	}

	switch block {
	case current_block:
		g.genObjCode(o)

	case sub_block:
		g.L(`{`)
		g.genObjCode(o)
		g.L(`}`)

	case else_block:
		g.L(`} else {`)
		g.genObjCode(o)
		//g.L(`}`)

	case elif_block:
		g.P(`} else if `)
		g.genObjCode(o)
		//g.L(`}`)
	}
}

func (g *generator) genObjCode(o *types.Obj) {
	switch {
	case o.Type.Kind == types.PTR && o.Has(rules.REQUIRED):
		g.gen_required_pointer_code(o)

	case o.Type.Kind == types.PTR && o.Has(rules.OPTIONAL):
		g.gen_optional_pointer_code(o)

	case o.Type.Kind == types.PTR:
		for o.Type.Kind == types.PTR {
			o = o.Type.Elem
		}
		g.genObjCode(o)

	case len(o.PreRules) > 0:
		g.gen_pre_rules_code(o)

	case len(o.IsRules) > 0:
		g.gen_is_rules_code(o)

	case o.Type.Kind == types.MAP:
		g.gen_map_code(o)
	case o.Type.Kind == types.ARRAY:
		g.gen_slice_code(o)
	case o.Type.Kind == types.SLICE:
		g.gen_slice_code(o)
	case o.Type.Kind == types.STRUCT:
		g.genStructCode(o)
	}
}
