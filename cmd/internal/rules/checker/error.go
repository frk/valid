package checker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/errors"
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

type Error struct {
	C ErrorCode

	// The struct field to which the error related, or nil.
	sf *types.StructField
	// The struct field referenced by a rule to which the error is associated, or nil.
	raf *types.StructField
	// The tag to which the error is related, or nil.
	tag *rules.TagNode
	// The type to which the error is related, or nil.
	ty *types.Type
	// The rule to which the error is related, or nil.
	r *rules.Rule
	// Set if the error is related to a combination of two rules, else nil.
	r2 *rules.Rule
	// The rule argument to which the error is related, or nil.
	ra *rules.Arg
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

func (e *Error) HasOriginalError() bool {
	return e.err != nil
}

func (e *Error) OriginalError() string {
	return strings.TrimRight(e.err.Error(), "\n")
}

func (e *Error) ErrType() string {
	return strings.TrimLeft(fmt.Sprintf("%T", e.err), "*")
}

func (e *Error) FieldPos() string {
	return types.GetPosition(e.sf)
}

func (e *Error) Field() string {
	return fmt.Sprintf("%s %s `%s`", e.sf.Name, e.sf.Obj.Type.TypeString(nil), e.sf.Tag)
}

func (e *Error) FieldName() string {
	return e.sf.Name
}

func (e *Error) FieldType() string {
	return e.sf.Obj.Type.TypeString(nil)
}

func (e *Error) TagKey() string {
	return "[" + e.tag.Key.String() + "]"
}

func (e *Error) TagElem() string {
	return "[]" + e.tag.Elem.String()
}

func (e *Error) TagSTKey() string {
	return e.tag.STKey() + ":"
}

func (e *Error) Type() string {
	return e.ty.TypeString(nil)
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

func (e *Error) RuleArgNum() string {
	return strconv.Itoa(len(e.r.Args))
}

func (e *Error) RuleArgValue() string {
	return e.ra.Value
}

func (e *Error) RuleSpecArgMin() string {
	return strconv.Itoa(e.r.Spec.ArgMin)
}

func (e *Error) RuleSpecArgMax() string {
	return strconv.Itoa(e.r.Spec.ArgMax)
}

func (e *Error) RuleSpecKind() string {
	return e.r.Spec.Kind.String()
}

func (e *Error) RuleFuncIdent() string {
	fn := specs.GetFunc(e.r.Spec)
	return fn.Type.Pkg.Name + "." + fn.Name
}

func (e *Error) RuleFuncType() string {
	fn := specs.GetFunc(e.r.Spec)
	return fn.Type.TypeString(nil)
}

func (e *Error) RuleFuncNameWithType() string {
	fn := specs.GetFunc(e.r.Spec)
	t := fn.Type.TypeString(nil)
	return fn.Name + strings.TrimPrefix(t, "func")
}

func (e *Error) RuleFuncIn0Type() string {
	fn := specs.GetFunc(e.r.Spec)
	return fn.Type.In[0].Type.TypeString(nil)
}

func (e *Error) RuleFuncOut0Type() string {
	fn := specs.GetFunc(e.r.Spec)
	return fn.Type.Out[0].Type.TypeString(nil)
}

func (e *Error) RuleFuncParamType() string {
	fn := specs.GetFunc(e.r.Spec)
	if fn.Type.IsVariadic && len(fn.Type.In) < *e.fpi {
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

func (e *Error) RuleArgType() string {
	if e.ra.IsFieldRef() {
		return e.raf.Obj.Type.TypeString(nil)
	}
	return e.ra.Type.String()
}

func (e *Error) RuleArgIsField() bool {
	return e.ra.IsFieldRef()
}

func (e *Error) RuleArgFieldType() string {
	return e.raf.Obj.Type.TypeString(nil)
}

type ErrorCode uint

const (
	_ ErrorCode = iota

	E_RULE_UNDEFINED // illegal use of undefined rule
	E_RULE_KEY       // illegal rule key on non-map field
	E_RULE_ELEM      // illegal rule elem on non-map/non-slice/non-array field
	E_RULE_ARGMIN    // number of rule arguments is less than min
	E_RULE_ARGMAX    // number of rule arguments is more than max

	E_FIELD_UNKNOWN // unknown field referenced by rule argument

	E_NOTNIL_TYPE // illegal rule "notnil" on non-nilable field

	E_OPTIONAL_CONFLICT // an optional rule is in conflict with a required rule

	E_ENUM_NONAME  // illegal rule "enum" on field with unnamed type
	E_ENUM_KIND    // illegal rule "enum" on field with type of non-basic kind
	E_ENUM_NOCONST // "enum" rule on field with type that has no "known" constants declared

	E_LENGTH_NOLEN   // illegal rule "len" on field with lenght-less type
	E_LENGTH_NORUNE  // illegal rule "runecount" on field with non-stringy and non-[]byte type
	E_LENGTH_ARGTYPE // bad argument type in LENGTH rule
	E_LENGTH_NOARG   // no valid arguments in LENGTH rule
	E_LENGTH_BOUNDS  // invalid bounds arguments in LENGTH rule

	E_RANGE_TYPE    // illegal rule "rng" on non-numeric field
	E_RANGE_NOARG   // missing one or both values in arguments of RANGE rule
	E_RANGE_BOUNDS  // invalid bounds arguments in RANGE rule
	E_RANGE_ARGTYPE // bad argument type in RANGE rule

	E_ORDERED_TYPE    // illegal ORDERED rule on non-numeric/non-string field
	E_ORDERED_ARGTYPE // bad argument type in ORDERED rule

	E_PREPROC_INTYPE  // bad PREPROC rule function's input type, incompatible with node
	E_PREPROC_OUTTYPE // bad PREPROC rule function's output type, incompatible with node
	E_PREPROC_ARGTYPE // bad argument type in PREPROC rule
	E_PREPROC_INVALID // invalid PREPROC rule

	E_FUNCTION_INTYPE   // bad FUNCTION rule function's input type, incompatible with node
	E_FUNCTION_ARGTYPE  // bad argument type in FUNCTION rule
	E_FUNCTION_ARGVALUE // bad argument value in FUNCTION rule

	E_METHOD_TYPE // illegal METHOD rule on type that does not have the specified method

	// TODO rename
	E_ARG_BADCMP // argument's type incompatible with field's type (for comparison)

)

func (e ErrorCode) ident() string {
	return fmt.Sprintf("%T_%d", e, uint(e))
}

func init() {
	errors.ParseTemplate(error_template)
}

var error_template = `

{{ define "` + E_RULE_UNDEFINED.ident() + `" -}}
{{ ERROR }} Undefined rule "{{R .RuleName}}" in "{{R .TagSTKey}}" struct tag.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_RULE_KEY.ident() + `" -}}
{{ ERROR }} Incompatible use of {{wb "key"}}-rule "{{R .TagKey}}" in "{{R .TagSTKey}}" struct tag` +
	` with a {{wb "non-map"}} type {{R .Type}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_RULE_ELEM.ident() + `" -}}
{{ ERROR }} Incompatible use of {{wb "elem"}}-rule "{{R .TagElem}}" in "{{R .TagSTKey}}" struct tag` +
	` with a {{wb "non-map/array/slice"}} type {{R .Type}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_RULE_ARGMIN.ident() + `" -}}
{{ ERROR }} Invalid number of arguments in "{{wb .RuleName}}". Expected at least {{wb .RuleSpecArgMin}}` +
	` argument(s), but instead got {{wb .RuleArgNum}} argument(s).
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_RULE_ARGMAX.ident() + `" -}}
{{ ERROR }} Invalid number of arguments in "{{wb .RuleName}}". Expected at most {{wb .RuleSpecArgMax}}` +
	` argument(s), but instead got {{wb .RuleArgNum}} argument(s).
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_FIELD_UNKNOWN.ident() + `" -}}
{{ ERROR }} Unknown field "{{wb .RuleArgValue}}" referenced in rule argument.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_NOTNIL_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "notnil"}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-nilable"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "notnil"}}" rule can ONLY be applied to fields with {{wb "nilable"}}` +
	` types like pointers, slices, maps, etc.
{{ end }}

{{ define "` + E_OPTIONAL_CONFLICT.ident() + `" -}}
{{ ERROR }} Conflicting use of "{{wb .RuleName}}" rule together with the "{{wb .Rule2Name}}"` +
	` rule in field {{wb .FieldName}} (type {{wb .FieldType}}).
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_ENUM_NONAME.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` of an {{R "unnamed"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with {{wb "named"}} types` +
	` of a {{wb "basic"}} kind.
{{ end }}

{{ define "` + E_ENUM_KIND.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" of a {{R "non-basic"}} kind.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with {{wb "named"}} types` +
	` of a {{wb "basic"}} kind.
{{ end }}

{{ define "` + E_ENUM_NOCONST.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb "enum"}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" that has {{R "no constants"}} declared.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb "enum"}}" rule can ONLY be applied to fields with types that have` +
	` {{wb "one or more constants"}} declared.
{{ end }}

{{ define "` + E_LENGTH_NOLEN.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` with type "{{wb .FieldType}}" that has {{R "no length"}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with types that have` +
	` a {{wb "length"}}.
{{ end }}

{{ define "` + E_LENGTH_NORUNE.ident() + `" -}}
{{ ERROR }} Illegal use of "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` with {{R "non-string/non-[]byte"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with types` +
	` whose underlying type is either {{wb "string"}} or {{wb "[]byte"}}.
{{ end }}

{{ define "` + E_LENGTH_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb "uint"}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The arguments to the "{{wb .RuleName}}" rule MUST be {{wb "non-negative integers"}}.
{{ end }}

{{ define "` + E_LENGTH_NOARG.ident() + `" -}}
{{ ERROR }} Invalid use of the "{{wb .RuleName}}" rule with {{R "no values"}} in arguments.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule MUST contain {{wb "at least one"}} value in arguments.
{{ end }}

{{ define "` + E_LENGTH_BOUNDS.ident() + `" -}}
{{ ERROR }} Invalid arguments "{{R .RuleArgs}}" in the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: When {{wb "two"}} values are provided in the arguments of the "{{wb .RuleName}}" rule,` +
	`{{NT}}  then those values are used as the {{wb "lower and upper bounds"}}` +
	` of that rule.{{NT}} For that reason the first value MUST be {{wb "less than"}} the second.
{{ end }}

{{ define "` + E_RANGE_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-numeric"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with {{wb "numeric"}}` +
	` types like int, uint8, float64, etc.
{{ end }}

{{ define "` + E_RANGE_NOARG.ident() + `" -}}
{{ ERROR }} Invalid use of the "{{wb .RuleName}}" rule with {{R "missing value(s)"}} in arguments.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule MUST contain {{wb "both"}} values in arguments.
{{ end }}

{{ define "` + E_RANGE_BOUNDS.ident() + `" -}}
{{ ERROR }} Invalid arguments "{{R .RuleArgs}}" in the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The first value MUST be {{wb "less than"}} the second.
{{ end }}

{{ define "` + E_RANGE_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .FieldTypeFamily}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The values in the "{{wb .RuleName}}" rule's arguments MUST be {{wb "comparable"}} to` +
	` the field {{wb .FieldName}} (type "{{wb .FieldType}}").
{{ end }}

{{ define "` + E_ORDERED_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` of a {{R "non-string"}}/{{R "non-numeric"}} type "{{wb .FieldType}}".
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule can ONLY be applied to fields with {{wb "string"}}` +
	` or {{wb "numeric"}} types like string, int32, float64, etc.
{{ end }}

{{ define "` + E_ORDERED_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .FieldTypeFamily}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "comparable"}} to` +
	` the field {{wb .FieldName}} (type "{{wb .FieldType}}").
{{ end }}

{{ define "` + E_PREPROC_INTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's 1st parameter type, i.e. {{wb .RuleFuncIn0Type}}.
{{ end }}

{{ define "` + E_PREPROC_OUTTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's return type, i.e. {{wb .RuleFuncOut0Type}}.
{{ end }}

{{ define "` + E_PREPROC_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .RuleFuncParamType}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "assignable"}} to` +
	`{{NT}}  the {{wb .RuleFuncIdent}} function's {{wb .FuncParamIdent}} parameter (type {{wb .RuleFuncParamType}}).
{{ end }}

{{ define "` + E_PREPROC_INVALID.ident() + `" -}}
{{ ERROR }} Cannot use "{{wb .RuleName}}" (kind {{R .RuleSpecKind}}) as a preprocessor.` +
	` Only preprocessor rules can be used in {{wb "pre:"}}"..." tags.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_FUNCTION_INTYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule with function {{wb .RuleFuncIdent}}` +
	` (type {{wb .RuleFuncType}}) in field {{wb .FieldName}} (type "{{wb .FieldType}}").
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The "{{wb .RuleName}}" rule's function {{wb .RuleFuncIdent}} can ONLY be applied to fields` +
	`{{NT}}  with a type that matches the function's 1st parameter type, i.e. {{wb .RuleFuncIn0Type}}.
{{ end }}

{{ define "` + E_FUNCTION_ARGTYPE.ident() + `" -}}
{{ ERROR }} Cannot use "{{R .RuleArgValue}}" (type {{R .RuleArgType}}) as {{wb .RuleFuncParamType}}` +
	` value in argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: The value in the "{{wb .RuleName}}" rule's argument MUST be {{wb "assignable"}} to` +
	`{{NT}}  the {{wb .RuleFuncIdent}} function's {{wb .FuncParamIdent}} parameter (type {{wb .RuleFuncParamType}}).
{{ end }}

{{ define "` + E_FUNCTION_ARGVALUE.ident() + `" -}}
{{ ERROR }} Cannot use value "{{R .RuleArgValue}}" as the {{wb .FuncParamIdent}} argument to the "{{wb .RuleName}}" rule.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
  > HINT: For a set of valid argument values, see the "{{wb .RuleName}}" rule's spec as defined` +
	`{{NT}}  in the config file or in the {{wb .RuleFuncIdent}} function's documentation.
{{- if .HasOriginalError}}
  > {{.ErrType}}: {{R (quote .OriginalError) }}
{{- end}}
{{ end }}

{{ define "` + E_METHOD_TYPE.ident() + `" -}}
{{ ERROR }} Illegal use of the "{{wb .RuleName}}" rule in field {{wb .FieldName}}` +
	` (type {{wb .FieldType}}). Type {{wb .FieldType}} does not implement` +
	` the method {{wb .RuleFuncNameWithType}}.
  > FILE: {{W .FieldPos}}
  > FIELD: {{W .Field}}
{{ end }}

{{ define "` + E_ARG_BADCMP.ident() + `" -}}
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
