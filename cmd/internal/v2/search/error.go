package search

import (
	"fmt"
	"strings"

	"github.com/frk/valid/cmd/internal/v2/errors"
)

type Error struct {
	C ErrorCode

	pkg  string
	name string
	err  error `cmp:"+"`
}

func (e *Error) Error() string {
	return errors.String(e.C.ident(), e)
}

func (e *Error) OriginalError() string {
	return strings.TrimRight(e.err.Error(), "\n")
}

func (e *Error) ErrType() string {
	return strings.TrimLeft(fmt.Sprintf("%T", e.err), "*")
}

func (e *Error) ObjName() string {
	return e.name
}

func (e *Error) PkgPath() string {
	return e.pkg
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	ERR_OBJECT_NOTFOUND // object (func or type) not found
	ERR_FUNC_NOTFOUND   // func not found
	ERR_PKG_NOTFOUND    // package not found
	ERR_PKG_LOADFAIL    // failed loading package
	ERR_PKG_ERROR       // package contains errors
)

func (e ErrorCode) ident() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `
{{ define "` + ERR_OBJECT_NOTFOUND.ident() + `" -}}
{{ ERROR }} Could not find top-level func nor type named "{{W .ObjName}}" in package "{{W .PkgPath}}".
{{ end }}

{{ define "` + ERR_FUNC_NOTFOUND.ident() + `" -}}
{{ ERROR }} Could not find function "{{W .ObjName}}" in package "{{W .PkgPath}}".
{{ end }}

{{ define "` + ERR_PKG_NOTFOUND.ident() + `" -}}
{{ ERROR }} Could not find package "{{W .PkgPath}}" for function "{{W .ObjName}}".
{{ end }}

{{ define "` + ERR_PKG_LOADFAIL.ident() + `" -}}
{{ ERROR }} Failed to load package "{{W .PkgPath}}" for function "{{W .ObjName}}":
  > {{.ErrType}}: {{Y (quote .OriginalError)}}
{{ end }}

{{ define "` + ERR_PKG_ERROR.ident() + `" -}}
{{ ERROR }} Loading of package "{{W .PkgPath}}" for function "{{W .ObjName}}" returned errors:
  > {{.ErrType}}: {{Y (quote .OriginalError)}}
{{ end }}
`
