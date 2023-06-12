package testdata

type VF_SimpleStruct struct {
	A int
	B string
	C bool
}

type VF_NonEmbeddedStructMember struct {
	A struct {
		X int
	}
}

type VF_EmbeddedExportedStruct struct {
	VF_SFG
}

type VF_EmbeddedUnexportedStruct struct {
	vf_sFG
}

type VF_TwoEmbeddedStructsWithCancellingMembers struct {
	VF_SFG
	VF_SF
}

type VF_EmbeddedStructsWithSameFieldsAtDifferentDepths struct {
	VF_SFGH3
	VF_SG1
	VF_SFG2
	VF_SF2
	L int
}

type VF_EmbeddedPointerStruct struct {
	*VF_SF
}

type VF_EmbeddedNotAPointer struct {
	VF_M
}

type VF_RecursiveEmbedding = VF_Rec1

type VF_RecursiveEmbedding2 = VF_Rec2

type VF_RecursiveEmbedding3 = VF_RS3

// 	testName: "RecursiveEmbedding3",
// 	val:      RS3{},
// 	expect: []structField{{
// 		name:  "RS2",
// 		index: []int{0},
// 	}, {
// 		name:  "RS1",
// 		index: []int{1},
// 	}, {
// 		name:  "i",
// 		index: []int{1, 0},
// 	}},
// }}

type VF1 struct {
	F string
}

type VF_SFG struct {
	F int
	G int
}

type VF_SFG1 struct {
	VF_SFG
}

type VF_SFG2 struct {
	VF_SFG1
}

type VF_SFGH struct {
	F int
	G int
	H int
}

type VF_SFGH1 struct {
	VF_SFGH
}

type VF_SFGH2 struct {
	VF_SFGH1
}

type VF_SFGH3 struct {
	VF_SFGH2
}

type VF_SF struct {
	F int
}

type VF_SF1 struct {
	VF_SF
}

type VF_SF2 struct {
	VF_SF1
}

type VF_SG struct {
	G int
}

type VF_SG1 struct {
	VF_SG
}

type vf_sFG struct {
	F int
	G int
}

type VF_RS1 struct {
	i int
}

type VF_RS2 struct {
	VF_RS1
}

type VF_RS3 struct {
	VF_RS2
	VF_RS1
}

type VF_M map[string]any

type VF_Rec1 struct {
	*VF_Rec2
}

type VF_Rec2 struct {
	F string
	*VF_Rec1
}
