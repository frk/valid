package check

import (
	"strings"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (c *checker) addFRefs(r *rules.Rule, ff types.FieldChain) error {
	for _, a := range r.Args {
		if !a.IsFieldRef() {
			continue
		}
		if err := c.addFRef(a, ff); err != nil {
			return c.err(err, errOpts{sf: ff.Last(), r: r, rs: c.Info.SpecMap[r]})
		}
	}
	return nil
}

func (c *checker) addFRef(a *rules.Arg, ff types.FieldChain) error {
	var root *types.Type
	switch a.Type {
	case rules.ARG_FIELD_ABS:
		root = c.v.Type
	case rules.ARG_FIELD_REL:
		if len(ff) > 1 {
			root = ff[len(ff)-2].Obj.Type
		} else {
			root = c.v.Type
		}
	}

	lf := ff.Last() // "leaf field"
	cur := root
	ref := strings.Split(a.Value, ".")
	for i, last := 0, len(ref)-1; i < len(ref); i++ {
		var sf *types.StructField
		var ok bool

		for _, f := range cur.VisibleFields() {
			if ref[i] == f.Name {
				sf, ok = f, true
				cur = sf.Obj.Type
				break
			}
		}

		if !ok {
			return &Error{C: E_FIELD_UNKNOWN, ty: ff.Last().Obj.Type, ra: a}
		}
		if i == last {
			c.Info.FRefMap[a] = sf
			c.Info.FDepMap[lf] = append(c.Info.FDepMap[lf], sf)
		}
	}
	return nil
}
