package analysis

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/isvalid/internal/testutil"
	"github.com/frk/tagutil"
)

var tdata = testutil.ParseTestdata("../testdata")

func testRunAnalysis(name string, t *testing.T) (*ValidatorStruct, error) {
	named, pos := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	vs, err := Run(tdata.Fset, named, pos, &Info{})
	if err != nil {
		return nil, err
	}
	return vs, nil
}

func TestAnalysisRun(t *testing.T) {
	// ...
	tests := []struct {
		Name     string
		want     *ValidatorStruct
		err      error
		printerr bool
	}{{
		Name: "AnalysisTestOK1_Validator",
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK1_Validator",
			Fields: []*StructField{{
				Name: "UserInput",
				Key:  "UserInput",
				Type: Type{
					Kind: TypeKindPtr,
					Elem: &Type{
						Name:       "UserInput",
						Kind:       TypeKindStruct,
						PkgPath:    "path/to/test",
						PkgName:    "testdata",
						PkgLocal:   "testdata",
						IsExported: true,
						Fields: []*StructField{{
							Name: "CountryCode", Key: "CountryCode",
							Type: Type{Kind: TypeKindString}, IsExported: true,
						}, {
							Name: "SomeVersion", Key: "SomeVersion",
							Type: Type{Kind: TypeKindString}, IsExported: true,
						}, {
							Name: "SomeValue", Key: "SomeValue",
							Type: Type{Kind: TypeKindString}, IsExported: true,
						}, {
							Name: "f1", Key: "f1",
							Tag:  tagutil.Tag{"is": []string{"required"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "required"}},
						}, {
							Name: "f2", Key: "f2",
							Tag:  tagutil.Tag{"is": []string{"required:@create"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "required", Params: []*RuleParam{
								{Value: "create", Kind: ParamKindContext},
							}}},
						}, {
							Name: "f3", Key: "f3",
							Tag:  tagutil.Tag{"is": []string{"required:#key"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "required", Params: []*RuleParam{
								{Value: "key", Kind: ParamKindGroupKey},
							}}},
						}, {
							Name: "f4", Key: "f4",
							Tag:  tagutil.Tag{"is": []string{"required:@create:#key"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "required", Params: []*RuleParam{
								{Value: "create", Kind: ParamKindContext},
								{Value: "key", Kind: ParamKindGroupKey},
							}}},
						}, {
							Name: "f5", Key: "f5",
							Tag:  tagutil.Tag{"is": []string{"email"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "email"}},
						}, {
							Name: "f6", Key: "f6",
							Tag:  tagutil.Tag{"is": []string{"url"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "url"}},
						}, {
							Name: "f7", Key: "f7",
							Tag:  tagutil.Tag{"is": []string{"uri"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "uri"}},
						}, {
							Name: "f8", Key: "f8",
							Tag:  tagutil.Tag{"is": []string{"pan"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "pan"}},
						}, {
							Name: "f9", Key: "f9",
							Tag:  tagutil.Tag{"is": []string{"cvv"}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "cvv"}},
						}, {
							Name: "F10", Key: "F10",
							Tag:  tagutil.Tag{"is": []string{"ssn"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ssn"}},
						}, {
							Name: "F11", Key: "F11",
							Tag:  tagutil.Tag{"is": []string{"ein"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ein"}},
						}, {
							Name: "F12", Key: "F12",
							Tag:  tagutil.Tag{"is": []string{"numeric"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "numeric"}},
						}, {
							Name: "F13", Key: "F13",
							Tag:  tagutil.Tag{"is": []string{"hex"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "hex"}},
						}, {
							Name: "F14", Key: "F14",
							Tag:  tagutil.Tag{"is": []string{"hexcolor"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "hexcolor"}},
						}, {
							Name: "F15", Key: "F15",
							Tag:  tagutil.Tag{"is": []string{"alphanum"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "alphanum"}},
						}, {
							Name: "F16", Key: "F16",
							Tag:  tagutil.Tag{"is": []string{"cidr"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "cidr"}},
						}, {
							Name: "F17", Key: "F17",
							Tag:  tagutil.Tag{"is": []string{"phone"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "phone"}},
						}, {
							Name: "F18", Key: "F18",
							Tag:  tagutil.Tag{"is": []string{"phone:us"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "phone", Params: []*RuleParam{
								{Value: "us", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F19", Key: "F19",
							Tag:  tagutil.Tag{"is": []string{"phone:&CountryCode"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "phone", Params: []*RuleParam{
								{Value: "CountryCode", Kind: ParamKindReference},
							}}},
						}, {
							Name: "F20", Key: "F20",
							Tag:  tagutil.Tag{"is": []string{"zip"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "zip"}},
						}, {
							Name: "F21", Key: "F21",
							Tag:  tagutil.Tag{"is": []string{"zip:deu"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "zip", Params: []*RuleParam{
								{Value: "deu", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F22", Key: "F22",
							Tag:  tagutil.Tag{"is": []string{"zip:&CountryCode"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "zip", Params: []*RuleParam{
								{Value: "CountryCode", Kind: ParamKindReference},
							}}},
						}, {
							Name: "F23", Key: "F23",
							Tag:  tagutil.Tag{"is": []string{"uuid"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "uuid"}},
						}, {
							Name: "F24", Key: "F24",
							Tag:  tagutil.Tag{"is": []string{"uuid:3"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "uuid", Params: []*RuleParam{
								{Value: "3", Kind: ParamKindLiteral, Type: ParamTypeUint},
							}}},
						}, {
							Name: "F25", Key: "F25",
							Tag:  tagutil.Tag{"is": []string{"uuid:v4"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "uuid", Params: []*RuleParam{
								{Value: "v4", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F26", Key: "F26",
							Tag:  tagutil.Tag{"is": []string{"uuid:&SomeVersion"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "uuid", Params: []*RuleParam{
								{Value: "SomeVersion", Kind: ParamKindReference},
							}}},
						}, {
							Name: "F27", Key: "F27",
							Tag:  tagutil.Tag{"is": []string{"ip"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ip"}},
						}, {
							Name: "F28", Key: "F28",
							Tag:  tagutil.Tag{"is": []string{"ip:4"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ip", Params: []*RuleParam{
								{Value: "4", Kind: ParamKindLiteral, Type: ParamTypeUint},
							}}},
						}, {
							Name: "F29", Key: "F29",
							Tag:  tagutil.Tag{"is": []string{"ip:v6"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ip", Params: []*RuleParam{
								{Value: "v6", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F30", Key: "F30",
							Tag:  tagutil.Tag{"is": []string{"ip:&SomeVersion"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "ip", Params: []*RuleParam{
								{Value: "SomeVersion", Kind: ParamKindReference},
							}}},
						}, {
							Name: "F31", Key: "F31",
							Tag:  tagutil.Tag{"is": []string{"mac"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "mac"}},
						}, {
							Name: "F32", Key: "F32",
							Tag:  tagutil.Tag{"is": []string{"mac:6"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "mac", Params: []*RuleParam{
								{Value: "6", Kind: ParamKindLiteral, Type: ParamTypeUint},
							}}},
						}, {
							Name: "F33", Key: "F33",
							Tag:  tagutil.Tag{"is": []string{"mac:v8"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "mac", Params: []*RuleParam{
								{Value: "v8", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F34", Key: "F34",
							Tag:  tagutil.Tag{"is": []string{"mac:&SomeVersion"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "mac", Params: []*RuleParam{
								{Value: "SomeVersion", Kind: ParamKindReference},
							}}},
						}, {
							Name: "F35", Key: "F35",
							Tag:  tagutil.Tag{"is": []string{"iso:1234"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "iso", Params: []*RuleParam{
								{Value: "1234", Kind: ParamKindLiteral, Type: ParamTypeUint},
							}}},
						}, {
							Name: "F36", Key: "F36",
							Tag:  tagutil.Tag{"is": []string{"rfc:1234"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "rfc", Params: []*RuleParam{
								{Value: "1234", Kind: ParamKindLiteral, Type: ParamTypeUint},
							}}},
						}, {
							Name: "F37", Key: "F37",
							Tag:  tagutil.Tag{"is": []string{`re:"^[a-z]+\[[0-9]+\]$"`}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "re", Params: []*RuleParam{
								{Value: `^[a-z]+\[[0-9]+\]$`, Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F38", Key: "F38",
							Tag:  tagutil.Tag{"is": []string{`re:"\w+"`}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "re", Params: []*RuleParam{
								{Value: `\w+`, Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F39", Key: "F39",
							Tag:  tagutil.Tag{"is": []string{`contains:foo bar`}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "contains", Params: []*RuleParam{
								{Value: "foo bar", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "F40", Key: "F40",
							Tag:  tagutil.Tag{"is": []string{`contains:&SomeValue`}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							Rules: []*Rule{{Name: "contains", Params: []*RuleParam{
								{Value: "SomeValue", Kind: ParamKindReference},
							}}},
						}, {
							Name: "f41", Key: "f41",
							Tag:  tagutil.Tag{"is": []string{`prefix:foo bar`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "prefix", Params: []*RuleParam{
								{Value: "foo bar", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "f42", Key: "f42",
							Tag:  tagutil.Tag{"is": []string{`prefix:&SomeValue`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "prefix", Params: []*RuleParam{
								{Value: "SomeValue", Kind: ParamKindReference},
							}}},
						}, {
							Name: "f43", Key: "f43",
							Tag:  tagutil.Tag{"is": []string{`suffix:foo bar`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "suffix", Params: []*RuleParam{
								{Value: "foo bar", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "f44", Key: "f44",
							Tag:  tagutil.Tag{"is": []string{`suffix:&SomeValue`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "suffix", Params: []*RuleParam{
								{Value: "SomeValue", Kind: ParamKindReference},
							}}},
						}, {
							Name: "f45", Key: "f45",
							Tag:  tagutil.Tag{"is": []string{`eq:foo bar`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "eq", Params: []*RuleParam{
								{Value: "foo bar", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "f46", Key: "f46",
							Tag:  tagutil.Tag{"is": []string{`eq:-123`}},
							Type: Type{Kind: TypeKindInt}, IsExported: false,
							Rules: []*Rule{{Name: "eq", Params: []*RuleParam{
								{Value: "-123", Kind: ParamKindLiteral, Type: ParamTypeNint},
							}}},
						}, {
							Name: "f47", Key: "f47",
							Tag:  tagutil.Tag{"is": []string{`eq:123.987`}},
							Type: Type{Kind: TypeKindFloat64}, IsExported: false,
							Rules: []*Rule{{Name: "eq", Params: []*RuleParam{
								{Value: "123.987", Kind: ParamKindLiteral, Type: ParamTypeFloat},
							}}},
						}, {
							Name: "f48", Key: "f48",
							Tag:  tagutil.Tag{"is": []string{`eq:&SomeValue`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "eq", Params: []*RuleParam{
								{Value: "SomeValue", Kind: ParamKindReference},
							}}},
						}, {
							Name: "f49", Key: "f49",
							Tag:  tagutil.Tag{"is": []string{`ne:foo bar`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "ne", Params: []*RuleParam{
								{Value: "foo bar", Kind: ParamKindLiteral, Type: ParamTypeString},
							}}},
						}, {
							Name: "f50", Key: "f50",
							Tag:  tagutil.Tag{"is": []string{`ne:-123`}},
							Type: Type{Kind: TypeKindInt}, IsExported: false,
							Rules: []*Rule{{Name: "ne", Params: []*RuleParam{
								{Value: "-123", Kind: ParamKindLiteral, Type: ParamTypeNint},
							}}},
						}, {
							Name: "f51", Key: "f51",
							Tag:  tagutil.Tag{"is": []string{`ne:123.987`}},
							Type: Type{Kind: TypeKindFloat64}, IsExported: false,
							Rules: []*Rule{{Name: "ne", Params: []*RuleParam{
								{Value: "123.987", Kind: ParamKindLiteral, Type: ParamTypeFloat},
							}}},
						}, {
							Name: "f52", Key: "f52",
							Tag:  tagutil.Tag{"is": []string{`ne:&SomeValue`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "ne", Params: []*RuleParam{
								{Value: "SomeValue", Kind: ParamKindReference}}}},
						}, {
							Name: "f53", Key: "f53",
							Tag:  tagutil.Tag{"is": []string{`gt:24`, `lt:128`}},
							Type: Type{Kind: TypeKindUint8}, IsExported: false,
							Rules: []*Rule{
								{Name: "gt", Params: []*RuleParam{
									{Value: "24", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
								{Name: "lt", Params: []*RuleParam{
									{Value: "128", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f54", Key: "f54",
							Tag:  tagutil.Tag{"is": []string{`gt:-128`, `lt:-24`}},
							Type: Type{Kind: TypeKindInt16}, IsExported: false,
							Rules: []*Rule{
								{Name: "gt", Params: []*RuleParam{
									{Value: "-128", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
								{Name: "lt", Params: []*RuleParam{
									{Value: "-24", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
							},
						}, {
							Name: "f55", Key: "f55",
							Tag:  tagutil.Tag{"is": []string{`gt:0.24`, `lt:1.28`}},
							Type: Type{Kind: TypeKindFloat32}, IsExported: false,
							Rules: []*Rule{
								{Name: "gt", Params: []*RuleParam{
									{Value: "0.24", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
								{Name: "lt", Params: []*RuleParam{
									{Value: "1.28", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
							},
						}, {
							Name: "f56", Key: "f56",
							Tag:  tagutil.Tag{"is": []string{`gte:24`, `lte:128`}},
							Type: Type{Kind: TypeKindUint8}, IsExported: false,
							Rules: []*Rule{
								{Name: "gte", Params: []*RuleParam{
									{Value: "24", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
								{Name: "lte", Params: []*RuleParam{
									{Value: "128", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f57", Key: "f57",
							Tag:  tagutil.Tag{"is": []string{`gte:-128`, `lte:-24`}},
							Type: Type{Kind: TypeKindInt16}, IsExported: false,
							Rules: []*Rule{
								{Name: "gte", Params: []*RuleParam{
									{Value: "-128", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
								{Name: "lte", Params: []*RuleParam{
									{Value: "-24", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
							},
						}, {
							Name: "f58", Key: "f58",
							Tag:  tagutil.Tag{"is": []string{`gte:0.24`, `lte:1.28`}},
							Type: Type{Kind: TypeKindFloat32}, IsExported: false,
							Rules: []*Rule{
								{Name: "gte", Params: []*RuleParam{
									{Value: "0.24", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
								{Name: "lte", Params: []*RuleParam{
									{Value: "1.28", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
							},
						}, {
							Name: "f59", Key: "f59",
							Tag:  tagutil.Tag{"is": []string{`min:24`, `max:128`}},
							Type: Type{Kind: TypeKindUint8}, IsExported: false,
							Rules: []*Rule{
								{Name: "min", Params: []*RuleParam{
									{Value: "24", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
								{Name: "max", Params: []*RuleParam{
									{Value: "128", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f60", Key: "f60",
							Tag:  tagutil.Tag{"is": []string{`min:-128`, `max:-24`}},
							Type: Type{Kind: TypeKindInt16}, IsExported: false,
							Rules: []*Rule{
								{Name: "min", Params: []*RuleParam{
									{Value: "-128", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
								{Name: "max", Params: []*RuleParam{
									{Value: "-24", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
							},
						}, {
							Name: "f61", Key: "f61",
							Tag:  tagutil.Tag{"is": []string{`min:0.24`, `max:1.28`}},
							Type: Type{Kind: TypeKindFloat32}, IsExported: false,
							Rules: []*Rule{
								{Name: "min", Params: []*RuleParam{
									{Value: "0.24", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
								{Name: "max", Params: []*RuleParam{
									{Value: "1.28", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
							},
						}, {
							Name: "f62", Key: "f62",
							Tag:  tagutil.Tag{"is": []string{`rng:24:128`}},
							Type: Type{Kind: TypeKindUint8}, IsExported: false,
							Rules: []*Rule{{Name: "rng", Params: []*RuleParam{
								{Value: "24", Kind: ParamKindLiteral, Type: ParamTypeUint},
								{Value: "128", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f63", Key: "f63",
							Tag:  tagutil.Tag{"is": []string{`rng:-128:-24`}},
							Type: Type{Kind: TypeKindInt16}, IsExported: false,
							Rules: []*Rule{{Name: "rng", Params: []*RuleParam{
								{Value: "-128", Kind: ParamKindLiteral, Type: ParamTypeNint},
								{Value: "-24", Kind: ParamKindLiteral, Type: ParamTypeNint}}},
							},
						}, {
							Name: "f64", Key: "f64",
							Tag:  tagutil.Tag{"is": []string{`rng:0.24:1.28`}},
							Type: Type{Kind: TypeKindFloat32}, IsExported: false,
							Rules: []*Rule{{Name: "rng", Params: []*RuleParam{
								{Value: "0.24", Kind: ParamKindLiteral, Type: ParamTypeFloat},
								{Value: "1.28", Kind: ParamKindLiteral, Type: ParamTypeFloat}}},
							},
						}, {
							Name: "f65", Key: "f65",
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f66", Key: "f66",
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f67", Key: "f67",
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f68", Key: "f68",
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "4", Kind: ParamKindLiteral, Type: ParamTypeUint},
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f69", Key: "f69",
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "4", Kind: ParamKindLiteral, Type: ParamTypeUint},
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f70", Key: "f70",
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "4", Kind: ParamKindLiteral, Type: ParamTypeUint},
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f71", Key: "f71",
							Tag:  tagutil.Tag{"is": []string{`len::28`}},
							Type: Type{Kind: TypeKindString}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "", Kind: ParamKindLiteral},
								{Value: "28", Kind: ParamKindLiteral, Type: ParamTypeUint}}},
							},
						}, {
							Name: "f72", Key: "f72",
							Tag:  tagutil.Tag{"is": []string{`len:4:`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}}, IsExported: false,
							Rules: []*Rule{{Name: "len", Params: []*RuleParam{
								{Value: "4", Kind: ParamKindLiteral, Type: ParamTypeUint},
								{Value: "", Kind: ParamKindLiteral}}},
							},
						}},
					},
				},
				IsExported: true,
			}},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := testRunAnalysis(tt.Name, t)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("Error: %v", e)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}

			tt.printerr = true
			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}
