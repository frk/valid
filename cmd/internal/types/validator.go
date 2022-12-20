package types

import (
	"go/types"
	"strings"

	"github.com/frk/valid/cmd/internal/search"
)

// Validator represents a validator struct type.
type Validator struct {
	// The struct type info.
	Type *Type
	// Info on the valid.ErrorConstructor or valid.ErrorAggregator
	// field of the validator struct type, or nil.
	ErrorHandlerField *ErrorHandlerField
	// Info on the validator type's method named "beforevalidate" (case insensitive), or nil.
	BeforeValidateMethod *MethodInfo
	// Info on the validator type's method named "aftervalidate" (case insensitive), or nil.
	AfterValidateMethod *MethodInfo
}

// ErrorHandlerField is the result of analyzing a validator struct's
// field whose type implements the valid.ErrorConstructor
// or valid.ErrorAggregator interface.
type ErrorHandlerField struct {
	// Name of the field (case preserved).
	Name string
	// Indicates whether or not the field's type implements
	// the valid.ErrorAggregator interface.
	IsAggregator bool
}

// MethodInfo represents the result of analysing a type's method.
type MethodInfo struct {
	// The name of the method (case preserved).
	Name string
}

// Validator analyzes the named type and returns its
// corresponding types.Validator representation.
//
// The method will panic if the named type is not a struct.
func AnalyzeValidator(named *types.Named, ast *search.AST) *Validator {
	if _, ok := named.Underlying().(*types.Struct); !ok {
		panic(named.Obj().Name() + " must be a struct type.")
	}
	if name := named.Obj().Name(); !strings.HasSuffix(strings.ToLower(name), "validator") {
		panic(name + " validator struct type has unsupported name suffix.")
	}

	v := new(Validator)
	v.Type = Analyze(named, ast)
	if name := nameOfValidateMethod(v.Type, "before"); len(name) > 0 {
		v.BeforeValidateMethod = &MethodInfo{Name: name}
	}
	if name := nameOfValidateMethod(v.Type, "after"); len(name) > 0 {
		v.AfterValidateMethod = &MethodInfo{Name: name}
	}
	for _, f := range v.Type.Fields {
		if f.Obj.Type.IsErrorConstructor() {
			v.ErrorHandlerField = new(ErrorHandlerField)
			v.ErrorHandlerField.Name = f.Name
		} else if f.Obj.Type.IsErrorAggregator() {
			v.ErrorHandlerField = new(ErrorHandlerField)
			v.ErrorHandlerField.Name = f.Name
			v.ErrorHandlerField.IsAggregator = true
		}
	}
	return v
}

// nameOfValidateMethod scans the method set of the given type for
// a method with the signature "func() error" and with its name equal
// to the given prefix+"validate" (case insensitive). If it finds a
// match it will return that method's full name (with case preserved),
// and if there's no match it will it will return an empty string.
func nameOfValidateMethod(t *Type, prefix string) string {
	name := strings.ToLower(prefix) + "validate"
	for _, m := range t.MethodSet {
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
