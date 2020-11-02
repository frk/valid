package analysis

import (
	"go/types"
)

// methoder represents a type with methods. It is implicitly implemented
// by *types.Interface and *types.Named.
type methoder interface {
	NumMethods() int
	Method(i int) *types.Func
}

// isString reports whether or not the given type is the builtin "string" type.
func isString(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.String && basic.Name() == "string"
}

// isBool reports whether or not the given type is the builtin "bool" type.
func isBool(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Bool && basic.Name() == "bool"
}

// isEmptyInterface reports whether or not the given type is the "interface{}" type.
func isEmptyInterface(typ types.Type) bool {
	iface, ok := typ.(*types.Interface)
	if !ok {
		return false
	}
	return iface.NumMethods() == 0
}

// isEmptyInterfaceSlice reports whether or not the given type is the "[]interface{}" type.
func isEmptyInterfaceSlice(typ types.Type) bool {
	if s, ok := typ.(*types.Slice); ok {
		return isEmptyInterface(s.Elem())
	}
	return false
}

// isError reports whether or not the given type is the "error" type.
func isError(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	pkg := named.Obj().Pkg()
	name := named.Obj().Name()
	return pkg == nil && name == "error"
}

// isErrorConstructor reports whether or not the given type implements the "isvalid.ErrorConstructor" interface.
func isErrorConstructor(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	mm := methoder(named)
	if iface, ok := named.Underlying().(*types.Interface); ok {
		mm = iface
	}

	var hasError bool
	for i := 0; i < mm.NumMethods(); i++ {
		m := mm.Method(i)
		switch m.Name() {
		case "Error": // Error(key string, val interface{}, rule string, args ...interface{}) error
			sig := m.Type().(*types.Signature)
			if !sig.Variadic() {
				return false
			}
			p, r := sig.Params(), sig.Results()
			if p.Len() != 4 || r.Len() != 1 {
				return false
			}
			if !isString(p.At(0).Type()) || !isEmptyInterface(p.At(1).Type()) ||
				!isString(p.At(2).Type()) || !isEmptyInterfaceSlice(p.At(3).Type()) ||
				!isError(r.At(0).Type()) {
				return false
			}
			hasError = true
		}
	}
	return hasError
}

// isErrorAggregator reports whether or not the given type implements the "isvalid.ErrorAggregator" interface.
func isErrorAggregator(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	mm := methoder(named)
	if iface, ok := named.Underlying().(*types.Interface); ok {
		mm = iface
	}

	var hasAddError, hasError bool
	for i := 0; i < mm.NumMethods(); i++ {
		m := mm.Method(i)
		switch m.Name() {
		case "Error": // Error(key string, val interface{}, rule string, args ...interface{})
			sig := m.Type().(*types.Signature)
			if !sig.Variadic() {
				return false
			}
			p, r := sig.Params(), sig.Results()
			if p.Len() != 4 || r.Len() != 0 {
				return false
			}
			if !isString(p.At(0).Type()) || !isEmptyInterface(p.At(1).Type()) ||
				!isString(p.At(2).Type()) || !isEmptyInterfaceSlice(p.At(3).Type()) {
				return false
			}
			hasAddError = true
		case "Out": // Out() error
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 0 || r.Len() != 1 {
				return false
			}
			if !isError(r.At(0).Type()) {
				return false
			}
			hasError = true
		}
	}
	return hasAddError && hasError
}
