package types

// NOTE: this implementation, and the corresponding tests, is a slightly
// modified and adapted version of https://pkg.go.dev/reflect#VisibleFields.
//
// VisibleFields returns all the visible fields in t, which must be
// a struct type or a pointer to a struct type. A field is defined as
// visible if it's accessible directly from the type's instance.
//
// The returned fields include fields inside anonymous struct members and
// unexported fields. They follow the same order found in the struct, with
// anonymous fields followed immediately by their promoted fields.
func (t *Type) VisibleFields() []*StructField {
	x := t

	if x.Kind == PTR {
		x = x.Elem.Type
	}
	if x.Kind != STRUCT {
		return nil
	}

	w := &visibleFieldsWalker{
		fields:   make([]*StructField, 0, len(x.Fields)),
		depth:    make(map[*StructField]int, len(x.Fields)),
		hidden:   make(map[*StructField]bool),
		byName:   make(map[string]*StructField, len(x.Fields)),
		visiting: make(map[*Type]bool),
	}
	w.walk(t, 0)

	// Remove all the fields that have been hidden.
	j := 0
	for i := range w.fields {
		f := w.fields[i]
		if w.hidden[f] {
			continue
		}
		if i != j {
			// A field has been removed. We need to shuffle
			// all the subsequent elements up.
			w.fields[j] = f
		}
		j++
	}
	return w.fields[:j]
}

type visibleFieldsWalker struct {
	fields   []*StructField
	depth    map[*StructField]int
	hidden   map[*StructField]bool
	byName   map[string]*StructField
	visiting map[*Type]bool
}

func (w *visibleFieldsWalker) walk(t *Type, depth int) {
	if w.visiting[t] {
		return
	}
	w.visiting[t] = true

	for i := 0; i < len(t.Fields); i++ {
		// Use a shallow copy, otherwise a struct that's embedded
		// multiple times at different depths will result in the
		// same field pointer to appear multiple times in w.fields.
		// And the subsequent hidden-field-removal process will,
		// because of the matching pointer, remove all instances
		// from w.fields, rather than keeping the shallowest one.
		v := *t.Fields[i]
		f := &v

		add := true
		if old, ok := w.byName[f.Name]; ok {
			switch oldepth := w.depth[old]; {
			case depth == oldepth:
				// hide old field because access
				// is ambiguous with the new field.
				w.hidden[old] = true
				add = false
			case depth < oldepth:
				// hide old field because its
				// deeper than the new field.
				w.hidden[old] = true
			case depth > oldepth:
				add = false
			}
		}

		if add {
			w.fields = append(w.fields, f)
			w.depth[f] = depth
			w.byName[f.Name] = f
		}

		if f.IsEmbedded {
			t := f.Obj.Type
			if t.Kind == PTR {
				t = t.Elem.Type
			}
			if t.Kind == STRUCT {
				w.walk(t, depth+1)
			}
		}
	}

	delete(w.visiting, t)
}
