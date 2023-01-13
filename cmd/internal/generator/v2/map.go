package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) gen_map_code(o *types.Obj) {
	k, v := o.Type.Key, o.Type.Elem
	if !k.HasRules() && !v.HasRules() {
		return
	}

	// TODO handle objects with only pre-rules
	// that are not pointers, and other variations

	switch {
	case k.HasRules() && !v.HasRules():
		g.L(`for ${0} := range ${1} {`, k, o)
		g.P(`	${0:g}`, k)
		g.L(`}`)

	case k.HasRules() && v.HasRules():
		g.L(`for ${0}, ${1} := range ${2} {`, k, v, o)
		g.P(`	${0:g}`, k)
		g.P(`	${0:g}`, v)
		g.L(`}`)

	case !k.HasRules() && v.HasRules():
		g.L(`for _, ${0} := range ${1} {`, v, o)
		g.P(`	${0:g}`, v)
		g.L(`}`)
	}
}
