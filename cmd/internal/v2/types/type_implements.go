package types

// Implements reports whether type T implements interface I.
func (T *Type) Implements(I *Type) bool {
	if I.Kind != INTERFACE {
		return false
	}

	if I.IsEmptyIface() {
		// all types implement the empty interface
		return true
	}

	if T.Kind == PTR {
		T = T.Elem
	}

iloop:
	for _, im := range I.MethodSet {
		for _, tm := range T.MethodSet {
			if im.Name != tm.Name {
				continue // try next tm
			}

			if !im.IsExported && im.Pkg != tm.Pkg {
				continue // try next tm
			}

			if !tm.Type.IsIdenticalTo(im.Type) {
				continue // try next tm
			}

			// tm matches im, go to next im
			continue iloop
		}

		// no method in T matched im
		return false
	}

	return true
}
