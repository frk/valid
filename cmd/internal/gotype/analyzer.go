package gotype

import (
	"go/types"

	"github.com/frk/tagutil"
)

// Analyzer maintains the state of the analysis.
type Analyzer struct {
	// The package of the type under analysis.
	//
	// NOTE: This is used primarily to resolve whether
	// the type's components need to be imported or not,
	// so it is important that the value's not changed
	// after the start of the analysis.
	pkg Pkg
	// Already analyzed types mapped by their package-qualified name.
	done map[string]*Type
}

// NewAnalyzer returns a new Analyzer instance with the p argument as
// the "root" package. The "root" package is primarily used to resolve
// wether the types encountered during analysis are imported or not.
func NewAnalyzer(p *types.Package) (a *Analyzer) {
	a = new(Analyzer)
	a.pkg.Path = p.Path()
	a.pkg.Name = p.Name()

	a.done = make(map[string]*Type)
	return a
}

// Object is used to analyze the type of a *types.Func object.
func (a *Analyzer) Object(obj types.Object) (u *Type) {
	u = a.Analyze(obj.Type())
	if pkg := obj.Pkg(); pkg != nil {
		u.Pkg.Path = pkg.Path()
		u.Pkg.Name = pkg.Name()
	}
	return u
}

// Analyze runs the analysis of the given types.Type t
// and returns the resulting gotype.Type representation.
func (a *Analyzer) Analyze(t types.Type) (u *Type) {
	u = new(Type)
	if named, ok := t.(*types.Named); ok {
		// check first if we've done this one before
		key := named.Obj().Name()
		if pkg := named.Obj().Pkg(); pkg != nil {
			key = pkg.Path() + key
		}
		if v, ok := a.done[key]; ok {
			return v
		}

		// save it asap so we don't redo it unnecessarily
		a.done[key] = u

		// NOTE this nil check is necessary for
		// "labels and objects in the Universe scope."
		if pkg := named.Obj().Pkg(); pkg != nil {
			u.Pkg.Path = pkg.Path()
			u.Pkg.Name = pkg.Name()
		}

		u.Name = named.Obj().Name()
		u.IsExported = named.Obj().Exported()
		u.Methods = a.analyzeMethods(named)
		t = named.Underlying()
	}

	switch T := t.(type) {
	case *types.Basic:
		switch T.Kind() {
		case types.Invalid:
			u.Kind = K_INVALID
		case types.Bool:
			u.Kind = K_BOOL
		case types.Int:
			u.Kind = K_INT
		case types.Int8:
			u.Kind = K_INT8
		case types.Int16:
			u.Kind = K_INT16
		case types.Int32:
			u.Kind = K_INT32
		case types.Int64:
			u.Kind = K_INT64
		case types.Uint:
			u.Kind = K_UINT
		case types.Uint8:
			u.Kind = K_UINT8
		case types.Uint16:
			u.Kind = K_UINT16
		case types.Uint32:
			u.Kind = K_UINT32
		case types.Uint64:
			u.Kind = K_UINT64
		case types.Uintptr:
			u.Kind = K_UINTPTR
		case types.Float32:
			u.Kind = K_FLOAT32
		case types.Float64:
			u.Kind = K_FLOAT64
		case types.Complex64:
			u.Kind = K_COMPLEX64
		case types.Complex128:
			u.Kind = K_COMPLEX128
		case types.String:
			u.Kind = K_STRING
		case types.UnsafePointer:
			u.Kind = K_UNSAFEPOINTER
		}
		u.IsRune = T.Name() == "rune"
		u.IsByte = T.Name() == "byte"
	case *types.Slice:
		u.Kind = K_SLICE
		u.Elem = a.Analyze(T.Elem())
	case *types.Array:
		u.Kind = K_ARRAY
		u.Elem = a.Analyze(T.Elem())
		u.ArrayLen = T.Len()
	case *types.Map:
		u.Kind = K_MAP
		u.Key = a.Analyze(T.Key())
		u.Elem = a.Analyze(T.Elem())
	case *types.Pointer:
		u.Kind = K_PTR
		u.Elem = a.Analyze(T.Elem())
	case *types.Interface:
		u.Kind = K_INTERFACE
		u.Methods = a.analyzeMethods(T)
	case *types.Signature:
		u.Kind = K_FUNC
		u.In, u.Out = a.analyzeSignature(T)
		u.IsVariadic = T.Variadic()
	case *types.Chan:
		u.Kind = K_CHAN
		// NOTE Channels aren't used for anything by the package
		// so we can ignore the element's type and the direction.
	case *types.Struct:
		u.Kind = K_STRUCT
		u.Fields = a.analyzeFields(T)
	}

	return u
}

func (a *Analyzer) analyzeFields(stype *types.Struct) (fields []*StructField) {
	for i := 0; i < stype.NumFields(); i++ {
		fvar := stype.Field(i)
		ftag := stype.Tag(i)
		tag := tagutil.New(ftag)

		f := new(StructField)
		f.Tag = ftag
		f.Name = fvar.Name()
		f.CanIgnore = (tag.First("is") == "-")
		f.IsEmbedded = fvar.Embedded()
		f.IsExported = fvar.Exported()
		f.Type = a.Analyze(fvar.Type())
		f.Var = fvar

		if pkg := fvar.Pkg(); pkg != nil {
			f.Pkg.Path = pkg.Path()
			f.Pkg.Name = pkg.Name()
		}

		fields = append(fields, f)
	}
	return fields
}

func (a *Analyzer) analyzeMethods(mm Methoder) (methods []*Method) {
	for i := 0; i < mm.NumMethods(); i++ {
		mo := mm.Method(i)

		m := new(Method)
		m.Name = mo.Name()
		m.Type = a.Analyze(mo.Type())
		m.IsExported = mo.Exported()

		// NOTE this nil check is necessary for
		// "labels and objects in the Universe scope."
		if pkg := mo.Pkg(); pkg != nil {
			m.Pkg.Path = pkg.Path()
			m.Pkg.Name = pkg.Name()
		}

		methods = append(methods, m)
	}
	return methods
}

func (a *Analyzer) analyzeSignature(sig *types.Signature) (in, out []*Var) {
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		v := params.At(i)
		in = append(in, &Var{
			Name: v.Name(),
			Type: a.Analyze(v.Type()),
		})
	}

	result := sig.Results()
	for i := 0; i < result.Len(); i++ {
		v := result.At(i)
		out = append(out, &Var{
			Name: v.Name(),
			Type: a.Analyze(v.Type()),
		})
	}
	return in, out
}
