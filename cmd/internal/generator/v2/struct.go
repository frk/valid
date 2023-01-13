package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genStructCode(o *types.Obj) {
	if !o.Type.HasRules() {
		return
	}

	if n := g.nptr(o); n > 1 {
		p := g.info.PtrMap[o]
		g.L(`$0 := *$1`, o, p)
	}

	for _, f := range o.Type.Fields {
		if !f.Obj.HasRules() {
			continue
		}

		g.genObjCode(f.Obj)
	}
}
