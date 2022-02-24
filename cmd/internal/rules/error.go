package rules

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/errors"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"

	"gopkg.in/yaml.v3"
)

type Error struct {
	C ErrorCode

	// The AST associated with the error.
	a *search.AST
	// Set if the error is directly related to the
	// main config, otherwise nil.
	c *config.Config
	// Set if the error is directly related to the
	// custom rule config, otherwise nil.
	rc *config.RuleConfig
	// Set if the error is directly related to the
	// config of a rule, otherwise nil.
	rs *config.RuleSpec
	// Set if the error is related to a specific
	// config argument (option), otherwise nil.
	rca *Arg
	// A config argument's index. Set if the error is related
	// to a specific spec argument (option), otherwise nil.
	rcai *int
	// A config argument's option key. Set if the error is related
	// to a specific spec argument (option), otherwise nil.
	rcak *string
	// Set if the error is related to the rule's
	// function, otherwise nil.
	ft *types.Func
	// Set if the error is related to a specific
	// struct field, otherwise nil.
	sf *gotype.StructField
	// The *types.Var instance associated with
	// the struct field.
	sfv *types.Var
	// Set if the error is directly related to
	// a specific type, otherwise nil.
	ty *gotype.Type
	// Set if the error is related to a tag,
	// otherwise nil.
	tag *Tag
	// Set if the error is related to a specific
	// rule, otherwise nil.
	r *Rule
	// Set if the error is related to a combination
	// of two rules, otherwise nil.
	r2 *Rule
	// Set if the error is related to a specific
	// rule argument, otherwise nil.
	ra *Arg
	// Set if the error is related to a specific
	// rule argument and that argument is a field
	// reference, otherwise nil.
	raf *gotype.StructField
	// A function parameter. Set if the error is related
	// to a specific rule function parameter, otherwise nil.
	fp *gotype.Var
	// A function parameter's index. Set if the error is related
	// to a specific rule function parameter, otherwise nil.
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

func (e *Error) CustomFuncIdent() string {
	return e.rc.Func.String()
}

func (e *Error) FuncPos() string {
	return e.a.FileAndLine(e.ft)
}

func (e *Error) FuncName() string {
	return e.ft.Name()
}

func (e *Error) FuncIdent() string {
	return e.ft.Pkg().Path() + "." + e.ft.Name()
}

func (e *Error) FuncType() string {
	return e.ft.Type().String()
}

func (e *Error) FuncParamNum() string {
	sig := e.ft.Type().(*types.Signature)
	return strconv.Itoa(sig.Params().Len() - 1)
}

func (e *Error) FuncParamWord() string {
	sig := e.ft.Type().(*types.Signature)
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
	return string(out)
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
	s := e.ft.Type().(*types.Signature)
	if s.Variadic() && s.Params().Len() < *e.fpi {
		return "..." + e.fp.Type.String()
	}
	return e.fp.Type.String()
}

func (e *Error) FieldPos() string {
	return e.a.FileAndLine(e.sfv)
}

func (e *Error) Field() string {
	return fmt.Sprintf("%s %s `%s`", e.sf.Name, e.sf.Type, e.sf.Tag)
}

func (e *Error) FieldName() string {
	return e.sf.Name
}

func (e *Error) FieldType() string {
	return e.sf.Type.String()
}

func (e *Error) FieldTypeFamily() string {
	switch {
	case e.sf.Type.Kind.IsFloat():
		return "float"
	case e.sf.Type.Kind.IsUnsigned():
		return "uint"
	case e.sf.Type.Kind.IsInteger():
		return "int"
	}
	return e.sf.Type.String()
}

func (e *Error) Tag() string {
	return e.tag.String()
}

func (e *Error) TagKey() string {
	return "[" + e.tag.Key.String() + "]"
}

func (e *Error) TagElem() string {
	return "[]" + e.tag.Elem.String()
}

func (e *Error) TagSTKey() string {
	return e.tag.stkey + ":"
}

func (e *Error) Type() string {
	return e.ty.String()
}

func (e *Error) Rule() string {
	return e.r.String()
}

func (e *Error) RuleName() string {
	return e.r.Name
}

func (e *Error) Rule2Name() string {
	return e.r2.Name
}

func (e *Error) RuleArgs() string {
	vv := make([]string, len(e.r.Args))
	for i := range e.r.Args {
		vv[i] = e.r.Args[i].Value
	}
	return strings.Join(vv, ":")
}

func (e *Error) RuleFuncIdent() string {
	return e.r.Spec.FType.Pkg.Name + "." + e.r.Spec.FName
}

func (e *Error) RuleFuncType() string {
	return e.r.Spec.FType.String()
}

func (e *Error) RuleFuncIn0Type() string {
	return e.r.Spec.FType.In[0].Type.String()
}

func (e *Error) RuleFuncOut0Type() string {
	return e.r.Spec.FType.Out[0].Type.String()
}

func (e *Error) RuleFuncParamType() string {
	if e.r.Spec.FType.IsVariadic && len(e.r.Spec.FType.In) < *e.fpi {
		return "..." + e.fp.Type.String()
	}
	return e.fp.Type.String()
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

func (e *Error) RuleArgValue() string {
	return e.ra.Value
}

func (e *Error) RuleArgType() string {
	if e.ra.Type == ARG_FIELD {
		return e.raf.Type.String()
	}
	return e.ra.Type.String()
}

func (e *Error) RuleArgIsField() bool {
	return e.ra.Type == ARG_FIELD
}

func (e *Error) RuleArgFieldType() string {
	return e.raf.Type.String()
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	ERR_CONFIG_FUNCID       // bad function identifier format
	ERR_CONFIG_FUNCSEARCH   // function search failed
	ERR_CONFIG_INVALID      // failed to unmarshal function's config.RuleConfig
	ERR_CONFIG_MISSING      // missing function's config.RuleConfig
	ERR_CONFIG_NONAME       // rule config with no rule name
	ERR_CONFIG_RESERVED     // illegal use of reserved rule name
	ERR_CONFIG_FUNCTYPE     // bad function signature for "is" rule
	ERR_CONFIG_PREFUNCTYPE  // bad function signature for "pre" rule
	ERR_CONFIG_PREPROCJOIN  // illegal use of cfg.JoinOp for "pre" rule
	ERR_CONFIG_PREPROCERROR // illegal use of cfg.Err for "pre" rule
	ERR_CONFIG_ARGNUM       // bad number of rule arguments
	ERR_CONFIG_ARGTYPE      // bad rule argument type
	ERR_CONFIG_ARGBOUNDS    // bad rule argument bounds

	ERR_RULE_UNDEFINED // illegal use of undefined rule
	ERR_RULE_KEY       // illegal rule key on non-map field
	ERR_RULE_ELEM      // illegal rule elem on non-map/non-slice/non-array field

	ERR_NOTNIL_TYPE // illegal rule "notnil" on non-nilable field

	ERR_OPTIONAL_CONFLICT // an optional rule is in conflict with a required rule

	ERR_ENUM_NONAME  // illegal rule "enum" on field with unnamed type
	ERR_ENUM_KIND    // illegal rule "enum" on field with type of non-basic kind
	ERR_ENUM_NOCONST // "enum" rule on field with type that has no "known" constants declared

	ERR_LENGTH_NOLEN   // illegal rule "len" on field with lenght-less type
	ERR_LENGTH_NORUNE  // illegal rule "runecount" on field with non-stringy and non-[]byte type
	ERR_LENGTH_ARGTYPE // bad argument type in LENGTH rule
	ERR_LENGTH_NOARG   // no valid arguments in LENGTH rule
	ERR_LENGTH_BOUNDS  // invalid bounds arguments in LENGTH rule

	ERR_RANGE_TYPE    // illegal rule "rng" on non-numeric field
	ERR_RANGE_NOARG   // missing one or both values in arguments of RANGE rule
	ERR_RANGE_BOUNDS  // invalid bounds arguments in RANGE rule
	ERR_RANGE_ARGTYPE // bad argument type in RANGE rule

	ERR_ORDERED_TYPE    // illegal ORDERED rule on non-numeric/non-string field
	ERR_ORDERED_ARGTYPE // bad argument type in ORDERED rule

	ERR_PREPROC_INTYPE  // bad PREPROC rule function's input type, incompatible with node
	ERR_PREPROC_OUTTYPE // bad PREPROC rule function's output type, incompatible with node
	ERR_PREPROC_ARGTYPE // bad argument type in PREPROC rule

	ERR_FUNCTION_INTYPE   // bad FUNCTION rule function's input type, incompatible with node
	ERR_FUNCTION_ARGTYPE  // bad argument type in FUNCTION rule
	ERR_FUNCTION_ARGVALUE // bad argument value in FUNCTION rule

	// TODO rename
	ERR_ARG_BADCMP // argument's type incompatible with field's type (for comparison)

)

func (e ErrorCode) ident() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `

{{ define "` + ERR_CONFIG_FUNCID.ident() + `" -}}
{{ ERROR }} Invalid custom rule function identifier "{{R .CustomFuncIdent}}" in config file.
  > CONFIG: {{W .ConfigFile}}
  > HINT: A valid function identifier consists of the function's full package path followed
         by a dot (".") and the function's name, i.e. "{{W "<package_path>.<func_name>"}}".
         For example "{{W "github.com/me/mod/pkg/foo.FuncName"}}" is a valid rule function identifier.
{{ end }}

{{ define "` + ERR_CONFIG_FUNCSEARCH.ident() + `" -}}
{{ .OriginalError }}
  > CONFIG: {{W .ConfigFile}}
{{ end }}

{{ define "` + ERR_CONFIG_INVALID.ident() + `" -}}
{{ ERROR }} Failed to unmarshal custom rule config for "{{W .FuncIdent}}".
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > {{.ErrType}}: {{R (quote .OriginalError) }}
{{ end }}

{{ define "` + ERR_CONFIG_MISSING.ident() + `" -}}
{{ ERROR }} Missing custom rule config for function "{{W .FuncIdent}}".
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
{{ end }}

{{ define "` + ERR_CONFIG_NONAME.ident() + `" -}}
{{ ERROR }} Missing "{{R "name"}}" in custom rule's config.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > FUNC: {{W .FuncIdent}}
{{ end }}

{{ define "` + ERR_CONFIG_RESERVED.ident() + `" -}}
{{ ERROR }} Illegal use of {{wb "reserved"}} name "{{R .CfgRuleName}}" for custom rule.
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ else -}}
  > FILE: {{W .FuncPos}}
  {{ end -}}
  > FUNC: {{W .FuncIdent}}
{{ end }}

{{ define "` + ERR_CONFIG_FUNCTYPE.ident() + `" -}}
{{ ERROR }} Invalid function signature {{R .FuncType}} for custom "{{wb .CfgRuleName}}" rule.
  > FILE: {{W .FuncPos}}
  > FUNC: {{W .FuncIdent}}
  > HINT: A custom rule function MUST have {{wb "at least one"}} parameter value and
          it MUST have {{wb "exactly one"}} result value which MUST be of type {{wb "bool"}}.
{{ end }}

{{ define "` + ERR_CONFIG_PREFUNCTYPE.ident() + `" -}}
{{ ERROR }} Invalid function signature {{R .FuncType}} for custom pre-processor "{{wb .CfgRuleName}}" rule.
  > FILE: {{W .FuncPos}}
  > FUNC: {{W .FuncIdent}}
  > HINT: A custom pre-processor function MUST have {{wb "at least one"}} parameter value and
          it MUST have {{wb "exactly one"}} result value which MUST be of a type {{wb "identical"}} to
          the function's {{wb "first"}} parameter type.
{{ end }}

{{ define "` + ERR_CONFIG_ARGNUM.ident() + `" -}}
{{ ERROR }} Incompatible number of args in "{{wb .CfgRuleName}}" rule config for "{{W .FuncIdent}}".
  Config specifies {{rb .CfgArgNum }} {{.CfgArgWord}}. Function specifies {{wb .FuncParamNum }} ` +
	`{{.FuncParamWord}} (not counting the primary parameter).
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ end -}}
  > FILE: {{W .FuncPos}}
{{ end }}

{{ define "` + ERR_CONFIG_ARGTYPE.ident() + `" -}}
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

{{ define "` + ERR_CONFIG_ARGBOUNDS.ident() + `" -}}
{{ ERROR }} Incompatible arg bounds in rule config for "{{W .FuncIdent}}".
  Config specifies {{rb .CfgArgBounds }} bounds. Function specifies {{wb .FuncParamNum }} ` +
	`{{.FuncParamWord}} (not counting the primary parameter).
  {{if .HasConfigFile -}}
  > CONFIG: {{W .ConfigFile}}
  {{ end -}}
  > FILE: {{W .FuncPos}}
{{ end }}

{{ define "` + ERR_CONFIG_PREPROCJOIN.ident() + `" -}}
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

{{ define "` + ERR_CONFIG_PREPROCERROR.ident() + `" -}}
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

{{ define "` + ERR_RULE_UNDEFINED.ident() + `" -}}
{{ ERROR }} Undefined rule "{{R .RuleName}}" in "{{R .TagSTKey}}" struct tag.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + ERR_RULE_KEY.ident() + `" -}}
{{ ERROR }} Incompatible use of {{wb "key"}}-rule "{{R .TagKey}}" in "{{R .TagSTKey}}" struct tag` +
	` with a {{wb "non-map"}} type {{R .Type}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + ERR_RULE_ELEM.ident() + `" -}}
{{ ERROR }} Incompatible use of {{wb "elem"}}-rule "{{R .TagElem}}" in "{{R .TagSTKey}}" struct tag` +
	` with a {{wb "non-map/array/slice"}} type {{R .Type}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + ERR_NOTNIL_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "notnil"}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-nilable"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "notnil"}}" rule can ONLY be applied to fields with {{wb "nilable"}}` +
	` types like pointers, slices, maps, etc.
{{ end }}

{{ define "` + ERR_OPTIONAL_CONFLICT.ident() + `" -}}
{{ ERROR }} Incompatible use of "{{wb .RuleName}}" rule together with the "{{wb .Rule2Name}}"` +
	` rule in field {{wb .FieldName}} (type {{wb .FieldType}}).
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + ERR_ENUM_NONAME.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` of an {{R "unnamed"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with {{wb "named"}} types` +
	` of a {{wb "basic"}} kind.
{{ end }}

{{ define "` + ERR_ENUM_KIND.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" of a {{R "non-basic"}} kind.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with {{wb "named"}} types` +
	` of a {{wb "basic"}} kind.
{{ end }}

{{ define "` + ERR_ENUM_NOCONST.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" that has {{R "no constants"}} declared.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with types that have` +
	` {{wb "one or more constants"}} declared.
{{ end }}

{{ define "` + ERR_LENGTH_NOLEN.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" that has {{R "no length"}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with types that have` +
	` a {{wb "length"}}.
{{ end }}

{{ define "` + ERR_LENGTH_NORUNE.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` with {{R "non-string/non-[]byte"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with types` +
	` whose underlying type is either {{wb "string"}} or {{wb "[]byte"}}.
{{ end }}

{{ define "` + ERR_LENGTH_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb "uint"}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The arguments to the "{{wb .RuleName}}" rule MUST be {{wb "non-negative integers"}}.
{{ end }}

{{ define "` + ERR_LENGTH_NOARG.ident() + `" -}}
{{ ERROR }} Invalid use of the "{{wb .RuleName}}" rule with {{R "no values"}} in arguments.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule MUST contain {{wb "at least one"}} value in arguments.
{{ end }}

{{ define "` + ERR_LENGTH_BOUNDS.ident() + `" -}}
{{ ERROR }} Invalid arguments "{{R .RuleArgs}}" in the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: When {{wb "two"}} values are provided in the arguments of the "{{wb .RuleName}}" rule,` +
	`{{NT}}  then those values are used as the {{wb "lower and upper bounds"}}` +
	` of that rule.{{NT}} For that reason the first value MUST be {{wb "less than"}} the second.
{{ end }}

{{ define "` + ERR_RANGE_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-numeric"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with {{wb "numeric"}}` +
	` types like int, uint8, float64, etc.
{{ end }}

{{ define "` + ERR_RANGE_NOARG.ident() + `" -}}
{{ ERROR }} Invalid use of the "{{wb .RuleName}}" rule with {{R "missing value(s)"}} in arguments.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule MUST contain {{wb "both"}} values in arguments.
{{ end }}

{{ define "` + ERR_RANGE_BOUNDS.ident() + `" -}}
{{ ERROR }} Invalid arguments "{{R .RuleArgs}}" in the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The first value MUST be {{wb "less than"}} the second.
{{ end }}

{{ define "` + ERR_RANGE_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .FieldTypeFamily}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The values in the "{{wb .RuleName}}" rule's arguments MUST be {{wb "comparable"}} to` +
	` the field {{wb .FieldName}} (type "{{wb .FieldType}}").
{{ end }}

{{ define "` + ERR_ORDERED_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-string"}}/{{R "non-numeric"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with {{wb "string"}}` +
	` or {{wb "numeric"}} types like string, int32, float64, etc.
{{ end }}

{{ define "` + ERR_ORDERED_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .FieldTypeFamily}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "comparable"}} to` +
	` the field {{wb .FieldName}} (type "{{wb .FieldType}}").
{{ end }}

{{ define "` + ERR_PREPROC_INTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's 1st parameter type, i.e. {{wb .RuleFuncIn0Type}}.
{{ end }}

{{ define "` + ERR_PREPROC_OUTTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's return type, i.e. {{wb .RuleFuncOut0Type}}.
{{ end }}

{{ define "` + ERR_PREPROC_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .RuleFuncParamType}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "assignable"}} to` +
	`{{NT}}  the {{wb .RuleFuncIdent}} function's {{wb .FuncParamIdent}} parameter (type {{wb .RuleFuncParamType}}).
{{ end }}

{{ define "` + ERR_FUNCTION_INTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's 1st parameter type, i.e. {{wb .RuleFuncIn0Type}}.
{{ end }}

{{ define "` + ERR_FUNCTION_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .RuleFuncParamType}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "assignable"}} to` +
	`{{NT}}  the {{wb .RuleFuncIdent}} function's {{wb .FuncParamIdent}} parameter (type {{wb .RuleFuncParamType}}).
{{ end }}

{{ define "` + ERR_FUNCTION_ARGVALUE.ident() + `" -}}
{{ ERROR }} Cannot use value "{{R .RuleArgValue}}" as argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: For a set of valid argument values, see the "{{wb .RuleName}}" rule's spec as defined` +
	`{{NT}}  in the config file or in the {{wb .RuleFuncIdent}} function's documentation.
{{ end }}

{{ define "` + ERR_ARG_BADCMP.ident() + `" -}}
{{ ERROR }} Invalid rule "{{wb .Rule}}", argument {{wb .RuleArgValue}} of type "{{R .RuleArgType}}"` +
	`{{if .RuleArgIsField}} and type "{{R .RuleArgFieldType}}"{{end}}` +
	` is NOT comparable to field {{wb .FieldName}} of type "{{R .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
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
