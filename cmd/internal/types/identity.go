package types

import (
	"sort"
)

// IsPtrOf reports whether type T is a pointer with its base type identical to type U.
func (T *Type) IsPtrOf(U *Type) bool {
	return T.Kind == PTR && T.Elem.Type.IsIdenticalTo(U)
}

// IsIdenticalTo reports whether type T is considered identical to type U.
func (T *Type) IsIdenticalTo(U *Type) bool {
	if T == U {
		return true
	}

	// if named compare names and package paths
	if T.IsNamed() || U.IsNamed() {
		// TODO(mkopriva): when its an instantiated type with different
		// type args then this would be incorrent, wouldn't it?
		return T.Name == U.Name && T.Pkg.Path == U.Pkg.Path
	}

	// if both are basic and unnamed and same kind, accept
	if T.Kind.IsBasic() && U.Kind.IsBasic() && !T.IsNamed() && !U.IsNamed() {
		if T.Kind == U.Kind {
			return true
		}
	}

	// different kinds, reject
	if T.Kind != U.Kind {
		return false
	}

	switch T.Kind {
	case PTR:
		// Two pointer types are identical if they have identical base types.
		return T.Elem.Type.IsIdenticalTo(U.Elem.Type)
	case ARRAY:
		// Two array types are identical if they have identical
		// element types and the same array length.
		return T.ArrayLen == U.ArrayLen && T.Elem.Type.IsIdenticalTo(U.Elem.Type)
	case SLICE:
		// Two slice types are identical if they have identical element types.
		return T.Elem.Type.IsIdenticalTo(U.Elem.Type)
	case MAP:
		// Two map types are identical if they have identical key and value types.
		return T.Key.Type.IsIdenticalTo(U.Key.Type) && T.Elem.Type.IsIdenticalTo(U.Elem.Type)
	case CHAN:
		// NOTE: not supported at the moment
		return false
	case STRUCT:
		// Two struct types are identical if they have the same sequence of fields,
		// and if corresponding fields have the same names, and identical types,
		// and identical tags. Lower-case field names from different packages are
		// always different.
		if len(T.Fields) != len(U.Fields) {
			return false
		}

		for i, f := range T.Fields {
			g := U.Fields[i]
			if f.Name != g.Name {
				return false
			}
			if !f.IsExported && f.Pkg.Path != g.Pkg.Path {
				return false
			}
			if !f.Obj.Type.IsIdenticalTo(g.Obj.Type) {
				return false
			}
		}
		return true
	case FUNC:
		// Two function types are identical if they have the same number of
		// parameters and result values, corresponding parameter and result types
		// are identical, and either both functions are variadic or neither is.
		// Parameter and result names are not required to match, and type
		// parameters are considered identical modulo renaming.

		// incompatible number of in/out parameters, reject
		if len(T.In) != len(U.In) || len(T.Out) != len(U.Out) {
			return false
		}
		// number of type parameters must also be identical
		if len(T.TypeParams) != len(U.TypeParams) {
			return false
		}

		if len(T.TypeParams) > 0 {
			// TODO(mkopriva): handle type parameters
		}

		// non-identical input parameter types, reject
		for i := range T.In {
			if !T.In[i].Type.IsIdenticalTo(U.In[i].Type) {
				return false
			}
		}
		// non-identical output parameter types, reject
		for i := range T.Out {
			if !T.Out[i].Type.IsIdenticalTo(U.Out[i].Type) {
				return false
			}
		}
		return true
	case INTERFACE:
		// Two interface types are identical if they describe the same type sets.
		// With the existing implementation restriction, this simplifies to:
		//
		// Two interface types are identical if they have the same set of methods with
		// the same names and identical function types, and if any type restrictions
		// are the same. Lower-case method names from different packages are always
		// different. The order of the methods is irrelevant.
		if len(T.MethodSet) != len(U.MethodSet) {
			return false
		}

		if len(T.TypeParams) > 0 {
			// TODO(mkopriva): handle type parameters
		}

		tmset := append([]*Method{}, T.MethodSet...)
		sort.Slice(tmset, func(i, j int) bool { return tmset[i].Name < tmset[j].Name })

		umset := append([]*Method{}, U.MethodSet...)
		sort.Slice(umset, func(i, j int) bool { return umset[i].Name < umset[j].Name })

		for i, m := range tmset {
			n := umset[i]
			if m.Name != n.Name {
				return false
			}
			if !m.IsExported && m.Pkg.Path != n.Pkg.Path {
				return false
			}
			if !m.Type.IsIdenticalTo(n.Type) {
				return false
			}
		}
		return true
	case UNION:
		// TODO(mkopriva): needs to be implemented
		return false
	}

	return false
}
