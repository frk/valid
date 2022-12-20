package generate

import (
	"github.com/frk/valid/cmd/internal/rules/spec"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genPtrCode(f *types.StructField, o *types.Obj) *types.Obj {
	if o.Type.Kind != types.PTR {
		return o
	}

	switch {
	case g.hasRequiredRule(o):
		g.wr.p(`if `)
		for o.Type.Kind == types.PTR {
			rs := g.info.RuleMap[o]
			r := rs.Is.Rules[0]
			s := g.info.SpecMap[r]

			o = o.Type.Elem
		}
		g.wr.ln(`}`)

	case g.hasOptionalRule(o):
		g.wr.p(`if `)
		for o.Type.Kind == types.PTR {
			// TODO
			o = o.Type.Elem
		}
		g.wr.ln(`}`)

	case g.hasNoguardRule(o):
	}

	return o
}
