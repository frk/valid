package analysis

import (
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type anError struct {
	Code errorCode

	a *analysis `cmp:"+"`
	// The field that's associated with the error.
	f *StructField `cmp:"+"`
	// The rule that's associated with the error.
	r *Rule `cmp:"+"`
	// The rule arg that's associated with the error.
	ra *RuleArg
	// The rule function that's associated with the error.
	fn *types.Func `cmp:"+"`
	// The original error
	err error `cmp:"-"`
}

func (e *anError) Error() string {
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (e *anError) VtorName() string {
	return e.a.named.Obj().Name()
}

func (e *anError) VtorFileAndLine() string {
	obj := e.a.named.Obj()
	pos := e.a.fset.Position(obj.Pos())
	return pos.Filename + ":" + strconv.Itoa(pos.Line)
}

func (e *anError) VtorErrorHandler() string {
	return e.a.validator.ErrorHandler.Name
}

func (e *anError) VtorContextOption() string {
	return e.a.validator.ContextOption.Name
}

func (e *anError) FieldName() string {
	return e.f.Name
}

func (e *anError) FieldType() string {
	if e.f.Type.Kind == TypeKindInvalid {
		if fv, ok := e.a.fieldVarMap[e.f]; ok {
			return fv.v.Type().String()
		}
	}
	return e.f.Type.String()
}

func (e *anError) FieldTagRaw(name string) string {
	if fv, ok := e.a.fieldVarMap[e.f]; ok {
		if tag, ok := reflect.StructTag(fv.tag).Lookup(name); ok {
			return "`is:\"" + tag + "\"`"
		}
	}
	return "``"
}

func (e *anError) FieldNameAndType() string {
	return e.f.Name + " " + e.FieldType()
}

func (e *anError) FileAndLine() string {
	if e.a != nil {
		if fv, ok := e.a.fieldVarMap[e.f]; ok {
			pos := e.a.fset.Position(fv.v.Pos())
			return pos.Filename + ":" + strconv.Itoa(pos.Line)
		}
	}
	return "[unknown-source-location]"
}

func (e *anError) RuleName() string {
	return e.r.Name
}

func (e *anError) RuleContext() string {
	return "@" + e.r.Context
}

func (e *anError) RuleArgValue() string {
	switch e.ra.Type {
	case ArgTypeBool, ArgTypeInt, ArgTypeFloat:
		return e.ra.Value
	case ArgTypeString, ArgTypeUnknown:
		return `"` + e.ra.Value + `"`
	case ArgTypeField:
		return "&" + e.ra.Value
	}

	panic("shouldn't reach")
	return ""
}

func (e *anError) RuleArgFieldKey() string {
	if e.ra.Type == ArgTypeField {
		return e.ra.Value
	}
	return ""
}

func (e *anError) RuleArgType() string {
	if e.ra.Type == ArgTypeField {
		sel := e.a.info.SelectorMap[e.ra.Value]
		return sel.Last().Type.String()
	}
	return e.ra.Type.String()
}

func (e *anError) RuleArgPos() (out string) {
	var pos int
	for i, ra := range e.r.Args {
		if e.ra == ra {
			pos = i + 1
			break
		}
	}

	switch pos % 100 {
	case 11, 12, 13:
		return strconv.Itoa(pos) + "th"
	}

	switch pos % 10 {
	case 1:
		return strconv.Itoa(pos) + "st"
	case 2:
		return strconv.Itoa(pos) + "nd"
	case 3:
		return strconv.Itoa(pos) + "rd"
	}
	return strconv.Itoa(pos) + "th"
}

func (e *anError) RuleArgs() (out string) {
	for _, ra := range e.r.Args {
		if ra.Type == ArgTypeField {
			out += ":&" + ra.Value
		} else {
			out += ":" + ra.Value
		}
	}
	return out
}

func (e *anError) RuleArgNum() string {
	return strconv.Itoa(len(e.r.Args))
}

func (e *anError) RuleArgNumWord() string {
	if len(e.r.Args) == 1 {
		return "argument"
	}
	return "arguments"
}

func (e *anError) RuleArgCount() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
		if !ok {
			return "<unknown-argument-count>"
		}
	}

	count := rt.argCount()
	if count.min == count.max {
		return strconv.Itoa(count.min)
	}
	if count.min == 0 && count.max > 0 {
		return "at most " + strconv.Itoa(count.max)
	}
	if count.min > 0 && count.max == -1 {
		return "at least " + strconv.Itoa(count.min)
	}
	if count.min > 0 && count.max > 0 {
		return "between " + strconv.Itoa(count.min) +
			" and " + strconv.Itoa(count.max)
	}
	return out
}

func (e *anError) RuleArgCountWord() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
		if !ok {
			return "arguments"
		}
	}

	count := rt.argCount()
	if count.min == 1 && count.max == 1 {
		return "argument"
	}
	return "arguments"
}

func (e *anError) FuncNameQualified() (out string) {
	return e.fn.Pkg().Name() + "." + e.fn.Name()
}

func (e *anError) FuncType() (out string) {
	return e.fn.Type().String()
}

func (e *anError) FuncFieldType() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
	}
	if fn, ok := rt.(RuleTypeFunc); ok {
		return fn.ArgTypes[0].String()
	}
	return "<unknown-func>"
}

func (e *anError) FuncArgType() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
	}

	switch rx := rt.(type) {
	case RuleTypeBasic:
		return e.f.Type.String()
	case RuleTypeFunc:
		var pos int
		for i, ra := range e.r.Args {
			if e.ra == ra {
				pos = i + 1
				break
			}
		}

		if rx.IsVariadic && pos >= (len(rx.ArgTypes)-1) {
			return rx.ArgTypes[len(rx.ArgTypes)-1].Elem.String()
		}
		if rx.LOp > 0 {
			return rx.ArgTypes[1].String()
		}
		return rx.ArgTypes[pos].String()
	}

	return "<unknown-arg-type>"
}

func (e *anError) Err() (out string) {
	return e.err.Error()
}

type errorCode uint8

func (e errorCode) name() string { return fmt.Sprintf("error_template_%d", e) }

const (
	_ errorCode = iota
	errRuleNameReserved
	errRuleFuncSignature
	errRuleConfOpts // TODO
	errValidatorNoField
	errRuleUnknown
	errRuleArgCount
	errRuleArgFieldUnknown
	errRuleArgValueRegexp
	errRuleArgValueUUIDVer
	errRuleArgValueIPVer
	errRuleArgValueMACVer
	errRuleArgValueCountryCode
	errRuleArgValueLanguageTag
	errRuleArgValueISONum
	errRuleArgValueRFCNum
	errRuleArgValueConflict
	errRuleArgValueBounds
	errRuleFieldNonNilable
	errRuleFieldLengthless
	errRuleFieldRuneless
	errRuleFieldNonNumeric
	errRuleFuncFieldType
	errRuleFuncArgType
	errRuleBasicArgType
	errRuleBasicArgTypeUint
	errContextOptionFieldRequired
	errErrorHandlerFieldConflict
	errContextOptionFieldConflict
	errContextOptionFieldType
	errRuleEnumType
	errRuleEnumTypeUnnamed
	errRuleEnumTypeNoConst
	errRuleKey
	errRuleElem
)

var error_template_string = `
{{ define "` + errRuleNameReserved.name() + `" -}}
{{R "ERROR:"}} Cannot use reserved rule name "{{R .RuleName}}" for custom rule.
{{ end }}

{{ define "` + errRuleFuncSignature.name() + `" -}}
{{R "ERROR:"}} Cannot use function {{R .FuncNameQualified}} of type {{R .FuncType}} as custom rule function.
  > A custom rule function must have {{R "at least one"}} parameter value and it must have {{R "exactly one"}}` +
	` result value which must be of type {{R "bool"}}.
{{ end }}

{{ define "` + errRuleConfOpts.name() + `" -}}
{{R "ERROR:"}} bad things happend!?
{{ end }}

{{ define "` + errValidatorNoField.name() + `" -}}
{{R "ERROR:"}} {{.VtorFileAndLine}}: 
  Cannot use type {{R .VtorName}} as a validator struct type.
  > A validator struct type must have {{R "at least one"}} field that can be validated.
{{ end }}

{{ define "` + errRuleUnknown.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}:
  Cannot use "{{R .RuleName}}" as rule of field {{R .FieldName}} in {{R .VtorName}}.
  > The value "{{R .RuleName}}" does not match the name of any registered rule.
{{ end }}

{{ define "` + errRuleArgCount.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with {{R .RuleArgNum}} {{.RuleArgNumWord}}.
  > The rule "{{R .RuleName}}" must have {{R .RuleArgCount}} {{.RuleArgCountWord}}.
{{ end }}

{{ define "` + errRuleArgFieldUnknown.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}:
  Cannot use rule argument {{R .RuleArgValue}} in rule "{{R .RuleName}}" of field {{R .FieldNameAndType}}.
  > The value {{R .RuleArgFieldKey}} does not match the key of any field in {{R .VtorName}}.
{{ end }}

{{ define "` + errRuleArgValueRegexp.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a valid and compilable regular expression.
  > Error received from "regexp" package compiler: {{R .Err}}
{{ end }}

{{ define "` + errRuleArgValueUUIDVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a valid UUID version.
{{ end }}

{{ define "` + errRuleArgValueIPVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a valid IP version.
{{ end }}

{{ define "` + errRuleArgValueMACVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a valid MAC version.
{{ end }}

{{ define "` + errRuleArgValueCountryCode.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a valid Alpha-2 or Alpha-3 country code.
{{ end }}

{{ define "` + errRuleArgValueLanguageTag.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be one of the supported language tags.
{{ end }}

{{ define "` + errRuleArgValueISONum.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a number representing a supported ISO standard.
{{ end }}

{{ define "` + errRuleArgValueRFCNum.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleArgValue}} as argument for rule "{{R .RuleName}}".
  > The argument to rule "{{R .RuleName}}" must be a number representing a supported RFC standard.
{{ end }}

{{ define "` + errRuleArgValueConflict.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use argument {{R .RuleArgValue}} in rule "{{R .RuleName}}" more than once.
  > The argument {{R .RuleArgValue}} is in conflict with another argument representing the same value in the same rule.
{{ end }}

{{ define "` + errRuleArgValueBounds.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use arguments "{{R .RuleArgs}}" with rule "{{R .RuleName}}".
  > The arguments to "{{R .RuleName}}" must represent a valid combination of lower and upper boundary.
  > Valid combinations of boundaries are: {{R}}{{.RuleName}}:<lower>:{{off}} | {{R}}{{.RuleName}}::<upper>{{off}} | {{R}}{{.RuleName}}:<lower>:<upper>{{off}}
{{ end }}

{{ define "` + errRuleFieldNonNilable.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with field {{R .FieldNameAndType}}.
  > The rule "{{R .RuleName}}" must be used with a field that can be nil.
{{ end }}

{{ define "` + errRuleFieldLengthless.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with field {{R .FieldNameAndType}}.
  > The rule "{{R .RuleName}}" must be used with a field that has length.
{{ end }}

{{ define "` + errRuleFieldRuneless.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with field {{R .FieldNameAndType}}.
  > The rule "{{R .RuleName}}" must be used with a []byte field or a string kind field.
{{ end }}

{{ define "` + errRuleFieldNonNumeric.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with field {{R .FieldNameAndType}}.
  > The rule "{{R .RuleName}}" must be used with a field of one of the numeric types.
{{ end }}

{{ define "` + errContextOptionFieldRequired.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use argument {{R .RuleContext}} with rule "{{R .RuleName}}" in field {{R .FieldName}} in {{R .VtorName}}.
  > An argument starting with {{R "@"}} denotes a "{{R "context attribute"}}" of the "{{R .RuleName}}" rule and, ` +
	`it requires that the target {{R .VtorName}} struct type has, at the root, a corresponding field named ` +
	`{{R "context"}} (case insensitive).
{{ end }}

{{ define "` + errErrorHandlerFieldConflict.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot have more than one "{{R "error handler"}}" field in a validator struct.
  > The field {{R .FieldName}} in {{R .VtorName}} is in conflict with the {{R .VtorErrorHandler}} field.
{{ end }}

{{ define "` + errContextOptionFieldConflict.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot have more than one "{{R "context"}}" field in a validator struct.
  > The field {{R .FieldName}} in {{R .VtorName}} is in conflict with the {{R .VtorContextOption}} field.
{{ end }}

{{ define "` + errContextOptionFieldType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use field {{R .FieldName}} of type {{R .FieldType}} as a "{{R "context"}}" field.
  > A "{{R "context"}}" field's type must be of the {{R "string"}} kind.
{{ end }}

{{ define "` + errRuleFuncFieldType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with field {{R .FieldName}} of type {{R .FieldType}}.
  > Rule "{{R .RuleName}}" requires a field with a type convertible to {{R .FuncFieldType}}.
{{ end }}

{{ define "` + errRuleFuncArgType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleArgValue}} of type {{R .RuleArgType}} as the {{R .RuleArgPos}} argument to the "{{R .RuleName}}" rule.
  > The {{R .RuleArgPos}} argument of the "{{R .RuleName}}" rule must be of a type convertible to {{R .FuncArgType}}.
{{ end }}

{{ define "` + errRuleBasicArgType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleArgValue}} of type {{R .RuleArgType}} as the {{R .RuleArgPos}}` +
	` argument to the "{{R .RuleName}}" rule of field {{R .FieldName}}.
  > The {{R .RuleArgPos}} argument of the "{{R .RuleName}}" rule must be of a type` +
	` convertible to the {{R .FieldName}} field's type {{R .FieldType}}.
{{ end }}

{{ define "` + errRuleBasicArgTypeUint.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleArgValue}} of type {{R .RuleArgType}} as the {{R .RuleArgPos}}` +
	` argument to the "{{R .RuleName}}" rule.
  > The {{R .RuleArgPos}} argument of the "{{R .RuleName}}" rule must be of a type` +
	` convertible to {{R "uint"}}.
{{ end }}

{{ define "` + errRuleEnumTypeUnnamed.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use "{{R .RuleName}}" rule with field {{R .FieldName}} of unnamed type {{R .FieldType}}.
  > The "{{R .RuleName}}" rule must be used with fields of {{R "named, basic-kind"}} types.
{{ end }}

{{ define "` + errRuleEnumType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use "{{R .RuleName}}" rule with field {{R .FieldName}} of non-basic-kind type {{R .FieldType}}.
  > The "{{R .RuleName}}" rule must be used with fields of {{R "named, basic-kind"}} types.
{{ end }}

{{ define "` + errRuleEnumTypeNoConst.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use "{{R .RuleName}}" rule with field {{R .FieldName}} of type {{R .FieldType}}.
  > Type {{R .FieldType}} has no exported constants.
{{ end }}

{{ define "` + errRuleKey.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use key-rule in tag {{R (.FieldTagRaw "is")}} with field {{R .FieldName}} of type {{R .FieldType}}.
  > A key-rule must have a corresponding map key in the field's type.
{{ end }}

{{ define "` + errRuleElem.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use elem-rule in tag {{R (.FieldTagRaw "is")}} with field {{R .FieldName}} of type {{R .FieldType}}.
  > An elem-rule must have a corresponding array/slice/map element in the field's type.
{{ end }}
` // `

var error_templates = template.Must(template.New("t").Funcs(template.FuncMap{
	// yellow color (terminal)
	"y":  func(v ...string) string { return getcolor("\033[0;33m", v) },
	"yb": func(v ...string) string { return getcolor("\033[1;33m", v) },
	"yi": func(v ...string) string { return getcolor("\033[3;33m", v) },
	"yu": func(v ...string) string { return getcolor("\033[4;33m", v) },
	// white color (terminal)
	"w":  func(v ...string) string { return getcolor("\033[0;37m", v) },
	"wb": func(v ...string) string { return getcolor("\033[1;37m", v) },
	"wi": func(v ...string) string { return getcolor("\033[3;37m", v) },
	"wu": func(v ...string) string { return getcolor("\033[4;37m", v) },
	// cyan color (terminal)
	"c":  func(v ...string) string { return getcolor("\033[0;36m", v) },
	"cb": func(v ...string) string { return getcolor("\033[1;36m", v) },
	"ci": func(v ...string) string { return getcolor("\033[3;36m", v) },
	"cu": func(v ...string) string { return getcolor("\033[4;36m", v) },

	/////////////////////////////////////////////////////////////////////////
	// High Intensity
	/////////////////////////////////////////////////////////////////////////

	// red color HI (terminal)
	"R":  func(v ...string) string { return getcolor("\033[0;91m", v) },
	"Rb": func(v ...string) string { return getcolor("\033[1;91m", v) },
	"Ri": func(v ...string) string { return getcolor("\033[3;91m", v) },
	"Ru": func(v ...string) string { return getcolor("\033[4;91m", v) },
	// green color HI (terminal)
	"G":  func(v ...string) string { return getcolor("\033[0;92m", v) },
	"Gb": func(v ...string) string { return getcolor("\033[1;92m", v) },
	"Gi": func(v ...string) string { return getcolor("\033[3;92m", v) },
	"Gu": func(v ...string) string { return getcolor("\033[4;92m", v) },
	// yellow color HI (terminal)
	"Y":  func(v ...string) string { return getcolor("\033[0;93m", v) },
	"Yb": func(v ...string) string { return getcolor("\033[1;93m", v) },
	"Yi": func(v ...string) string { return getcolor("\033[3;93m", v) },
	"Yu": func(v ...string) string { return getcolor("\033[4;93m", v) },
	// blue color HI (terminal)
	"B":  func(v ...string) string { return getcolor("\033[0;94m", v) },
	"Bb": func(v ...string) string { return getcolor("\033[1;94m", v) },
	"Bi": func(v ...string) string { return getcolor("\033[3;94m", v) },
	"Bu": func(v ...string) string { return getcolor("\033[4;94m", v) },
	// cyan color HI (terminal)
	"C":  func(v ...string) string { return getcolor("\033[0;96m", v) },
	"Cb": func(v ...string) string { return getcolor("\033[1;96m", v) },
	"Ci": func(v ...string) string { return getcolor("\033[3;96m", v) },
	"Cu": func(v ...string) string { return getcolor("\033[4;96m", v) },
	// white color HI (terminal)
	"W":  func(v ...string) string { return getcolor("\033[0;97m", v) },
	"Wb": func(v ...string) string { return getcolor("\033[1;97m", v) },
	"Wi": func(v ...string) string { return getcolor("\033[3;97m", v) },
	"Wu": func(v ...string) string { return getcolor("\033[4;97m", v) },

	// no color (terminal)
	"off": func() string { return "\033[0m" },

	"raw": func(s string) string { return "`" + s + "`" },
	"Up":  strings.ToUpper,
}).Parse(error_template_string))

func getcolor(c string, v []string) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v\033[0m", c, stringsStringer(v))
	}
	return c
}

type stringsStringer []string

func (s stringsStringer) String() string {
	return strings.Join([]string(s), "")
}
