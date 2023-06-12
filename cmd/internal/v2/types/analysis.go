package types

import (
	"fmt"
	"go/types"
	"sync"

	"github.com/frk/valid/cmd/internal/search"
)

var _ = fmt.Println

// Analyze runs the analyzer of the given types.Type t
// and returns the corresponding Type representation.
func Analyze(t types.Type, ast *search.AST) *Type {
	a := &analyzer{ast: ast}
	return a.analyzeType(t)
}

func AnalyzeObject(obj types.Object, ast *search.AST) *Type {
	a := &analyzer{ast: ast}
	t := a.analyzeType(obj.Type())
	if pkg := obj.Pkg(); pkg != nil {
		t.Pkg.Path = pkg.Path()
		t.Pkg.Name = pkg.Name()
	}
	return t
}

////////////////////////////////////////////////////////////////////////////////

type analyzer struct {
	ast      *search.AST
	visiting map[string]*Type
}

func (a *analyzer) analyzeType(t types.Type) (u *Type) {
	u = new(Type)

	var key string
	if named, ok := t.(*types.Named); ok {
		// check for recursive types
		key = named.Obj().Name()
		if pkg := named.Obj().Pkg(); pkg != nil {
			key = pkg.Path() + key
		}
		if v, ok := a.visiting[key]; ok {
			v.IsRecursive = true
			return v
		} else {
			if a.visiting == nil {
				a.visiting = make(map[string]*Type)
			}
			a.visiting[key] = u
		}

		// NOTE this nil check is necessary for
		// "labels and objects in the Universe scope."
		if pkg := named.Obj().Pkg(); pkg != nil {
			u.Pkg.Path = pkg.Path()
			u.Pkg.Name = pkg.Name()
		}

		u.Name = named.Obj().Name()
		u.IsExported = named.Obj().Exported()
		u.MethodSet = a.analyzeMethodSet(named)
		u.TypeArgs = a.analyzeTypeArgs(named)
		u.TypeParams = a.analyzeTypeParams(named)
		if o := named.Origin(); o != nil && o != t {
			u.Origin = a.analyzeType(o)
		}
		t = named.Underlying()
	}

	switch T := t.(type) {
	case *types.Basic:
		switch T.Kind() {
		case types.Invalid:
			u.Kind = INVALID
		case types.Bool:
			u.Kind = BOOL
		case types.Int:
			u.Kind = INT
		case types.Int8:
			u.Kind = INT8
		case types.Int16:
			u.Kind = INT16
		case types.Int32:
			u.Kind = INT32
		case types.Int64:
			u.Kind = INT64
		case types.Uint:
			u.Kind = UINT
		case types.Uint8:
			u.Kind = UINT8
		case types.Uint16:
			u.Kind = UINT16
		case types.Uint32:
			u.Kind = UINT32
		case types.Uint64:
			u.Kind = UINT64
		case types.Uintptr:
			u.Kind = UINTPTR
		case types.Float32:
			u.Kind = FLOAT32
		case types.Float64:
			u.Kind = FLOAT64
		case types.Complex64:
			u.Kind = COMPLEX64
		case types.Complex128:
			u.Kind = COMPLEX128
		case types.String:
			u.Kind = STRING
		case types.UnsafePointer:
			u.Kind = UNSAFEPOINTER
		}
		u.IsRune = T.Name() == "rune"
		u.IsByte = T.Name() == "byte"
	case *types.Slice:
		u.Kind = SLICE
		u.Elem = a.analyzeType(T.Elem())
	case *types.Array:
		u.Kind = ARRAY
		u.Elem = a.analyzeType(T.Elem())
		u.ArrayLen = T.Len()
	case *types.Map:
		u.Kind = MAP
		u.Key = a.analyzeType(T.Key())
		u.Elem = a.analyzeType(T.Elem())
	case *types.Pointer:
		u.Kind = PTR
		u.Elem = a.analyzeType(T.Elem())
	case *types.Interface:
		u.Kind = INTERFACE
		u.MethodSet = a.analyzeMethodSet(T)
		u.Embeddeds = a.analyzeEmbeddeds(T)
	case *types.Signature:
		u.Kind = FUNC
		u.In, u.Out = a.analyzeSignature(T)
		u.TypeParams = a.analyzeTypeParams(T)
		u.IsVariadic = T.Variadic()
	case *types.Chan:
		u.Kind = CHAN
		// NOTE Channels aren't used for anything by the package
		// so we can ignore the element's type and the direction.
	case *types.Struct:
		u.Kind = STRUCT
		u.Fields = a.analyzeFields(T)
	case *types.TypeParam:
		p := a.analyzeType(T.Constraint())
		*u = *p
	case *types.Union:
		u.Kind = UNION
		u.Terms = a.analyzeTerms(T)
	}

	if key != "" {
		delete(a.visiting, key)
	}
	return u
}

func (a *analyzer) analyzeFields(stype *types.Struct) (fields []*Field) {
	for i := 0; i < stype.NumFields(); i++ {
		fvar := stype.Field(i)
		ftag := stype.Tag(i)

		f := new(Field)
		f.Tag = ftag
		f.Name = fvar.Name()
		f.IsEmbedded = fvar.Embedded()
		f.IsExported = fvar.Exported()
		f.Type = a.analyzeType(fvar.Type())

		if pkg := fvar.Pkg(); pkg != nil {
			f.Pkg.Path = pkg.Path()
			f.Pkg.Name = pkg.Name()
		}

		fields = append(fields, f)

		storePosition(f, a.ast.FileAndLine(fvar))
	}
	return fields
}

func (a *analyzer) analyzeMethodSet(t types.Type) (methods []*Method) {
	if _, ok := t.(*types.Named); ok {
		t = types.NewPointer(t)
	}

	mset := types.NewMethodSet(t)
	for i := 0; i < mset.Len(); i++ {
		obj := mset.At(i).Obj()

		m := new(Method)
		m.Name = obj.Name()
		m.Type = a.analyzeType(obj.Type())
		m.IsExported = obj.Exported()

		// NOTE this nil check is necessary for
		// "labels and objects in the Universe scope."
		if pkg := obj.Pkg(); pkg != nil {
			m.Pkg.Path = pkg.Path()
			m.Pkg.Name = pkg.Name()
		}

		// determine whether the method's declared with a pointer receiver
		if sig, ok := obj.Type().(*types.Signature); ok {
			if recv := sig.Recv(); recv != nil {
				_, m.IsPtr = recv.Type().(*types.Pointer)
			}
		}

		methods = append(methods, m)
	}

	return methods
}

func (a *analyzer) analyzeEmbeddeds(iface *types.Interface) (embeddeds []*Type) {
	for i := 0; i < iface.NumEmbeddeds(); i++ {
		t := a.analyzeType(iface.EmbeddedType(i))
		embeddeds = append(embeddeds, t)
	}
	return embeddeds
}

func (a *analyzer) analyzeTerms(union *types.Union) (terms []*Term) {
	for i := 0; i < union.Len(); i++ {
		ut := union.Term(i)

		t := new(Term)
		t.Tilde = ut.Tilde()
		t.Type = a.analyzeType(ut.Type())
		terms = append(terms, t)
	}
	return terms
}

func (a *analyzer) analyzeTypeArgs(named *types.Named) (args []*Type) {
	list := named.TypeArgs()
	if list == nil {
		return nil
	}

	for i := 0; i < list.Len(); i++ {
		t := a.analyzeType(list.At(i))
		args = append(args, t)
	}
	return args
}

func (a *analyzer) analyzeTypeParams(t namedOrSignature) (params []*TypeParam) {
	list := t.TypeParams()
	if list == nil {
		return nil
	}

	for i := 0; i < list.Len(); i++ {
		tp := list.At(i)

		p := new(TypeParam)
		if typeName := tp.Obj(); typeName != nil {
			p.Name = typeName.Name()
			if pkg := typeName.Pkg(); pkg != nil {
				p.Pkg.Path = pkg.Path()
				p.Pkg.Name = pkg.Name()
			}
		}
		p.Constraint = a.analyzeType(tp.Constraint())

		params = append(params, p)
	}
	return params
}

func (a *analyzer) analyzeSignature(sig *types.Signature) (in, out []*Var) {
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		v := params.At(i)
		in = append(in, &Var{
			Name: v.Name(),
			Type: a.analyzeType(v.Type()),
		})
	}

	result := sig.Results()
	for i := 0; i < result.Len(); i++ {
		v := result.At(i)
		out = append(out, &Var{
			Name: v.Name(),
			Type: a.analyzeType(v.Type()),
		})
	}
	return in, out
}

////////////////////////////////////////////////////////////////////////////////
// helpers

type namedOrSignature interface {
	TypeParams() *types.TypeParamList
}

////////////////////////////////////////////////////////////////////////////////
// cache

var typeCache = struct {
	sync.RWMutex
	// m maps package-path qualified type names
	// to an already analyzed type.
	m map[string]*Type
}{m: make(map[string]*Type)}
