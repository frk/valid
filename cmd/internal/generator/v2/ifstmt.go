package generate

import (
	"github.com/frk/valid/cmd/internal/rules/spec"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) ifStmt(rs *rules.Set) {
	g.wr.p(`if `)
	for _, r := range rs.Is.Rules {
		s := g.info.SpecMap[r]
		g.ruleCond(r, s)

		// if ptr; if elem has identical rules?

		g.wr.ln(` {`)
		g.genError(r, s)
	}
	g.wr.ln(`}`)
}

func (g *generator) ruleCond(r *rules.Rule, s *spec.Spec) {
	switch s.Kind {
	case spec.REQUIRED:
		switch r.Name {
		case "notnil":
			g.genNotnilExpr(r, s)
		case "required":
			g.genRequiredExpr(r, s)
		}
	}
}
