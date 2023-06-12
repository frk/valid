package types

import (
	"github.com/frk/tagutil"
)

// Field describes a single struct field in a struct.
type Field struct {
	// The package to which the field belongs.
	Pkg Pkg
	// Name of the field.
	Name string
	// The field's type.
	Type *Type
	// The field's raw, uparsed struct tag.
	Tag string
	// Indicates whether or not the field is embedded.
	IsEmbedded bool
	// Indicates whether or not the field is exported.
	IsExported bool
}

// CanSkip reports whether the field should
// be skipped by the type checker and generator
func (f *Field) CanSkip(pkg Pkg) bool {
	return f.Name == "_" || tagutil.New(f.Tag).First("is") == "-" ||
		(!f.IsExported && !f.IsEmbedded && f.Pkg != pkg)
}

// FieldChain is a list of fields that represents a chain of fields where
// the 0th field is the "root" field and the len-1 field is the "leaf" field.
type FieldChain []*Field

// First returns the first, the root field in the chain.
func (ff FieldChain) First() *Field {
	return ff[0]
}

// Last returns the last, the leaf field in the chain.
func (ff FieldChain) Last() *Field {
	if len(ff) > 0 {
		return ff[len(ff)-1]
	}
	return nil
}

// CopyWith returns a copy of the receiver with f appended to the end.
func (ff FieldChain) CopyWith(f *Field) FieldChain {
	ff2 := make(FieldChain, len(ff)+1)
	copy(ff2, ff)
	ff2[len(ff2)-1] = f
	return ff2
}

func (ff FieldChain) String() (s string) {
	for _, f := range ff {
		s += f.Name + "."
	}
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

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
func (t *Type) VisibleFields() []*Field {
	x := t

	if x.Kind == PTR {
		x = x.Elem
	}
	if x.Kind != STRUCT {
		return nil
	}

	w := &visibleFieldsWalker{
		fields:   make([]*Field, 0, len(x.Fields)),
		depth:    make(map[*Field]int, len(x.Fields)),
		hidden:   make(map[*Field]bool),
		byName:   make(map[string]*Field, len(x.Fields)),
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
	fields   []*Field
	depth    map[*Field]int
	hidden   map[*Field]bool
	byName   map[string]*Field
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
			t := f.Type
			if t.Kind == PTR {
				t = t.Elem
			}
			if t.Kind == STRUCT {
				w.walk(t, depth+1)
			}
		}
	}

	delete(w.visiting, t)
}
