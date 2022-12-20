package types

// IsAssignableTo reports whether a value of type
// V is assignable to a variable of type T.
func (V *Type) IsAssignableTo(T *Type) bool {
	// if identical, accept
	if T.IsIdenticalTo(V) {
		return true
	}

	// if t is interface{}, accept
	if T.IsEmptyIface() {
		return true
	}

	// if T is interface{ ... } and V implements T, accept
	if T.Kind == INTERFACE && V.Implements(T) {
		return true
	}

	// if at least one is unnamed, check their underlying types
	if !T.IsNamed() || !V.IsNamed() {
		return T.Underlying().IsIdenticalTo(V.Underlying())
	}

	return false
}
