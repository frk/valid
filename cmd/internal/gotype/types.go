package gotype

import (
	"go/types"
	"strconv"
	"strings"
)

// Pkg describes a type's package.
type Pkg struct {
	// The package import path.
	Path string
	// The package's name.
	Name string
}

// Type is the representation of a Go type.
type Type struct {
	// The type's package.
	Pkg Pkg
	// The name of a named type or empty string for unnamed types
	Name string
	// The kind of the go type.
	Kind Kind
	// Indicates whether or not the field is exported.
	IsExported bool
	// If the base type's an array type, this field will hold the array's length.
	ArrayLen int64
	// If kind is func, indicates whether or not the function is variadic.
	IsVariadic bool
	// Indicates whether or not the type is the "byte" alias type.
	IsByte bool
	// Indicates whether or not the type is the "rune" alias type.
	IsRune bool
	// If kind is map, Key will hold the info on the map's key type.
	Key *Type
	// If kind is map, Elem will hold the info on the map's value type.
	// If kind is ptr, Elem will hold the info on pointed-to type.
	// If kind is slice/array, Elem will hold the info on slice/array element type.
	Elem *Type
	// If kind is func, In & Out will hold the function's parameter and result types.
	In, Out []*Var
	// If kind is struct, Fields will hold the list of the struct's fields.
	Fields []*StructField
	// The method set of a named type or an interface type.
	Methods []*Method
}

// Reports whether or not the type's kind is is one of the provided kinds.
func (t Type) Is(kinds ...Kind) bool {
	for _, k := range kinds {
		if t.Kind == k {
			return true
		}
	}
	return false
}

// Indicates whether or not the type is an empty interface type.
func (t *Type) IsEmptyInterface() bool {
	return t.Kind == K_INTERFACE && len(t.Methods) == 0
}

// IsIncluded reports whether or not the Type was
// declared in the github.com/frk/valid package.
func (t *Type) IsIncluded() bool {
	return t.Pkg.Path == "github.com/frk/valid"
}

// IsGoString reports whether or not t is the Go builtin string type.
func (t *Type) IsGoString() bool {
	return t.Pkg == Pkg{} && t.Kind == K_STRING && t.Name == ""
}

// IsGoError reports whether or not t is the Go builtin error type.
func (t *Type) IsGoError() bool {
	return t.Pkg == Pkg{} && t.Kind == K_INTERFACE && t.Name == "error"
}

// IsGoAny reports whether or not t is the Go builtin any/interface{} type.
func (t *Type) IsGoAny() bool {
	return t.Pkg == Pkg{} && t.Kind == K_INTERFACE &&
		(t.Name == "any" || len(t.Methods) == 0)
}

// IsGoAnySlice reports whether or not t is the Go builtin []any/[]interface{} type.
func (t *Type) IsGoAnySlice() bool {
	if t.Kind == K_SLICE {
		return t.Elem.IsGoAny()
	}
	return false
}

// IsComparable reports wether or not a value of the Go
// type represented by t is comparable.
func (t *Type) IsComparable() bool {
	if t.Kind == K_MAP || t.Kind == K_SLICE || t.Kind == K_FUNC {
		return false
	}
	if t.Kind == K_ARRAY {
		return t.Elem.IsComparable()
	}
	if t.Kind == K_STRUCT {
		for _, f := range t.Fields {
			if !f.Type.IsComparable() {
				return false
			}
		}
	}
	return true
}

// IsNilable reports wether or not a value of the Go
// type represented by t can be set to nil.
func (t *Type) IsNilable() bool {
	return t.Kind == K_PTR ||
		t.Kind == K_SLICE ||
		t.Kind == K_MAP ||
		t.Kind == K_INTERFACE ||
		t.Kind == K_FUNC ||
		t.Kind == K_CHAN
}

// HasLength reports whether or not the Go type
// represented by t has a length.
func (t *Type) HasLength() bool {
	return t.Kind == K_STRING ||
		t.Kind == K_ARRAY ||
		t.Kind == K_SLICE ||
		t.Kind == K_MAP ||
		t.Kind == K_CHAN
}

// IsValid reports whether or not the "IsValid() bool"
// method belongs to the method set of the type t.
func (t *Type) HasIsValid() bool {
	for _, m := range t.Methods {
		if m.Name == "IsValid" &&
			len(m.Type.In) == 0 &&
			len(m.Type.Out) == 1 &&
			m.Type.Out[0].Type.Kind == K_BOOL {
			return true
		}
	}
	return false
}

// Reports whether or not the type t represents a pointer type of u.
func (t *Type) PtrOf(u *Type) bool {
	return t.Kind == K_PTR && t.Elem.IsIdentical(u)
}

// Reports whether the types represented by t and u are equal. Note that this
// does not handle unnamed struct, interface (non-empty), func, and channel types.
func (t *Type) IsIdentical(u *Type) bool {
	// named with same name and same package, accept
	if t.Name != "" {
		if t.Name == u.Name && t.Pkg.Path == u.Pkg.Path {
			return true
		}
	}

	// different kinds, reject
	if t.Kind != u.Kind {
		return false
	}

	// different names, reject
	if t.Name != u.Name {
		return false
	}

	// different packages, reject
	if t.Pkg.Path != u.Pkg.Path {
		return false
	}

	// channel, reject
	if t.Kind == K_CHAN {
		return false
	}

	// unnamed struct, reject
	if t.Kind == K_STRUCT && t.Name == "" {
		return false
	}

	// unnamed
	switch t.Kind {
	case K_ARRAY:
		return t.ArrayLen == u.ArrayLen && t.Elem.IsIdentical(u.Elem)
	case K_MAP:
		return t.Key.IsIdentical(u.Key) && t.Elem.IsIdentical(u.Elem)
	case K_SLICE, K_PTR:
		return t.Elem.IsIdentical(u.Elem)
	case K_INTERFACE:
		// TODO range over the methods and compare those
		return t.IsEmptyInterface() && u.IsEmptyInterface()
	case K_FUNC:
		// incompatible number of in/out parameters, reject
		if len(t.In) != len(u.In) || len(t.Out) != len(u.Out) {
			return false
		}
		// non-identical input parameter types, reject
		for i := range t.In {
			if !t.In[i].Type.IsIdentical(u.In[i].Type) {
				return false
			}
		}
		// non-identical output parameter types, reject
		for i := range t.Out {
			if !t.Out[i].Type.IsIdentical(u.Out[i].Type) {
				return false
			}
		}
	}

	// accept
	return true
}

// Reports whether or not a value of type t needs to be converted
// before it can be assigned to a variable of type u.
func (t *Type) NeedsConversion(u *Type) bool {
	if u.IsIdentical(t) {
		return false
	}
	if u.Kind == K_INTERFACE {
		return false
	}
	return true
}

// CanError reports that the type, if it *is* a K_FUNC type,
// has error as its last return value type.
func (t *Type) CanError() bool {
	if t.Kind != K_FUNC {
		return false
	}
	if n := len(t.Out); n > 0 && t.Out[n-1].Type.IsGoError() {
		return true
	}
	return false
}

// String retruns a string representation of the t Type.
func (t Type) TypeString(pkg *Pkg) string {
	if len(t.Name) > 0 {
		if pkg == nil || *pkg != t.Pkg {
			return t.Pkg.Name + "." + t.Name
		}
		return t.Name
	}

	if t.IsByte {
		return "byte"
	} else if t.IsRune {
		return "rune"
	} else if t.Kind.IsBasic() {
		return _kindstring[t.Kind]
	}

	switch t.Kind {
	case K_ARRAY:
		return "[" + strconv.FormatInt(t.ArrayLen, 10) + "]" + t.Elem.TypeString(pkg)
	case K_INTERFACE:
		if !t.IsEmptyInterface() {
			return "interface{ ... }"
		}
		return "interface{}"
	case K_MAP:
		return "map[" + t.Key.TypeString(pkg) + "]" + t.Elem.TypeString(pkg)
	case K_PTR:
		return "*" + t.Elem.TypeString(pkg)
	case K_SLICE:
		return "[]" + t.Elem.TypeString(pkg)
	case K_STRUCT:
		if len(t.Fields) > 0 {
			return "struct{ ... }"
		}
		return "struct{}"
	case K_CHAN:
		return "<chan>"
	case K_FUNC:
		in := make([]string, len(t.In))
		for i := range t.In {
			in[i] = t.In[i].Type.TypeString(pkg)
		}
		out := make([]string, len(t.Out))
		for i := range t.Out {
			out[i] = t.Out[i].Type.TypeString(pkg)
		}

		s := "func(" + strings.Join(in, ", ") + ")"
		if len(out) > 1 {
			s += " (" + strings.Join(out, ", ") + ")"
		} else if len(out) == 1 {
			s += " " + out[0]
		}
		return s
	}
	return "<unknown>"
}

////////////////////////////////////////////////////////////////////////////////
// Func
////////////////////////////////////////////////////////////////////////////////

// Func is used to represent a function.
type Func struct {
	Name string
	Type *Type
}

////////////////////////////////////////////////////////////////////////////////
// Var
////////////////////////////////////////////////////////////////////////////////

// Var is used to represent a function's parameters and results.
type Var struct {
	Name string
	Type *Type
}

func (v *Var) ShallowCopy() *Var {
	u := *v
	return &u
}

////////////////////////////////////////////////////////////////////////////////
// Fields
////////////////////////////////////////////////////////////////////////////////

// StructField describes a single struct field in a struct.
type StructField struct {
	// The package to which the field belongs.
	Pkg Pkg
	// Name of the field.
	Name string
	// The field's type.
	Type *Type
	// The field's raw, uparsed struct tag.
	Tag string
	// Indicates that the tag `is:"-"` was used.
	CanIgnore bool
	// Indicates whether or not the field is embedded.
	IsEmbedded bool
	// Indicates whether or not the field is exported.
	IsExported bool
	// Var holds a reference to the *types.Var
	// representation of the field.
	Var *types.Var `cmp:"+"`
}

func (f *StructField) CanAccess(from Pkg) bool {
	if f.Name == "_" {
		return false
	}
	if f.CanIgnore {
		return false
	}
	if !f.IsExported && !f.IsEmbedded && f.Pkg != from {
		return false
	}
	return true
}

// FieldSelector is a list of fields that represents a chain of selectors where
// the 0th field is the "root" field and the len-1 field is the "leaf" field.
type FieldSelector []*StructField

// Last returns the last field in the selector.
func (s FieldSelector) Last() *StructField {
	return s[len(s)-1]
}

// CopyWith returns a copy of the receiver with f appended to the end.
func (s FieldSelector) CopyWith(f *StructField) FieldSelector {
	s2 := make(FieldSelector, len(s)+1)
	copy(s2, s)
	s2[len(s2)-1] = f
	return s2
}

////////////////////////////////////////////////////////////////////////////////
// Methods
////////////////////////////////////////////////////////////////////////////////

// Method describes a single method in the method set of a named type or interface.
type Method struct {
	// The package to which the method belongs.
	Pkg Pkg
	// The name of the method.
	Name string
	// The method's type.
	Type *Type
	// Indicates whether or not the method is exported.
	IsExported bool
}

// Methoder represents a type with methods. It is implicitly
// implemented by *types.Interface and *types.Named.
type Methoder interface {
	NumMethods() int
	Method(i int) *types.Func
}

////////////////////////////////////////////////////////////////////////////////
// Kind
////////////////////////////////////////////////////////////////////////////////

// Kind indicates the specific kind of a Go type.
type Kind uint

const (
	// basic
	K_INVALID Kind = iota

	_basic_kind_start
	K_BOOL

	_numeric_kind_start // int/uint/float
	_integer_kind_start // int
	K_INT
	K_INT8
	K_INT16
	K_INT32
	K_INT64
	_integer_kind_end

	_unsigned_kind_start // uint
	K_UINT
	K_UINT8
	K_UINT16
	K_UINT32
	K_UINT64
	K_UINTPTR
	_unsigned_kind_end

	K_FLOAT32
	K_FLOAT64
	_numeric_kind_end

	K_COMPLEX64
	K_COMPLEX128
	K_STRING
	K_UNSAFEPOINTER
	_basic_kind_end

	// non-basic
	K_ARRAY     // try to validate individual elements
	K_INTERFACE // try to validate ... ???
	K_MAP       // try to validate individual elements
	K_PTR       // try to validate the element
	K_SLICE     // try to validate the individual elements
	K_STRUCT    // try to validate the individual fields
	K_CHAN      // don't validate
	K_FUNC      // don't validate

	// alisases (basic)
	K_BYTE = K_UINT8
	K_RUNE = K_INT32
)

// Reports whether or not k is of a basic kind.
func (k Kind) IsBasic() bool { return _basic_kind_start < k && k < _basic_kind_end }

// Reports whether or not k is of the numeric kind, note that this
// does not include the complex64 and complex128 kinds.
func (k Kind) IsNumeric() bool { return _numeric_kind_start < k && k < _numeric_kind_end }

// Reports whether or not k is one of the int types. (does not include unsigned integers)
func (k Kind) IsInteger() bool { return _integer_kind_start < k && k < _integer_kind_end }

// Reports whether or not k is one of the uint types.
func (k Kind) IsUnsigned() bool { return _unsigned_kind_start < k && k < _unsigned_kind_end }

// Reports whether or not k is one of the float types.
func (k Kind) IsFloat() bool { return k == K_FLOAT32 || k == K_FLOAT64 }

// BasicString returns a string representation of the basic kind k.
func (k Kind) BasicString() string {
	if k.IsBasic() {
		return _kindstring[k]
	}
	return "<unknown>"
}

func (k Kind) String() string {
	if int(k) < len(_kindstring) {
		return _kindstring[k]
	}
	return "<unknown>"
}

var _kindstring = [...]string{
	K_INVALID:    "<invalid>",
	K_BOOL:       "bool",
	K_INT:        "int",
	K_INT8:       "int8",
	K_INT16:      "int16",
	K_INT32:      "int32",
	K_INT64:      "int64",
	K_UINT:       "uint",
	K_UINT8:      "uint8",
	K_UINT16:     "uint16",
	K_UINT32:     "uint32",
	K_UINT64:     "uint64",
	K_UINTPTR:    "uintptr",
	K_FLOAT32:    "float32",
	K_FLOAT64:    "float64",
	K_COMPLEX64:  "complex64",
	K_COMPLEX128: "complex128",
	K_STRING:     "string",

	// ...
	K_ARRAY:     "array",
	K_INTERFACE: "interface",
	K_MAP:       "map",
	K_PTR:       "ptr",
	K_SLICE:     "slice",
	K_STRUCT:    "struct",
	K_CHAN:      "chan",
	K_FUNC:      "func",
}

////////////////////////////////////////////////////////////////////////////////
// Assignment
////////////////////////////////////////////////////////////////////////////////

type AssignmentStatus uint

const (
	ASSIGNMENT_INVALID AssignmentStatus = iota // cannot assign
	ASSIGNMENT_CONVERT                         // can assign but needs explicit converstion
	ASSIGNMENT_OK                              // can assign as is
)

// CanConvert reports whether or not a value of type u
// can be converted to a value of type t.
func (t *Type) CanConvert(u *Type) bool {
	return t.CanAssign(u) > ASSIGNMENT_INVALID
}

// CanAssign reports whether or not a value of type u can be
// assigned to a variable of type t. Note that this does not
// handle unnamed struct, interface, func, and channel types.
func (t *Type) CanAssign(u *Type) AssignmentStatus {
	// if same, accept
	if u.IsIdentical(t) {
		return ASSIGNMENT_OK
	}

	// if t is interface{}, accept
	if t.IsEmptyInterface() {
		return ASSIGNMENT_OK
	}
	// if t is interface{ ... } and u implements t, accept
	if t.Kind == K_INTERFACE && u.Implements(t) {
		return ASSIGNMENT_OK
	}

	// same basic kind, accept
	if t.Kind == u.Kind && t.Kind.IsBasic() {
		// TODO if u is unnamed, e.g. a rule arg constant
		// then this should return
		return ASSIGNMENT_CONVERT
	}

	// both numeric, accept
	if t.Kind.IsNumeric() && u.Kind.IsNumeric() {
		return ASSIGNMENT_CONVERT
	}

	// string from []byte, []rune, []uint8, and []int32, accept
	if t.Kind == K_STRING && u.Kind == K_SLICE && u.Elem.Name == "" &&
		(u.Elem.Kind == K_UINT8 || u.Elem.Kind == K_INT32) {
		return ASSIGNMENT_CONVERT
	}
	// string to []byte, []rune, []uint8, and []int32, accept
	if u.Kind == K_STRING && t.Kind == K_SLICE && t.Elem.Name == "" &&
		(t.Elem.Kind == K_UINT8 || t.Elem.Kind == K_INT32) {
		return ASSIGNMENT_CONVERT
	}

	// element types (and key & len) of non-basic are equal, accept
	if t.Kind == u.Kind && !t.Kind.IsBasic() {
		switch t.Kind {
		case K_ARRAY:
			if t.ArrayLen == u.ArrayLen && t.Elem.IsIdentical(u.Elem) {
				return ASSIGNMENT_CONVERT
			}
		case K_MAP:
			if t.Key.IsIdentical(u.Key) && t.Elem.IsIdentical(u.Elem) {
				return ASSIGNMENT_CONVERT
			}
		case K_SLICE, K_PTR:
			if t.Elem.IsIdentical(u.Elem) {
				return ASSIGNMENT_CONVERT
			}
		}
	}
	return ASSIGNMENT_INVALID
}

func (t *Type) Implements(u *Type) bool {
	if u.Kind != K_INTERFACE {
		return false
	}

	if t.Kind == K_PTR {
		t = t.Elem
	}

methodLoop:
	for _, um := range u.Methods {
		for _, tm := range t.Methods {
			if um.Name != tm.Name {
				continue // try next
			}
			if !um.IsExported && um.Pkg != tm.Pkg {
				continue // try next
			}
			if !tm.Type.IsIdentical(um.Type) {
				continue // try next
			}

			// tm matches um
			continue methodLoop
		}

		// no method in t matched um
		return false
	}

	return true
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
func (t *Type) VisibleFields() []*StructField {
	x := t

	if x.Kind == K_PTR {
		x = x.Elem
	}
	if x.Kind != K_STRUCT {
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
			t := f.Type
			if t.Kind == K_PTR {
				t = t.Elem
			}
			if t.Kind == K_STRUCT {
				w.walk(t, depth+1)
			}
		}
	}

	delete(w.visiting, t)
}
