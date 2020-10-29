package generator

import (
	"strconv"
	"strings"

	"github.com/frk/isvalid/internal/analysis"

	GO "github.com/frk/ast/golang"
)

type errTextOption uint

const (
	errTextNoArg errTextOption = iota
	errTextWithOneArg
	errTextWithOredQuotedArgs
	errTextWithAndedArgs
	errTextForLenRule
)

type errText struct {
	text   string
	option errTextOption
}

func (t errText) errText(f *analysis.StructField, r *analysis.Rule) GO.ExprNode {
	switch t.option {
	case errTextWithOneArg:
		text := f.Key + " " + t.text + ": "
		if a := r.Args[0]; a.Type == analysis.ArgTypeString {
			text += strconv.Quote(a.Value)
		} else {
			text += a.Value
		}
		return GO.ValueLit(strconv.Quote(text))
	case errTextWithOredQuotedArgs:
		text := f.Key + " " + t.text + ": " + joinArgValues(f, r.Args, true, " or ")
		return GO.ValueLit(strconv.Quote(text))
	case errTextWithAndedArgs:
		text := f.Key + " " + t.text + ": " + joinArgValues(f, r.Args, false, " and ")
		return GO.ValueLit(strconv.Quote(text))
	case errTextForLenRule:
		text := f.Key + " " + t.text
		if len(r.Args) == 1 {
			text += ": " + r.Args[0].Value
		} else { // len(r.Args) == 2 is assumed
			a1, a2 := r.Args[0], r.Args[1]
			if len(a1.Value) > 0 && len(a2.Value) == 0 {
				text += " at least: " + a1.Value
			} else if len(a1.Value) == 0 && len(a2.Value) > 0 {
				text += " at most: " + a2.Value
			} else {
				text += " between: " + a1.Value + " and " + a2.Value + " (inclusive)"
			}
		}
		return GO.StringLit(text)
	}
	return GO.StringLit(f.Key + " " + t.text)
}

var errTextMap = map[string]errText{
	"required": {text: "is required"},
	"notnil":   {text: "cannot be nil"},
	"email":    {text: "must be a valid email"},
	"url":      {text: "must be a valid URL"},
	"uri":      {text: "must be a valid URI"},
	"pan":      {text: "must be a valid PAN"},
	"cvv":      {text: "must be a valid CVV"},
	"ssn":      {text: "must be a valid SSN"},
	"ein":      {text: "must be a valid EIN"},
	"numeric":  {text: "must contain only digits [0-9]"},
	"hex":      {text: "must be a valid hexadecimal string"},
	"hexcolor": {text: "must be a valid hex color code"},
	"alphanum": {text: "must be an alphanumeric string"},
	"cidr":     {text: "must be a valid CIDR"},
	"phone":    {text: "must be a valid phone number"},
	"zip":      {text: "must be a valid zip code"},
	"uuid":     {text: "must be a valid UUID"},
	"ip":       {text: "must be a valid IP"},
	"mac":      {text: "must be a valid MAC"},
	"iso":      {text: "must be a valid ISO", option: errTextWithOneArg},
	"rfc":      {text: "must be a valid RFC", option: errTextWithOneArg},
	"re":       {text: "must match the regular expression", option: errTextWithOneArg},
	"prefix":   {text: "must be prefixed with", option: errTextWithOredQuotedArgs},
	"suffix":   {text: "must be suffixed with", option: errTextWithOredQuotedArgs},
	"contains": {text: "must contain substring", option: errTextWithOredQuotedArgs},
	"eq":       {text: "must be equal to", option: errTextWithOredQuotedArgs},
	"ne":       {text: "must not be equal to", option: errTextWithOredQuotedArgs},
	"gt":       {text: "must be greater than", option: errTextWithOneArg},
	"lt":       {text: "must be less than", option: errTextWithOneArg},
	"gte":      {text: "must be greater than or equal to", option: errTextWithOneArg},
	"lte":      {text: "must be less than or equal to", option: errTextWithOneArg},
	"min":      {text: "must be greater than or equal to", option: errTextWithOneArg},
	"max":      {text: "must be less than or equal to", option: errTextWithOneArg},
	"rng":      {text: "must be between", option: errTextWithAndedArgs},
	"len":      {text: "must be of length", option: errTextForLenRule},
}

func joinArgValues(f *analysis.StructField, args []*analysis.RuleArg, quote bool, sep string) string {
	typ := f.Type
	for typ.Kind == analysis.TypeKindPtr {
		typ = *typ.Elem
	}

	vals := make([]string, len(args))
	for i, a := range args {
		val := a.Value
		if a.Type == analysis.ArgTypeString && typ.Kind.IsNumeric() && len(val) == 0 {
			val = "0"
		}

		if quote {
			vals[i] = strconv.Quote(val)
		} else {
			vals[i] = val
		}
	}
	return strings.Join(vals, sep)
}
