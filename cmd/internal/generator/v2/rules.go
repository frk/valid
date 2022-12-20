package generate

import (
	"github.com/frk/valid/cmd/internal/rules/spec"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

//func (g *generator) ifStmt(rs *rules.Set) {
//	g.wr.p(`if `)
//	for _, r := range rs.Is.Rules {
//		s := g.info.SpecMap[r]
//		g.ruleCond(r, s)
//
//		// if ptr; if elem has identical rules?
//
//		g.wr.ln(` {`)
//		g.genError(r, s)
//	}
//	g.wr.ln(`}`)
//}
//
//func (g *generator) gen(r *rules.Rule, s *spec.Spec) {
//	switch s.Kind {
//	case spec.REQUIRED:
//		switch r.Name {
//		case "notnil":
//			g.genNotnilExpr(r, s)
//		case "required":
//			g.genRequiredExpr(r, s)
//		}
//	}
//}

func (g *generator) genNotnilExpr(r *rules.Rule, s *spec.Spec) {
	switch g.obj.Type.Kind {
	case types.SLICE, types.MAP, types.INTERFACE:
		g.wr.p(`$x == nil`)
	case types.PTR, types.CHAN, types.FUNC:
		g.wr.p(`$x == nil`)
	default:
		panic("shouldn't reach")
	}
}

func (g *generator) genRequiredExpr(r *rules.Rule, s *spec.Spec) {
	switch g.obj.Type.Kind {
	case types.STRING:
		g.wr.p(`$x == ""`)
	case types.SLICE, types.MAP:
		g.wr.p(`len($x) == 0`)
	case types.INT, types.INT8, types.INT16, types.INT32, types.INT64:
		g.wr.p(`$x == 0`)
	case types.UINT, types.UINT8, types.UINT16, types.UINT32, types.UINT64:
		g.wr.p(`$x == 0`)
	case types.FLOAT32, types.FLOAT64:
		g.wr.p(`$x == 0.0`)
	case types.BOOL:
		g.wr.p(`$x == false`)
	case types.INTERFACE:
		g.wr.p(`$x == nil`)
	case types.PTR:
		g.wr.p(`$x == nil`)
	case types.STRUCT:
		if g.obj.Type.Pkg == g.file.pkg {
			g.wr.p(`$x == ($0{})`, g.obj.Type.Name)
		} else {
			p := g.file.addImport(g.obj.Type.Pkg)
			g.wr.p(`$x == ($0.$1{})`, p.name, g.obj.Type.Name)
		}
	default:
		panic("shouldn't reach")
	}
}
