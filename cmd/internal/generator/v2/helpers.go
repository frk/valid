package generate

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) hasRequiredRule(o *types.Obj) bool {
	for _, r := range o.IsRules {
		if r.Spec.Kind == rules.REQUIRED {
			return true
		}
	}
	return false
}

func (g *generator) hasOptionalRule(o *types.Obj) bool {
	for _, r := range o.IsRules {
		if r.Spec.Kind == rules.OPTIONAL {
			return true
		}
	}
	return false
}

func (g *generator) hasNoguardRule(o *types.Obj) bool {
	for _, r := range o.IsRules {
		if r.Spec.Kind == rules.NOGUARD {
			return true
		}
	}
	return false
}
