package types

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
	// Indicates that the type is recursive.
	IsRecursive bool
	// If kind is map, Key will hold the info on the map's key type.
	Key *Type
	// If kind is map, Elem will hold the info on the map's value type.
	// If kind is ptr, Elem will hold the info on pointed-to type.
	// If kind is slice/array, Elem will hold the info on slice/array element type.
	Elem *Type
	// The method set of a named type or an interface type.
	MethodSet []*Method
	Embeddeds []*Type
	// If kind is func, In & Out will hold the
	// function's parameter and result types.
	In, Out []*Var
	// If kind is struct, Fields will hold the
	// list of the struct's fields.
	Fields []*Field
	// If the Type is an instantiated named type then
	// Origin points to the original generic type.
	Origin *Type
	// If the Type is an instantiated named type then
	// TypeArgs is the list of type arguments.
	TypeArgs []*Type
	// If the Type is a generic named type or a generic function
	// signature then TypeParams is the list of type parameters.
	TypeParams []*TypeParam
	// If the Type is a union then Terms holds
	// the union's of terms.
	Terms []*Term
}

type TypeParam struct {
	// The type param's package.
	Pkg Pkg
	// The type name of the type param.
	Name string
	// The constraint specified for the type param.
	Constraint *Type
}

type TypeParamer interface {
	TypeParams() *types.TypeParamList
}

type Term struct {
	// Indicates whether or not the term
	// was declared with a tilde.
	Tilde bool
	// The term's type.
	Type *Type
}

// Underlying returns a shallow copy of t with its
// name, package, and method set removed, ...
func (t *Type) Underlying() *Type {
	u := *t
	u.Pkg = Pkg{}
	u.Name = ""
	u.IsByte = false
	u.IsRune = false
	if u.Kind != INTERFACE {
		u.MethodSet = nil
		u.Embeddeds = nil
	}
	return &u
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
	case ARRAY:
		return "[" + strconv.FormatInt(t.ArrayLen, 10) + "]" + t.Elem.TypeString(pkg)
	case INTERFACE:
		if !t.IsEmptyIface() {
			return "interface{ ... }"
		}
		return "interface{}"
	case MAP:
		return "map[" + t.Key.TypeString(pkg) + "]" + t.Elem.TypeString(pkg)
	case PTR:
		return "*" + t.Elem.TypeString(pkg)
	case SLICE:
		return "[]" + t.Elem.TypeString(pkg)
	case STRUCT:
		if len(t.Fields) > 0 {
			return "struct{ ... }"
		}
		return "struct{}"
	case CHAN:
		return "<chan>"
	case FUNC:
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
// Methods
////////////////////////////////////////////////////////////////////////////////

// Method describes a single method in the method set of a named type or interface.
type Method struct {
	// The package to which the method belongs.
	Pkg Pkg
	// The name of the method.
	Name string
	// The method's type (func signature).
	Type *Type
	// Indicates whether or not the method is exported.
	IsExported bool
	// Indicates whether or not the method's receiver is a pointer.
	IsPtr bool
}
