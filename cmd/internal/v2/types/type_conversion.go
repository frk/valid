package types

// IsConvertibleTo reports whether a value of
// type V can be converted to a value of type T.
func (V *Type) IsConvertibleTo(T *Type) bool {
	if V.IsAssignableTo(T) {
		return true
	}

	// basic with same kind, accept
	if V.Kind.IsBasic() && V.Kind == T.Kind {
		return true
	}

	// non-basic with identical underlying type, accept
	if !V.Kind.IsBasic() && V.Kind == T.Kind {
		Vu := V.Underlying()
		Tu := T.Underlying()
		if Vu.IsIdenticalTo(Tu) {
			return true
		}
	}

	// "V and T are unnamed pointer types and their pointer base types
	// have identical underlying types if tags are ignored
	// and their pointer base types are not type parameters"
	if V.Kind == PTR && T.Kind == PTR && !V.IsNamed() && !T.IsNamed() {
		Vubase := V.Elem.Underlying()
		Tubase := T.Elem.Underlying()
		if Vubase.IsIdenticalTo(Tubase) {
			return true
		}
	}

	// "V and T are both integer or floating point types"
	if V.Kind.IsNumeric() && T.Kind.IsNumeric() {
		return true
	}

	// "V and T are both complex types"
	if V.Kind.IsComplex() && T.Kind.IsComplex() {
		return true
	}

	// "V is an integer or a slice of bytes or runes and T is a string type"
	if (V.Kind.IsInteger() || V.Kind.IsUnsigned() || V.IsBytesOrRunes()) && T.Kind == STRING {
		return true
	}

	// "V is a string and T is a slice of bytes or runes"
	if V.Kind == STRING && T.IsBytesOrRunes() {
		return true
	}

	// "V is a slice, T is an array or pointer-to-array type,
	// and the slice and array types have identical element types."
	if V.Kind == SLICE {
		switch T.Kind {
		case ARRAY:
			if V.Elem.IsIdenticalTo(T.Elem) {
				// NOTE: conversion of slices to arrays requires go1.20 or later
				return true
			}
		case PTR:
			if A := T.Elem; A.Kind == ARRAY {
				if V.Elem.IsIdenticalTo(A.Elem) {
					// NOTE: conversion of slices to array pointers requires go1.17 or later
					return true
				}
			}
		}
	}

	// TODO(mkopriva): add support for type parameters.
	return false
}
