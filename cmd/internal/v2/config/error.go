package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/v2/errors"

	"gopkg.in/yaml.v3"
)

type Error struct {
	C ErrorCode

	// the working directory
	dir string
	// the file path of the config file
	file string
	// the target type for which decoding failed
	tt interface{}
	// the key of the failed config value
	key string
	// the failed config value
	val string
	// the failed yaml node
	node *yaml.Node `cmp:"+"`
	// the original error
	err error `cmp:"+"`
}

func (e *Error) Error() string {
	return errors.String(e.C.Id(), e)
}

func (e *Error) WDir() string {
	return e.dir
}

func (e *Error) File() string {
	return e.file
}

func (e *Error) FilePos() string {
	if e.file != "" && e.node != nil {
		return e.file + ":" + strconv.Itoa(e.node.Line)
	}
	return e.file
}

func (e *Error) OriginalError() string {
	if e.err != nil {
		return strings.TrimRight(e.err.Error(), "\n")
	}
	return ""
}

func (e *Error) ErrType() string {
	return strings.TrimLeft(fmt.Sprintf("%T", e.err), "*")
}

func (e *Error) TargetType() string {
	return strings.TrimLeft(fmt.Sprintf("%T", e.tt), "*")
}

func (e *Error) Key() string {
	return e.key
}

func (e *Error) Value() string {
	return e.val
}

func (e *Error) YAMLNodeLine() string {
	return strconv.Itoa(e.node.Line)
}

func (e *Error) YAMLNode() string {
	n := e.node.Tag
	if e.node.Kind == yaml.ScalarNode {
		n = fmt.Sprintf("%q", e.node.Value)
	}
	return fmt.Sprintf("YAML %s on line %d", n, e.node.Line)
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	ERR_FLAG_DECODE    // flag decoding failed
	ERR_YAML_ERROR     // error from gopkg.in/yaml.v3
	ERR_YAML_TYPE      // unsupported yaml type
	ERR_YAML_FILE      // yaml file is invalid
	ERR_OBJECT_IDENT   // invalid format for object identifier
	ERR_JOIN_OP        // invalid join op value
	ERR_CONFIG_FILE    // config file unusable
	ERR_WORK_DIR       // working_directory unusable
	ERR_FILE_ITEM      // file_list item unusable
	ERR_OUTNAME_FORMAT // invalid output name format
	ERR_PATTERN        // invalid regular expression
	ERR_FKEY_TAG       // invalid field key tag
	ERR_FKEY_SEP       // invalid field key separator
	ERR_RULE_NONAME    // rule with no name
	ERR_RULE_NOFUNC    // missing rule func
	ERR_RULE_DUPNAME   // duplicate rule name
	ERR_TYPE_NONAME    // type with no name
	ERR_TYPE_DUPNAME   // duplicate type name
)

func (e ErrorCode) Id() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `
{{ define "config_error_meta" }}
{{- with .FilePos }}
  > FILE: {{W . }}
{{- end }}
{{- with .WDir }}
  > WDIR: {{W . }}
{{- end }}
{{- with .OriginalError }}
  > {{$.ErrType}}: {{Y (quote .)}}
{{- end }}
{{ end }}

{{ define "` + ERR_FLAG_DECODE.Id() + `" -}}
{{ ERRCFG }} Decoding flag value "{{W .Value}}" into "{{W .TargetType}}" resulted in an error.
  > {{.ErrType}}: {{Y (quote .OriginalError)}}
{{ end }}

{{ define "` + ERR_YAML_ERROR.Id() + `" -}}
{{ ERRCFG }} Unmarshaling YAML on line {{.YAMLNodeLine}} into "{{W .TargetType}}" resulted in an error.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_YAML_TYPE.Id() + `" -}}
{{ ERRCFG }} Cannot unmarshal {{.YAMLNode}} into type "{{W .TargetType}}".
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_YAML_FILE.Id() + `" -}}
{{ ERRCFG }} Trying to unmarshal config file resulted in an error.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_OBJECT_IDENT.Id() + `" -}}
{{ ERRCFG }} Invalid {{Wi "object_identifier"}} format: {{(quote .Value)}}
{{- template "config_error_meta" . -}}
{{""}}  > HINT: An object_identifier string MUST consist of a Go package's import path followed` +
	`{{NT}}  by a dot (".") and the name of a type or function that is declared in that package.
{{ end }}

{{ define "` + ERR_JOIN_OP.Id() + `" -}}
{{ ERRCFG }} Invalid {{Wi "join_op"}} value: {{(quote .Value)}}
{{- template "config_error_meta" . -}}
{{""}}  > HINT: A join_op value MUST be one of the following: "AND", "OR", "NOT".
{{ end }}

{{ define "` + ERR_CONFIG_FILE.Id() + `" -}}
{{ ERRCFG }} Trying to resolve the config file resulted in an error:
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_WORK_DIR.Id() + `" -}}
{{ ERRCFG }} Trying to use the working directory resulted in an error:
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_FILE_ITEM.Id() + `" -}}
{{ ERRCFG }} Unable to use the provided file "{{W .Value}}".
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_OUTNAME_FORMAT.Id() + `" -}}
{{ ERRCFG }} The value "{{W .Value}}" is not a valid "{{W .Key}}".
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_PATTERN.Id() + `" -}}
{{ ERRCFG }} The value "{{W .Value}}" is not a valid "{{W .Key}}" (regular expression).
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_FKEY_TAG.Id() + `" -}}
{{ ERRCFG }} The value "{{W .Value}}" is not a valid "{{W .Key}}".
{{- template "config_error_meta" . -}}
{{""}}  > HINT: A field key tag MUST match the "{{W ` + "`^(?:[A-Za-z_]\\w*)?$`" + `}}" regular expression.
{{ end }}

{{ define "` + ERR_FKEY_SEP.Id() + `" -}}
{{ ERRCFG }} The value "{{W .Value}}" is not a valid "{{W .Key}}".
{{- template "config_error_meta" . -}}
{{""}}  > HINT: A field key separator MUST be a single byte.
{{ end }}

{{ define "` + ERR_RULE_NONAME.Id() + `" -}}
{{ ERRCFG }} The rule (index {{W .Value}}) is missing a name.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_RULE_NOFUNC.Id() + `" -}}
{{ ERRCFG }} The rule (index {{W .Value}}) has no associated func value.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_RULE_DUPNAME.Id() + `" -}}
{{ ERRCFG }} The rule name {{W .Value}} is already taken by another rule.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_TYPE_NONAME.Id() + `" -}}
{{ ERRCFG }} The type (index {{W .Value}}) is missing a name.
{{- template "config_error_meta" . -}}
{{ end }}

{{ define "` + ERR_TYPE_DUPNAME.Id() + `" -}}
{{ ERRCFG }} The type name {{W .Value}} is already taken by another type.
{{- template "config_error_meta" . -}}
{{ end }}

`
