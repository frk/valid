package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) gen_slice_code(o *types.Obj) {
	e := o.Type.Elem
	if !e.HasRules() {
		return
	}

	// declare index if value has preprocs but is not a pointer
	with_index := (len(e.PreRules) > 0 && e.Type.Kind != types.PTR)
	// declare elem if it has rules other than the preprocs
	with_elem := len(e.IsRules) > 0 || e.Type.HasRules()

	switch {
	case with_index && !with_elem:
		g.vars[e] = g.vars[o] + "[i]"
		g.L(`for i := range ${0} {`, o)
		g.P(`	${0} = ${1[@]}`, e, e.PreRules)
		g.L(`}`)

	case with_index && with_elem:
		pre := e.PreRules
		e.PreRules = nil

		// TODO: am not happy with how the vars need
		// to be handled need to figure out a better
		// way to do this, to make code compact.
		v := g.vars[e]
		g.L(`for i, ${0} := range ${1} {`, e, o)
		g.vars[e] = g.vars[o] + "[i]"
		g.L(`	${0} = ${1[@]}`, e, pre)
		g.vars[e] = v
		g.P(`	${0:g}`, e)
		g.L(`}`)

	case !with_index && with_elem:
		g.L(`for _, ${0} := range ${1} {`, e, o)
		g.P(`	${0:g}`, e)
		g.L(`}`)
	}
}
