package types

import (
	"testing"

	"github.com/frk/compare"
)

func TestVisibleFields(t *testing.T) {
	type structField struct {
		name string
	}

	tests := []struct {
		name string
		want []structField
	}{{
		name: "VF_SimpleStruct",
		want: []structField{
			{name: "A"},
			{name: "B"},
			{name: "C"},
		},
	}, {
		name: "VF_NonEmbeddedStructMember",
		want: []structField{
			{name: "A"},
		},
	}, {
		name: "VF_EmbeddedExportedStruct",
		want: []structField{
			{name: "VF_SFG"},
			{name: "F"},
			{name: "G"},
		},
	}, {
		name: "VF_EmbeddedUnexportedStruct",
		want: []structField{
			{name: "vf_sFG"},
			{name: "F"},
			{name: "G"},
		},
	}, {
		name: "VF_TwoEmbeddedStructsWithCancellingMembers",
		want: []structField{
			{name: "VF_SFG"},
			{name: "G"},
			{name: "VF_SF"},
		},
	}, {
		name: "VF_EmbeddedStructsWithSameFieldsAtDifferentDepths",
		want: []structField{
			{name: "VF_SFGH3"},
			{name: "VF_SFGH2"},
			{name: "VF_SFGH1"},
			{name: "VF_SFGH"},
			{name: "H"},
			{name: "VF_SG1"},
			{name: "VF_SG"},
			{name: "G"},
			{name: "VF_SFG2"},
			{name: "VF_SFG1"},
			{name: "VF_SFG"},
			{name: "VF_SF2"},
			{name: "VF_SF1"},
			{name: "VF_SF"},
			{name: "L"},
		},
	}, {
		name: "VF_EmbeddedPointerStruct",
		want: []structField{
			{name: "VF_SF"},
			{name: "F"},
		},
	}, {
		name: "VF_EmbeddedNotAPointer",
		want: []structField{
			{name: "VF_M"},
		},
	}, {
		name: "VF_RecursiveEmbedding",
		want: []structField{
			{name: "VF_Rec2"},
			{name: "F"},
			{name: "VF_Rec1"},
		},
	}, {
		name: "VF_RecursiveEmbedding2",
		want: []structField{
			{name: "F"},
			{name: "VF_Rec1"},
			{name: "VF_Rec2"},
		},
	}, {
		name: "VF_RecursiveEmbedding3",
		want: []structField{
			{name: "VF_RS2"},
			{name: "VF_RS1"},
			{name: "i"},
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ty := Analyze(test_type(tt.name), &test_ast)

			fields := ty.VisibleFields()
			got := make([]structField, len(fields))
			for i := range fields {
				got[i].name = fields[i].Name
			}

			if err := compare.Compare(got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
