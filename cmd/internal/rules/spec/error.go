package spec

import (
	"fmt"
	stdtypes "go/types"
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/errors"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"

	"gopkg.in/yaml.v3"
)

type Error struct {
	C ErrorCode

	// The AST associated with the error.
	a *search.AST
	// Set if the error is directly related to the
	// main config, otherwise nil.
	c *config.Config
	// Set if the error is related to the rule's
	// function, otherwise nil.
	ft *stdtypes.Func
	// Set if the error is directly related to the
	// config of a rule, otherwise nil.
	rs *config.RuleSpec
	// Set if the error is related to a specific
	// config argument (option), otherwise nil.
	rca *rules.Arg
	// A config argument's index. Set if the error is related
	// to a specific spec argument (option), otherwise nil.
	rcai *int
	// A config argument's option key. Set if the error is related
	// to a specific spec argument (option), otherwise nil.
	rcak *string
	// The function parameter to which the error is related, or nil.
	fp *types.Var
	// The function parameter's index with which the error is related, or nil.
	fpi *int
	// the actual error, or nil
	err error
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

func (e *Error) HasConfigFile() bool {
	return e.c != nil
}

func (e *Error) ConfigFile() string {
	return e.c.File.Value
}

func (e *Error) FuncIdent() string {
	return e.ft.Pkg().Path() + "." + e.ft.Name()
}

func (e *Error) FuncType() string {
	return e.ft.Type().String()
}

func (e *Error) FuncPos() string {
	return e.a.FileAndLine(e.ft)
}

func (e *Error) FuncParamNum() string {
	sig := e.ft.Type().(*stdtypes.Signature)
	return strconv.Itoa(sig.Params().Len() - 1)
}

func (e *Error) FuncParamWord() string {
	sig := e.ft.Type().(*stdtypes.Signature)
	if (sig.Params().Len() - 1) == 1 {
		return "parameter"
	}
	return "parameters"
}

func (e *Error) CfgRuleName() string {
	return e.rs.Name
}

func (e *Error) CfgJoinOp() string {
	return e.rs.JoinOp.String()
}

func (e *Error) CfgErrYAML() string {
	out, err := yaml.Marshal(e.rs.Error)
	if err != nil {
		return "<invalid>" // not expected to happen
	}
	return strings.TrimSpace(string(out))
}

func (e *Error) CfgArgNum() string {
	return strconv.Itoa(len(e.rs.Args))
}

func (e *Error) CfgArgWord() string {
	if len(e.rs.Args) == 1 {
		return "arg"
	}
	return "args"
}

func (e *Error) CfgArgBounds() string {
	var min, max, sep string
	if e.rs.ArgMin != nil {
		min = "ArgMin=" + strconv.FormatUint(uint64(*e.rs.ArgMin), 10)
	}
	if e.rs.ArgMax != nil {
		max = "ArgMax=" + strconv.FormatInt(int64(*e.rs.ArgMax), 10)
	}
	if len(min) > 0 && len(max) > 0 {
		sep = ", "
	}
	return min + sep + max
}

func (e *Error) CfgArgValue() string {
	return e.rca.Value
}

func (e *Error) CfgArgType() string {
	return e.rca.Type.String()
}

func (e *Error) CfgArgKey() string {
	return *e.rcak
}

func (e *Error) CfgFuncParamType() string {
	s := e.ft.Type().(*stdtypes.Signature)
	if s.Variadic() && s.Params().Len() < *e.fpi {
		return "..." + e.fp.Type.TypeString(nil)
	}
	return e.fp.Type.TypeString(nil)
}

func (e *Error) FuncParamIdent() string {
	// if it's a named parameter use the name
	// if it's unnamed, used the index
	if len(e.fp.Name) > 0 {
		return e.fp.Name
	}
	return e.FuncParamIndex()
}

func (e *Error) FuncParamIndex() string {
	// +1 change index start from 0 to 1 to make
	//    the error message more human friendly
	// +1 for the field's position
	i := *e.fpi + 1 + 1
	s := strconv.Itoa(i)

	switch i % 100 {
	case 11, 12, 13:
		s += "th"
	default:
		switch i % 10 {
		case 1:
			s += "st"
		case 2:
			s += "nd"
		case 3:
			s += "rd"
		default:
			s += "th"
		}
	}

	return s
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	E_CONFIG_FUNCID       // bad function identifier format
	E_CONFIG_FUNCSEARCH   // function search failed
	E_CONFIG_INVALID      // failed to unmarshal function's config.RuleConfig
	E_CONFIG_MISSING      // missing function's config.RuleConfig
	E_CONFIG_NONAME       // rule config with no rule name
	E_CONFIG_RESERVED     // illegal use of reserved rule name
	E_CONFIG_FUNCTYPE     // bad function signature for "is" rule
	E_CONFIG_PREFUNCTYPE  // bad function signature for "pre" rule
	E_CONFIG_PREPROCJOIN  // illegal use of cfg.JoinOp for "pre" rule
	E_CONFIG_PREPROCERROR // illegal use of cfg.Err for "pre" rule
	E_CONFIG_ARGNUM       // bad number of rule arguments
	E_CONFIG_ARGTYPE      // bad rule argument type
	E_CONFIG_ARGBOUNDS    // bad rule argument bounds
)

func (e ErrorCode) ident() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `

{{ define "` + E_CONFIG_FUNCID.ident() + `" -}}
{{ ERROR }} Invalid custom rule function identifier "{{R .CustomFuncIdent}}" in config file.
  > CONFIG: {{W .ConfigFile}}
  > HINT: A valid function identifier consists of the function's full package path followed
         by a dot (".") and the function's name, i.e. "{{W "<package_path>.<func_name>"}}".
         For example "{{W "github.com/me/mod/pkg/foo.FuncName"}}" is a valid rule function identifier.
{{ end }}

{{ define "` + E_CONFIG_FUNCSEARCH.ident() + `" -}}
{{ .OriginalError }}
  > CONFIG: {{W .ConfigFile}}
{{ end }}

{{ define "` + E_CONFIG_INVALID.ident() + `" -}}
{{ ERROR }} Failed to unmarshal custom rule config for "{{W .FuncIdent}}".
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > {{.ErrType}}: {{R (quote .OriginalError) }}
{{ end }}

{{ define "` + E_CONFIG_MISSING.ident() + `" -}}
{{ ERROR }} Missing custom rule config for function "{{W .FuncIdent}}".
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
{{ end }}

{{ define "` + E_CONFIG_NONAME.ident() + `" -}}
{{ ERROR }} Missing "{{R "name"}}" in custom rule's config.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > FUNC: {{W .FuncIdent}}
{{ end }}

{{ define "` + E_CONFIG_RESERVED.ident() + `" -}}
{{ ERROR }} Illegal use of {{wb "reserved"}} name "{{R .CfgRuleName}}" for custom rule.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > FUNC: {{W .FuncIdent}}
{{ end }}

{{ define "` + E_CONFIG_FUNCTYPE.ident() + `" -}}
{{ ERROR }} Invalid function signature {{R .FuncType}} for custom "{{wb .CfgRuleName}}" rule.
  > FILE: {{W .FuncPos}}
  > FUNC: {{W .FuncIdent}}
  > HINT: A custom rule function MUST have {{wb "at least one"}} parameter value and it MUST have either:
	- {{wb "one"}} result value of type {{wb "bool"}}, e.g. {{wb "func(v string) (ok bool)"}}.
	- or {{wb "two"}} result values where the first one is of type {{wb "bool"}} and the second one
	  is of type {{wb "error"}}, e.g. {{wb "func(v string) (ok bool, err error)"}}.
{{ end }}

{{ define "` + E_CONFIG_PREFUNCTYPE.ident() + `" -}}
{{ ERROR }} Invalid function signature {{R .FuncType}} for custom pre-processor "{{wb .CfgRuleName}}" rule.
  > FILE: {{W .FuncPos}}
  > FUNC: {{W .FuncIdent}}
  > HINT: A custom pre-processor function MUST have {{wb "at least one"}} parameter value and
          it MUST have {{wb "exactly one"}} result value which MUST be of a type {{wb "identical"}} to
          the function's {{wb "first"}} parameter type.
{{ end }}

{{ define "` + E_CONFIG_ARGNUM.ident() + `" -}}
{{ ERROR }} Incompatible number of args in "{{wb .CfgRuleName}}" rule config for "{{W .FuncIdent}}".
  Config specifies {{rb .CfgArgNum }} {{.CfgArgWord}}. Function specifies {{wb .FuncParamNum }} ` +
	`{{.FuncParamWord}} (not counting the primary parameter).
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ end -}}
  > FILE: {{W .FuncPos}}
{{ end }}

{{ define "` + E_CONFIG_ARGTYPE.ident() + `" -}}
{{ ERROR }} Incompatible arg type in "{{wb .CfgRuleName}}" config for "{{W .FuncIdent}}".
  Config specifies {{with .CfgArgKey}}option "{{R .}}":{{else}}default option {{end}}` +
	`"{{R .CfgArgValue}}" (type {{R .CfgArgType}}) as the argument value for the function's ` +
	`{{wb .FuncParamIdent}} parameter (type {{wb .CfgFuncParamType}}).
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ end -}}
  > FILE: {{W .FuncPos}}
  > HINT: The arguments in the config MUST be {{wb "assignable"}} to the function's` +
	` corresponding parameters.
{{ end }}

{{ define "` + E_CONFIG_ARGBOUNDS.ident() + `" -}}
{{ ERROR }} Incompatible arg bounds in rule config for "{{W .FuncIdent}}".
  Config specifies {{rb .CfgArgBounds }} bounds. Function specifies {{wb .FuncParamNum }} ` +
	`{{.FuncParamWord}} (not counting the primary parameter).
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ end -}}
  > FILE: {{W .FuncPos}}
{{ end }}

{{ define "` + E_CONFIG_PREPROCJOIN.ident() + `" -}}
{{ ERROR }} Illegal use of {{R "JoinOp"}}:{{R .CfgJoinOp}} in config for the custom ` +
	`pre-processor "{{wb .CfgRuleName}}" rule.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > HINT: Pre-processor rules DO NOT support joins, therefore the JoinOp SHOULD ` +
	`be omitted from the rule's config.
{{ end }}

{{ define "` + E_CONFIG_PREPROCERROR.ident() + `" -}}
{{ ERROR }} Illegal use of {{R "Err"}}:{{R .CfgErrYAML}} in config for the custom ` +
	`pre-processor "{{wb .CfgRuleName}}" rule.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > HINT: Pre-processor rules DO NOT support errors, therefore the Err SHOULD ` +
	`be omitted from the rule's config.
{{ end }}

` //`

////////////////////////////////////////////////////////////////////////////////
// helpers
////////////////////////////////////////////////////////////////////////////////

func extendError(err error, x func(*Error)) error {
	if e, ok := err.(*Error); ok {
		x(e)
	}
	return err
}
