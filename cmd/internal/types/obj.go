package types

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
)

type Obj struct {
	Type     *Type
	IsRules  rules.List `cmp:"-"`
	PreRules rules.List `cmp:"-"`
}

func (o *Obj) Has(kk ...rules.Kind) bool {
	return o.IsRules.Has(kk...)
}

// HasRules reports whether or not o, or any of its child
// objects, have rules other than OPTIONAL and "noguard".
func (o *Obj) HasRules() bool {
	if o == nil {
		return false
	}
	if len(o.PreRules) > 0 {
		return true
	}
	for _, r := range o.IsRules {
		if r.Spec.Kind != rules.OPTIONAL && r.Spec.Kind != rules.NOGUARD {
			return true
		}
	}
	return o.Type.HasRules()
}

func (t *Type) HasRules() bool {
	if t == nil {
		return false
	}
	if key := t.Key; key != nil && key.HasRules() {
		return true
	}
	if ele := t.Elem; ele != nil && ele.HasRules() {
		return true
	}
	for i := range t.Fields {
		if t.Fields[i].Obj.HasRules() {
			return true
		}
	}
	return false
}
