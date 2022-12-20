package types

// Reports whether or not the type's kind is is one of the provided kinds.
func (t Type) Is(kinds ...Kind) bool {
	for _, k := range kinds {
		if t.Kind == k {
			return true
		}
	}
	return false
}

// IsNamed reports whether or not the type has a name.
func (t *Type) IsNamed() bool {
	return t.Name != ""
}

// IsUnnamed is identical to !T.IsNamed().
func (t *Type) IsUnnamed() bool {
	return !t.IsNamed()
}

func (t *Type) IsBytesOrRunes() bool {
	return t.Kind == SLICE &&
		(t.Elem.Type.Kind == BYTE || t.Elem.Type.Kind == RUNE)
}

// Indicates whether or not the type is an empty interface type.
func (t *Type) IsEmptyIface() bool {
	return t.Kind == INTERFACE && len(t.MethodSet) == 0
}

// IsEmptyIfaceSlice reports whether or not t is the Go builtin []interface{} type.
func (t *Type) IsEmptyIfaceSlice() bool {
	if t.Kind == SLICE {
		return t.Elem.Type.IsEmptyIfaceSlice()
	}
	return false
}

// IsIncluded reports whether or not the Type was
// declared in the github.com/frk/valid package.
func (t *Type) IsIncluded() bool {
	return t.Pkg.Path == "github.com/frk/valid"
}

// IsGoString reports whether or not t is the Go builtin string type.
func (t *Type) IsGoString() bool {
	return t.Pkg == Pkg{} && t.Kind == STRING && t.Name == ""
}

// IsGoError reports whether or not t is the Go builtin error type.
func (t *Type) IsGoError() bool {
	return t.Pkg == Pkg{} && t.Kind == INTERFACE && t.Name == "error"
}

// IsGoAny reports whether or not t is the Go builtin any/interface{} type.
func (t *Type) IsGoAny() bool {
	return t.Pkg == Pkg{} && t.Kind == INTERFACE &&
		(t.Name == "any" || len(t.MethodSet) == 0)
}

// IsGoAnySlice reports whether or not t is the Go builtin []any/[]interface{} type.
func (t *Type) IsGoAnySlice() bool {
	if t.Kind == SLICE {
		return t.Elem.Type.IsGoAny()
	}
	return false
}

// IsNilable reports wether or not a value of the Go
// type represented by t can be set to nil.
func (t *Type) IsNilable() bool {
	return t.Kind == PTR ||
		t.Kind == SLICE ||
		t.Kind == MAP ||
		t.Kind == INTERFACE ||
		t.Kind == FUNC ||
		t.Kind == CHAN
}

// HasLength reports whether or not the Go type
// represented by t has a length.
func (t *Type) HasLength() bool {
	return t.Kind == STRING ||
		t.Kind == ARRAY ||
		t.Kind == SLICE ||
		t.Kind == MAP ||
		t.Kind == CHAN
}

// IsValid reports whether or not the "IsValid() bool"
// method belongs to the method set of the type t.
func (t *Type) HasIsValid() bool {
	for _, m := range t.MethodSet {
		if m.Name == "IsValid" &&
			len(m.Type.In) == 0 &&
			len(m.Type.Out) == 1 &&
			m.Type.Out[0].Type.Kind == BOOL {
			return true
		}
	}
	return false
}

// CanError reports that the type, if it *is* a FUNC type,
// has error as its last return value type.
func (t *Type) CanError() bool {
	if t.Kind != FUNC {
		return false
	}
	if n := len(t.Out); n > 0 && t.Out[n-1].Type.IsGoError() {
		return true
	}
	return false
}

// IsComparable reports wether or not a value of the Go
// type represented by t is comparable.
func (t *Type) IsComparable() bool {
	if t.Kind == MAP || t.Kind == SLICE || t.Kind == FUNC {
		return false
	}
	if t.Kind == ARRAY {
		return t.Elem.Type.IsComparable()
	}
	if t.Kind == STRUCT {
		for _, f := range t.Fields {
			if !f.Obj.Type.IsComparable() {
				return false
			}
		}
	}
	return true
}

// IsErrorAggregator reports whether or not the given
// type implements the "ErrorAggregator" interface.
func (t *Type) IsErrorAggregator() bool {
	var hasAddError, hasError bool
	for _, m := range t.MethodSet {
		switch m.Name {

		// Error(key string, val any, rule string, args ...any)
		case "Error":
			sig := m.Type // signature
			if !sig.IsVariadic {
				return false
			}
			in, out := sig.In, sig.Out
			if len(in) != 4 || len(out) != 0 {
				return false
			}
			if !in[0].Type.IsGoString() || !in[1].Type.IsGoAny() ||
				!in[2].Type.IsGoString() || !in[3].Type.IsGoAnySlice() {
				return false
			}
			hasAddError = true

		// Out() error
		case "Out":
			sig := m.Type // signature
			in, out := sig.In, sig.Out
			if len(in) != 0 || len(out) != 1 {
				return false
			}
			if !out[0].Type.IsGoError() {
				return false
			}
			hasError = true
		}
	}
	return hasAddError && hasError
}

// IsErrorConstructor reports whether or not the given
// type implements the "ErrorConstructor" interface.
func (t *Type) IsErrorConstructor() bool {
	var hasError bool
	for _, m := range t.MethodSet {
		switch m.Name {

		// Error(key string, val any, rule string, args ...any) error
		case "Error":
			if m.Type.IsErrorConstructorFunc() {
				hasError = true
			}
		}
	}
	return hasError
}

// IsErrorConstructorFunc reports whether or not the given
// type matches the "ErrorConstructor" function signature.
//
//	func(key string, val any, rule string, args ...any) error
//
func (t *Type) IsErrorConstructorFunc() bool {
	if t.Kind != FUNC {
		return false
	}
	if !t.IsVariadic {
		return false
	}

	in, out := t.In, t.Out
	if len(in) != 4 || len(out) != 1 {
		return false
	}
	if !in[0].Type.IsGoString() || !in[1].Type.IsGoAny() ||
		!in[2].Type.IsGoString() || !in[3].Type.IsGoAnySlice() ||
		!out[0].Type.IsGoError() {
		return false
	}
	return true
}
