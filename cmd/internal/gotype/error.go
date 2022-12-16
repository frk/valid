package gotype

import (
	"fmt"
	"go/types"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/errors"
)

type Error struct {
	C ErrorCode
	// the identifier of the object that is associated with the error
	oid config.ObjectIdent
	// the go/types.Object associated with the error
	obj types.Object
	// the original error, if any
	err error `cmp:"+"`
}

func (e *Error) Error() string {
	return errors.String(e.C.ident(), e)
}

func (e *Error) ObjId() string {
	return e.oid.Pkg + "." + e.oid.Name
}

func (e *Error) Object() string {
	return fmt.Sprintf("%T", e.obj)
}

func (e *Error) Type() string {
	return e.obj.Type().String()
}

func (e *Error) OriginalError() string {
	return e.err.Error()
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	ERR_OBJECT_SEARCH            // searching for object failed
	ERR_ERROR_CONSTRUCTOR_OBJECT // error constructor object not *types.Func
	ERR_ERROR_CONSTRUCTOR_TYPE   // error constructor func does not have correct signature
	ERR_ERROR_AGGREGATOR_OBJECT  // error aggregator object not *types.TypeName
	ERR_ERROR_AGGREGATOR_TYPE    // error aggregator type does not implement interface
)

func (e ErrorCode) ident() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `
{{ define "` + ERR_OBJECT_SEARCH.ident() + `" -}}
{{.OriginalError -}}
{{ end }}

{{ define "` + ERR_ERROR_CONSTRUCTOR_OBJECT.ident() + `" -}}
{{ ERROR }} Cannot use "{{W .ObjId}}" (object {{.Object}}) as a global error constructor func.
{{ end }}

{{ define "` + ERR_ERROR_CONSTRUCTOR_TYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{W .ObjId}}" as a global error constructor func.
  > func has signature "{{wb .Type}}"
  > required signature "{{wb}}func(key string, val any, rule string, args ...any) error{{off}}".
{{ end }}

{{ define "` + ERR_ERROR_AGGREGATOR_OBJECT.ident() + `" -}}
{{ ERROR }} Cannot use "{{W .ObjId}}" (object {{.Object}}) as a global error aggregator type.
{{ end }}

{{ define "` + ERR_ERROR_AGGREGATOR_TYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{W .ObjId}}" as a global error aggregator type.
  > The type does NOT implement the required interface.
    {{wb}}interface {
	Error(key string, val any, rule string, args ...any)
	Out() error
    }{{off}}
{{ end }}
`
