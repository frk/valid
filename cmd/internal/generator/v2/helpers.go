package generate

import (
	"github.com/frk/valid/cmd/internal/rules/spec"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) hasRequiredRule(o *types.Obj) bool {
	if rs, ok := g.info.RuleMap[o]; ok {
		for _, r := range rs.Is.Rules {
			s := g.info.SpecMap[r]
			if s.Kind == spec.REQUIRED {
				return true
			}
		}
	}
	return false
}

func (g *generator) hasOptionalRule(o *types.Obj) bool {
	if rs, ok := g.info.RuleMap[o]; ok {
		for _, r := range rs.Is.Rules {
			s := g.info.SpecMap[r]
			if s.Kind == spec.OPTIONAL {
				return true
			}
		}
	}
	return false
}

func (g *generator) hasNoguardRule(o *types.Obj) bool {
	if rs, ok := g.info.RuleMap[o]; ok {
		for _, r := range rs.Is.Rules {
			s := g.info.SpecMap[r]
			if s.Kind == spec.NOGUARD {
				return true
			}
		}
	}
	return false
}
