package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/frk/tagutil"
)

type anError struct {
	Code errorCode
	// Name of the validator that caused the error.
	VtorName string
	// The file in which the validator that caused the error is declared.
	VtorFileName string
	// The line at which the validator that caused the error is declared.
	VtorFileLine int `cmp:"+"`
	// The name of the field which caused the error, may be empty.
	FieldName string
	// The type of the field which caused the error.
	FieldType string
	// The tag of the field which caused the error.
	FieldTag tagutil.Tag
	// The file in which the field that caused the error is defined.
	FieldFileName string
	// The line at which the field that caused the error is defined.
	FieldFileLine int `cmp:"+"`
	// The rule that caused the error.
	RuleName string
	// The rule arg that caused the error.
	RuleArg *RuleArg
	// The original error
	Err error `cmp:"-"`
}

func (e *anError) Error() string {
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (e *anError) FileAndLine() string {
	return e.FieldFileName + ":" + strconv.Itoa(e.FieldFileLine)
}

func (e *anError) VtorFileAndLine() string {
	return e.VtorFileName + ":" + strconv.Itoa(e.VtorFileLine)
}

type errorCode uint8

func (e errorCode) name() string { return fmt.Sprintf("error_template_%d", e) }

const (
	_ errorCode = iota
	errEmptyValidator
	errRuleNameUnavailable
	errRuleUnknown
	errRuleContextUnknown
	errRuleArgNum
	errRuleArgType
	errRuleArgTypeUint
	errRuleArgTypeNint
	errRuleArgTypeFloat
	errRuleArgTypeString
	errRuleArgTypeReference
	errRuleArgTypeReferenceKind
	errRuleArgValueRegexp
	errRuleArgValueUUIDVer
	errRuleArgValueIPVer
	errRuleArgValueMACVer
	errRuleArgValueCountryCode
	errRuleArgValueLen
	errRuleArgValueISOStd
	errRuleArgValueRFCStd
	errRuleArgValueConflict
	errRuleArgValueParseFloat
	errRuleArgValueParseUint
	errRuleArgValueBounds
	errTypeNil
	errTypeLength
	errTypeNumeric
	errTypeKind
	errFieldKeyUnknown
	errFieldKeyConflict
	errErrorHandlerFieldConflict
	errContextOptionFieldConflict
	errContextOptionFieldRequired
	errContextOptionFieldType
	errRuleFuncParamCount
	errRuleFuncResultType
	errRuleFuncFieldArgType
	errRuleFuncRuleArgCount
	errRuleFuncRuleArgType
	errRuleEnumTypeUnnamed
	errRuleEnumTypeNoConst
)

var error_template_string = `
{{ define "` + errEmptyValidator.name() + `" -}}
{{.VtorFileAndLine}}: {{R .VtorName}}
  > must have at least one field to validate.
{{ end }}

{{ define "` + errRuleNameUnavailable.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "rule name not available."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown rule."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleContextUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown rule context."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgNum.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad number of rule args."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeUint.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (uint)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeNint.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (nint)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeFloat.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (float)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeString.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (string)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeReference.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (reference)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgTypeReferenceKind.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg type (reference)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueRegexp.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (regexp)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueUUIDVer.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (uuid)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueIPVer.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (ip)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueMACVer.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (mac)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueCountryCode.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (country code)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueLen.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (len)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueISOStd.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (ISO)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueRFCStd.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value (RFC)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueConflict.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting rule arg values."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueParseFloat.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Parsing rule arg value as float failed."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueParseUint.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Parsing rule arg value as uint failed."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleArgValueBounds.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule arg value bounds."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errTypeNil.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Field type is not nilable."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errTypeLength.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Field type has no length."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errTypeNumeric.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Field type is not numeric."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errFieldKeyUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown field key."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errFieldKeyConflict.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting field key."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errErrorHandlerFieldConflict.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting error handler field."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errContextOptionFieldConflict.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting context field."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errContextOptionFieldRequired.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Missing context field."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errContextOptionFieldType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad context field type."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleFuncParamCount.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule func parameter/result count."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleFuncResultType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule func result type."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleFuncFieldArgType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Incompatible field type and rule func argument type."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleFuncRuleArgCount.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Incompatible rule func parameter count and rule argument count."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleFuncRuleArgType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Incompatible rule func parameter type and rule argument type."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleEnumTypeUnnamed.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "\"enum\" rule with unnamed type is not allowed."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}

{{ define "` + errRuleEnumTypeNoConst.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Type of field with \"enum\" rule has no constants declared."}}
	TODO {{R .FieldName}} {{R .FieldType}}
{{ end }}
` // `

var error_templates = template.Must(template.New("t").Funcs(template.FuncMap{
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
