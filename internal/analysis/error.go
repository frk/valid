package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

type anError struct {
	Code        errorCode
	StructField *StructField
	Rule        *Rule
	RuleParam   *RuleParam
	FileName    string
	FileLine    int
	Err         error
}

func (e *anError) Error() string {
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (e *anError) FileAndLine() string {
	return e.FileName + ":" + strconv.Itoa(e.FileLine)
}

func (e *anError) FieldName() string {
	return e.StructField.Name
}

type errorCode uint8

func (e errorCode) name() string { return fmt.Sprintf("error_template_%d", e) }

const (
	_ errorCode = iota
	errRuleUnknown
	errRuleContextUnknown
	errRuleParamKind
	errRuleParamKindConflict
	errRuleParamTypeUint
	errRuleParamTypeNint
	errRuleParamTypeFloat
	errRuleParamTypeString
	errRuleParamValueRegexp
	errRuleParamValueUUID
	errRuleParamValueIP
	errRuleParamValueMAC
	errRuleParamValueCountryCode
	errRuleParamValueLen
	errTypeLength
	errTypeNumeric
	errTypeString
	errFieldKeyUnknown
	errFieldKeyConflict
)

var error_template_string = `
{{ define "` + errRuleUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown rule."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleContextUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown rule context."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamKind.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramater kind."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamKindConflict.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting rule parameter kind."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamTypeUint.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule parameter type (uint)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamTypeNint.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule parameter type (nint)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamTypeFloat.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule parameter type (float)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamTypeString.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule parameter type (string)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamValueRegexp.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramter value (regexp)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamValueUUID.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramter value (uuid)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamValueIP.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramter value (ip)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamValueMAC.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramter value (mac)."}}
	TODO {{R .FieldName}}
{{ end }}

{{ define "` + errRuleParamValueCountryCode.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad rule paramter value (country code)."}}
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

{{ define "` + errTypeString.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Field type is not string"}}
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
