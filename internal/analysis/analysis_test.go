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
							Name: "F1", Key: "F1",
							Tag: tagutil.Tag{"is": []string{"len:8:512"}},
							Type: Type{
								Kind: TypeKindString,
							},
							IsExported: true,
							Rules: []*Rule{{
								Name: "len",
								Params: []*RuleParam{
									{Value: "8", Type: ParamTypeUint},
									{Value: "512", Type: ParamTypeUint},
								},
							}},
						}, {
							Name: "F2", Key: "F2",
							Tag: tagutil.Tag{"is": []string{"min:8", "max:512"}},
							Type: Type{
								Kind: TypeKindInt,
							},
							IsExported: true,
							Rules: []*Rule{
								{Name: "min", Params: []*RuleParam{
									{Value: "8", Type: ParamTypeUint},
								}},
								{Name: "max", Params: []*RuleParam{
									{Value: "512", Type: ParamTypeUint},
								}},
							},
						}, {
							Name: "F3", Key: "F3",
							Tag: tagutil.Tag{"is": []string{`re:"\w+"`}},
							Type: Type{
								Kind: TypeKindString,
							},
							IsExported: true,
							Rules: []*Rule{
								{Name: "re", Params: []*RuleParam{
									{Value: `\w+`, Type: ParamTypeString},
								}},
							},
						}, {
							Name: "F4", Key: "F4",
							Tag: tagutil.Tag{"is": []string{`ne:"foo"bar"`, `ne:"baz"`}},
							Type: Type{
								Kind: TypeKindString,
							},
							IsExported: true,
							Rules: []*Rule{
								{Name: "ne", Params: []*RuleParam{
									{Value: `foo"bar`, Type: ParamTypeString},
								}},
								{Name: "ne", Params: []*RuleParam{
									{Value: `baz`, Type: ParamTypeString},
								}},
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
				t.Errorf("%v - %#v %v", e, err, err)
				//t.Errorf("Error: %v", e)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}

			//tt.printerr = true
			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}
