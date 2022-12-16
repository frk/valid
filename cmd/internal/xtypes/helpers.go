package xtypes

import (
	"strings"
)

// IsErrorAggregator reports whether or not the given
// type implements the "ErrorAggregator" interface.
func IsErrorAggregator(t *Type) bool {
	var hasAddError, hasError bool
	for _, m := range t.Methods {
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
func IsErrorConstructor(t *Type) bool {
	var hasError bool
	for _, m := range t.Methods {
		switch m.Name {

		// Error(key string, val any, rule string, args ...any) error
		case "Error":
			if IsErrorConstructorFunc(m.Type) {
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
func IsErrorConstructorFunc(t *Type) bool {
	if t.Kind != K_FUNC {
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

// nameOfValidateMethod scans the method set of the given type for
// a method with the signature "func() error" and with its name equal
// to the given prefix+"validate" (case insensitive). If it finds a
// match it will return that method's full name (with case preserved),
// and if there's no match it will it will return an empty string.
func nameOfValidateMethod(t *Type, prefix string) string {
	name := strings.ToLower(prefix) + "validate"
	for _, m := range t.Methods {
		if strings.ToLower(m.Name) == name {
			sig := m.Type // signature
			in, out := sig.In, sig.Out
			if len(in) != 0 || len(out) != 1 {
				return ""
			}
			if !out[0].Type.IsGoError() {
				return ""
			}
			return m.Name
		}
	}
	return ""
}
