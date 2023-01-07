package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genStructBlock(t *types.Type, block blockType) {
	switch block {
	case sub_block:
		g.L(`{`)
		g.genStructCode(t)
		g.L(`}`)
	case else_block:
		g.RL(`} else {`)
		g.genStructCode(t)
		g.L(`}`)
	default:
		g.genStructCode(t)
	}
}

func (g *generator) genStructCode(t *types.Type) {
	x := g.vars["x"]
	for _, f := range t.Fields {
		if !f.Obj.HasRules() {
			continue
		}

		g.fsel(x, f)
		// if x is not a pointer, then we can do a simple x.f
		// if x is single pointer, then we can also do a simple x.f
		// if x is multi pointer, then we need to do (x).f; or f := x; f.f
		//
		// TODO: how to represent x such that:
		// - it generates the correct expression
		// - it can be easily set/reset to the desired expression
		// - it can be used not only for struct fields, but also for rules...
		//
		// ... basically if we need to access a member (field/method) of
		// the current expression, the expression must be at MOST one pointer
		// deep if it is a non-interface type, and if it is an interface
		// type then it cannot be a pointer at all...
		g.genObjCode(f, f.Obj)
	}
	g.vars["x"] = x
}
