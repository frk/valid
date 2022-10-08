package gotype

import (
	"go/types"
	"strings"
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
// corresponding gotype.Validator representation.
//
// The method will panic if the named type is not a struct.
func (an *Analyzer) Validator(named *types.Named) *Validator {
	if _, ok := named.Underlying().(*types.Struct); !ok {
		panic(named.Obj().Name() + " must be a struct type.")
	}
	if name := named.Obj().Name(); !strings.HasSuffix(strings.ToLower(name), "validator") {
		panic(name + " validator struct type has unsupported name suffix.")
	}

	v := new(Validator)
	v.Type = an.Analyze(named)
	if name := nameOfValidateMethod(v.Type, "before"); len(name) > 0 {
		v.BeforeValidateMethod = &MethodInfo{Name: name}
	}
	if name := nameOfValidateMethod(v.Type, "after"); len(name) > 0 {
		v.AfterValidateMethod = &MethodInfo{Name: name}
	}
	for _, f := range v.Type.Fields {
		if IsErrorConstructor(f.Type) {
			v.ErrorHandlerField = new(ErrorHandlerField)
			v.ErrorHandlerField.Name = f.Name
		} else if IsErrorAggregator(f.Type) {
			v.ErrorHandlerField = new(ErrorHandlerField)
			v.ErrorHandlerField.Name = f.Name
			v.ErrorHandlerField.IsAggregator = true
		}
	}
	return v
}
