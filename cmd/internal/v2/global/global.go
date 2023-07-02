package global

import (
	go_types "go/types"

	"github.com/frk/valid/cmd/internal/v2/config"
	"github.com/frk/valid/cmd/internal/v2/source"
	"github.com/frk/valid/cmd/internal/v2/types"
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

func Init(cfg config.Config, src *source.Source) error {
	if v := cfg.ErrorHandling.Constructor; v.IsSet {
		obj, err := src.FindObject(v.Pkg, v.Name)
		if err != nil {
			return &Error{C: E_OBJECT_SEARCH, oid: v, err: err}
		}
		fn, ok := obj.(*go_types.Func)
		if !ok {
			return &Error{C: E_ERROR_CONSTRUCTOR_OBJECT, oid: v, obj: obj}
		}

		f := types.AnalyzeFunc(fn, src)
		if !f.Type.IsErrorConstructorFunc() {
			return &Error{C: E_ERROR_CONSTRUCTOR_TYPE, oid: v, obj: obj}
		}
		ErrorConstructor = f
	}
	if v := cfg.ErrorHandling.Aggregator; v.IsSet {
		obj, err := src.FindObject(v.Pkg, v.Name)
		if err != nil {
			return &Error{C: E_OBJECT_SEARCH, oid: v, err: err}
		}
		tn, ok := obj.(*go_types.TypeName)
		if !ok {
			return &Error{C: E_ERROR_AGGREGATOR_OBJECT, oid: v, obj: obj}
		}
		named, ok := tn.Type().(*go_types.Named)
		if !ok {
			return &Error{C: E_ERROR_AGGREGATOR_OBJECT, oid: v, obj: obj}
		}

		t := types.Analyze(named, src)
		if !t.IsErrorAggregator() {
			return &Error{C: E_ERROR_AGGREGATOR_TYPE, oid: v, obj: obj}
		}

		ErrorAggregator = t
	}

	return nil
}
