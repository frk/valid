package testdata

type Test_ERR_ARG_BADCMP_1_Validator struct {
	F int `is:"eq:42:64:foo:-22"`
}

type Test_ERR_ARG_BADCMP_2_Validator struct {
	F int `is:"eq:123:&S.G"`
	S struct {
		G string
	}
}

type Test_ERR_ARG_BADCMP_3_Validator struct {
	F int `is:"eq:0.03"`
}

type Test_ERR_ARG_BADCMP_4_Validator struct {
	F uint `is:"ne:-345"`
}

type Test_ERR_ARG_BADCMP_5_Validator struct {
	F bool `is:"eq:1"`
}

type Test_ERR_ARG_BADCMP_6_Validator struct {
	F int `is:"eq:true"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_comparable_Validator struct {
	// boolean rule arguments
	Bool1 bool   `is:"eq:true"`
	Bool2 t_bool `is:"ne:false"`

	// string rule arguments
	String1 string      `is:"eq:foo"`
	String2 t_string    `is:"ne:bar"`
	String3 []byte      `is:"eq:foo"`
	String4 []uint8     `is:"ne:bar"`
	String5 []rune      `is:"eq:foo"`
	String6 []int32     `is:"eq:bar"`
	String7 interface{} `is:"eq:foo"`

	// positive integer rule arguments
	Int1 int     `is:"eq:123456789"`
	Int2 uint    `is:"ne:987654321"`
	Int3 float32 `is:"eq:123456789"`
	Int4 float64 `is:"ne:987654321"`
	Int5 t_int   `is:"eq:123456789"`
	Int6 t_uint  `is:"ne:987654321"`
	Int7 t_float `is:"eq:123456789"`

	// negative integer rule arguments
	NegInt1 int     `is:"eq:-123456789"`
	NegInt2 float32 `is:"eq:-123456789"`
	NegInt3 float64 `is:"ne:-987654321"`
	NegInt4 t_int   `is:"eq:-123456789"`
	NegInt5 t_float `is:"eq:-123456789"`

	// float rule arguments
	Float1 float32 `is:"eq:0.123"`
	Float2 float64 `is:"eq:0.123"`
	Float3 t_float `is:"ne:-0.123"`

	// field rule arguments
	Field1 float32 `is:"eq:&Float1"`
	Field2 t_int   `is:"ne:&Int5"`
	Field3 t_bool  `is:"eq:&Bool1"`
	Field4 string  `is:"eq:&String2"`
	Field5 []byte  `is:"eq:&String1"`
}
