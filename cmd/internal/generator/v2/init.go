package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
)

// genInit produces a top-level init function.
func (g *generator) genInit() {
	if len(g.file.reMap) == 0 {
		return
	}

	fn := specs.RegisterRegexpFunc()
	reList := make([]string, len(g.file.reMap))
	for re, idx := range g.file.reMap {
		reList[idx] = re
	}

	g.L(`func init() {`)
	for _, re := range reList {
		g.L("$0(`$1`)", fn, re)
	}
	g.L(`}`)
	g.L(``)
}
