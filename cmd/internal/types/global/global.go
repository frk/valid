package global

import (
	stdtypes "go/types"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"
)

var (
	ErrorConstructor *types.Func
	ErrorAggregator  *types.Type
)

// used by tests
func Unset() {
	ErrorConstructor = nil
	ErrorAggregator = nil
}

//
func Init(cfg config.Config, a *search.AST) error {
	if v := cfg.ErrorHandling.Constructor; v.IsSet {
		obj, err := search.FindObject(v.Pkg, v.Name, a)
		if err != nil {
			return &Error{C: E_OBJECT_SEARCH, oid: v, err: err}
		}
		fn, ok := obj.(*stdtypes.Func)
		if !ok {
			return &Error{C: E_ERROR_CONSTRUCTOR_OBJECT, oid: v, obj: obj}
		}

		f := types.AnalyzeFunc(fn, a)
		if !f.Type.IsErrorConstructorFunc() {
			return &Error{C: E_ERROR_CONSTRUCTOR_TYPE, oid: v, obj: obj}
		}
		ErrorConstructor = f
	}
	if v := cfg.ErrorHandling.Aggregator; v.IsSet {
		obj, err := search.FindObject(v.Pkg, v.Name, a)
		if err != nil {
			return &Error{C: E_OBJECT_SEARCH, oid: v, err: err}
		}
		tn, ok := obj.(*stdtypes.TypeName)
		if !ok {
			return &Error{C: E_ERROR_AGGREGATOR_OBJECT, oid: v, obj: obj}
		}
		named, ok := tn.Type().(*stdtypes.Named)
		if !ok {
			return &Error{C: E_ERROR_AGGREGATOR_OBJECT, oid: v, obj: obj}
		}

		t := types.Analyze(named, a)
		if !t.IsErrorAggregator() {
			return &Error{C: E_ERROR_AGGREGATOR_TYPE, oid: v, obj: obj}
		}

		ErrorAggregator = t
	}

	return nil
}
