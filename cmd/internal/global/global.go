package global

import (
	"go/types"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"
)

var (
	ErrorConstructor *gotype.Func
	ErrorAggregator  *gotype.Type
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
			return &Error{C: ERR_OBJECT_SEARCH, oid: v, err: err}
		}
		fn, ok := obj.(*types.Func)
		if !ok {
			return &Error{C: ERR_ERROR_CONSTRUCTOR_OBJECT, oid: v, obj: obj}
		}

		a := gotype.NewAnalyzer(obj.Pkg())
		t := a.Object(obj)
		if !gotype.IsErrorConstructorFunc(t) {
			return &Error{C: ERR_ERROR_CONSTRUCTOR_TYPE, oid: v, obj: obj}
		}

		ErrorConstructor = &gotype.Func{
			Name: fn.Name(),
			Type: t,
		}
	}
	if v := cfg.ErrorHandling.Aggregator; v.IsSet {
		obj, err := search.FindObject(v.Pkg, v.Name, a)
		if err != nil {
			return &Error{C: ERR_OBJECT_SEARCH, oid: v, err: err}
		}
		if _, ok := obj.(*types.TypeName); !ok {
			return &Error{C: ERR_ERROR_AGGREGATOR_OBJECT, oid: v, obj: obj}
		}

		a := gotype.NewAnalyzer(obj.Pkg())
		t := a.Object(obj)
		if !gotype.IsErrorAggregator(t) {
			return &Error{C: ERR_ERROR_AGGREGATOR_TYPE, oid: v, obj: obj}
		}

		ErrorAggregator = t
	}

	return nil
}
