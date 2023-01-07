package types

import (
	"github.com/frk/tagutil"
)

// StructField describes a single struct field in a struct.
type StructField struct {
	// The package to which the field belongs.
	Pkg Pkg
	// Name of the field.
	Name string
	// The field's type.
	Obj *Obj
	// The field's raw, uparsed struct tag.
	Tag string
	// Indicates whether or not the field is embedded.
	IsEmbedded bool
	// Indicates whether or not the field is exported.
	IsExported bool
}

// CanSkip reports whether the field should
// be skipped by the type checker and generator
func (f *StructField) CanSkip(pkg Pkg) bool {
	return f.Name == "_" || tagutil.New(f.Tag).First("is") == "-" ||
		(!f.IsExported && !f.IsEmbedded && f.Pkg != pkg)
}

// FieldChain is a list of fields that represents a chain of fields where
// the 0th field is the "root" field and the len-1 field is the "leaf" field.
type FieldChain []*StructField

// First returns the first, the root field in the chain.
func (ff FieldChain) First() *StructField {
	return ff[0]
}

// Last returns the last, the leaf field in the chain.
func (ff FieldChain) Last() *StructField {
	if len(ff) > 0 {
		return ff[len(ff)-1]
	}
	return nil
}

// CopyWith returns a copy of the receiver with f appended to the end.
func (ff FieldChain) CopyWith(f *StructField) FieldChain {
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
