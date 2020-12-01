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
	// The rule option that's associated with the error.
	opt *RuleOption
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

func (e *anError) RuleOptionValue() string {
	switch e.opt.Type {
	case OptionTypeBool, OptionTypeInt, OptionTypeFloat:
		return e.opt.Value
	case OptionTypeString, OptionTypeUnknown:
		return `"` + e.opt.Value + `"`
	case OptionTypeField:
		return "&" + e.opt.Value
	}

	panic("shouldn't reach")
	return ""
}

func (e *anError) RuleOptionFieldKey() string {
	if e.opt.Type == OptionTypeField {
		return e.opt.Value
	}
	return ""
}

func (e *anError) RuleOptionType() string {
	if e.opt.Type == OptionTypeField {
		sel := e.a.info.SelectorMap[e.opt.Value]
		return sel.Last().Type.String()
	}
	return e.opt.Type.String()
}

func (e *anError) RuleOptionPos() (out string) {
	var pos int
	for i, opt := range e.r.Options {
		if e.opt == opt {
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

func (e *anError) RuleOptions() (out string) {
	for _, opt := range e.r.Options {
		if opt.Type == OptionTypeField {
			out += ":&" + opt.Value
		} else {
			out += ":" + opt.Value
		}
	}
	return out
}

func (e *anError) RuleOptionNum() string {
	return strconv.Itoa(len(e.r.Options))
}

func (e *anError) RuleOptionNumWord() string {
	if len(e.r.Options) == 1 {
		return "option"
	}
	return "options"
}

func (e *anError) RuleOptionCount() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
		if !ok {
			return "<unknown-option-count>"
		}
	}

	count := rt.optCount()
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

func (e *anError) RuleOptionCountWord() (out string) {
	rt, ok := e.a.conf.customTypeMap[e.r.Name]
	if !ok {
		rt, ok = defaultRuleTypeMap[e.r.Name]
		if !ok {
			return "options"
		}
	}

	count := rt.optCount()
	if count.min == 1 && count.max == 1 {
		return "option"
	}
	return "options"
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
		return fn.FieldArgType.String()
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
		argtypes := append([]Type{rx.FieldArgType}, rx.OptionArgTypes...)

		var pos int
		for i, opt := range e.r.Options {
			if e.opt == opt {
				pos = i + 1
				break
			}
		}

		if rx.IsVariadic && pos >= (len(argtypes)-1) {
			return argtypes[len(argtypes)-1].Elem.String()
		}
		if rx.LOp > 0 {
			return argtypes[1].String()
		}
		return argtypes[pos].String()
	}

	return "<unknown-option-type>"
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
	errRuleOptionCount
	errRuleOptionFieldUnknown
	errRuleOptionValueRegexp
	errRuleOptionValueUUIDVer
	errRuleOptionValueIPVer
	errRuleOptionValueMACVer
	errRuleOptionValueCountryCode
	errRuleOptionValueLanguageTag
	errRuleOptionValueISONum
	errRuleOptionValueRFCNum
	errRuleOptionValueConflict
	errRuleOptionValueBounds
	errRuleFieldNonNilable
	errRuleFieldLengthless
	errRuleFieldRuneless
	errRuleFieldNonNumeric
	errRuleFuncFieldType
	errRuleFuncOptionType
	errRuleBasicOptionType
	errRuleBasicOptionTypeUint
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

{{ define "` + errRuleOptionCount.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use rule "{{R .RuleName}}" with {{R .RuleOptionNum}} {{.RuleOptionNumWord}}.
  > The rule "{{R .RuleName}}" must have {{R .RuleOptionCount}} {{.RuleOptionCountWord}}.
{{ end }}

{{ define "` + errRuleOptionFieldUnknown.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}:
  Cannot use rule option {{R .RuleOptionValue}} in rule "{{R .RuleName}}" of field {{R .FieldNameAndType}}.
  > The value {{R .RuleOptionFieldKey}} does not match the key of any field in {{R .VtorName}}.
{{ end }}

{{ define "` + errRuleOptionValueRegexp.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a valid and compilable regular expression.
  > Error received from "regexp" package compiler: {{R .Err}}
{{ end }}

{{ define "` + errRuleOptionValueUUIDVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a valid UUID version.
{{ end }}

{{ define "` + errRuleOptionValueIPVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a valid IP version.
{{ end }}

{{ define "` + errRuleOptionValueMACVer.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a valid MAC version.
{{ end }}

{{ define "` + errRuleOptionValueCountryCode.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a valid Alpha-2 or Alpha-3 country code.
{{ end }}

{{ define "` + errRuleOptionValueLanguageTag.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be one of the supported language tags.
{{ end }}

{{ define "` + errRuleOptionValueISONum.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a number representing a supported ISO standard.
{{ end }}

{{ define "` + errRuleOptionValueRFCNum.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use value {{R .RuleOptionValue}} as option for rule "{{R .RuleName}}".
  > The option to rule "{{R .RuleName}}" must be a number representing a supported RFC standard.
{{ end }}

{{ define "` + errRuleOptionValueConflict.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use option {{R .RuleOptionValue}} in rule "{{R .RuleName}}" more than once.
  > The option {{R .RuleOptionValue}} is in conflict with another option representing the same value in the same rule.
{{ end }}

{{ define "` + errRuleOptionValueBounds.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use options "{{R .RuleOptions}}" with rule "{{R .RuleName}}".
  > The options to "{{R .RuleName}}" must represent a valid combination of lower and upper boundary.
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
  Cannot use option {{R .RuleContext}} with rule "{{R .RuleName}}" in field {{R .FieldName}} in {{R .VtorName}}.
  > An option starting with {{R "@"}} denotes a "{{R "context attribute"}}" of the "{{R .RuleName}}" rule and, ` +
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

{{ define "` + errRuleFuncOptionType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleOptionValue}} of type {{R .RuleOptionType}} as the {{R .RuleOptionPos}} option to the "{{R .RuleName}}" rule.
  > The {{R .RuleOptionPos}} option of the "{{R .RuleName}}" rule must be of a type convertible to {{R .FuncArgType}}.
{{ end }}

{{ define "` + errRuleBasicOptionType.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleOptionValue}} of type {{R .RuleOptionType}} as the {{R .RuleOptionPos}}` +
	` option to the "{{R .RuleName}}" rule of field {{R .FieldName}}.
  > The {{R .RuleOptionPos}} option of the "{{R .RuleName}}" rule must be of a type` +
	` convertible to the {{R .FieldName}} field's type {{R .FieldType}}.
{{ end }}

{{ define "` + errRuleBasicOptionTypeUint.name() + `" -}}
{{R "ERROR:"}} {{.FileAndLine}}: 
  Cannot use {{R .RuleOptionValue}} of type {{R .RuleOptionType}} as the {{R .RuleOptionPos}}` +
	` option to the "{{R .RuleName}}" rule.
  > The {{R .RuleOptionPos}} option of the "{{R .RuleName}}" rule must be of a type` +
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
