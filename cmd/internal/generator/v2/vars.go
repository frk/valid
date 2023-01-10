package generate

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) prepVars() {
	t := g.info.Validator.Type
	g.v = &types.Obj{Type: t}

	g.vars = make(map[*types.Obj]string)
	g.vars[g.v] = "v"

	g.prepObjVars(g.v)
}

func (g *generator) prepObjVars(o *types.Obj) {
	switch t := o.Type; t.Kind {
	case types.PTR:
		g.vars[t.Elem] = "*" + g.vars[o]
		g.prepObjVars(t.Elem)

	case types.ARRAY, types.SLICE:
		if g.vars[o] != "e" {
			g.vars[t.Elem] = "e"
		} else {
			g.vars[t.Elem] = "e2"
		}
		g.prepObjVars(t.Elem)

	case types.MAP:
		if g.vars[o] != "k" {
			g.vars[t.Key] = "k"
		} else {
			g.vars[t.Key] = "k2"
		}
		g.prepObjVars(t.Key)

		if g.vars[o] != "e" {
			g.vars[t.Elem] = "e"
		} else {
			g.vars[t.Elem] = "e2"
		}
		g.prepObjVars(t.Elem)

	case types.STRUCT:
		v := o
		if n := g.nptr(o); n == 1 {
			v = g.info.PtrMap[o]
		} else if n > 1 {
			// during the second pass, in this scenario, an assignment
			// statement "f := <expr>" will be produced of the generator.
			g.vars[v] = "f"
		}

		for _, f := range t.VisibleFields() {
			g.vars[f.Obj] = g.vars[v] + "." + f.Name
			g.prepObjVars(f.Obj)
		}
	}
}
