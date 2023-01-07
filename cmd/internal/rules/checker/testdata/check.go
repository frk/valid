package testdata

import (
	"github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg"
)

type Test_E_FIELD_UNKNOWN_1_Validator struct {
	F int `is:"gt:&num"`
}

type Test_E_FIELD_UNKNOWN_2_Validator struct {
	F string `pre:"p4:&num"`
}

type Test_E_FIELD_UNKNOWN_3_Validator struct {
	num int
	X   struct {
		F int `is:"gt:.num"`
	}
}

type Test_E_FIELD_UNKNOWN_4_Validator struct {
	num int
	X   struct {
		F string `pre:"p4:.num"`
	}
}

type Test_E_RULE_ARGMIN_1_Validator struct {
	F int `is:"gt"`
}

type Test_E_RULE_ARGMIN_2_Validator struct {
	F string `pre:"p4"`
}

type Test_E_RULE_ARGMAX_1_Validator struct {
	F int `is:"gt:4:5"`
}

type Test_E_RULE_ARGMAX_2_Validator struct {
	F string `pre:"p4:1:2:3"`
}

type Test_E_PREPROC_INVALID_1_Validator struct {
	F string `pre:"p0"`
}

// -----------------------------------------------------------------------------
// "required" test cases

type Test_E_NOTNIL_TYPE_1_Validator struct {
	F string `is:"notnil"`
}

type Test_E_NOTNIL_TYPE_2_Validator struct {
	F bool `is:"notnil"`
}

type Test_E_NOTNIL_TYPE_3_Validator struct {
	F float64 `is:"notnil"`
}

// -----------------------------------------------------------------------------
// "comparable" test cases

type Test_E_ARG_BADCMP_1_Validator struct {
	F int `is:"eq:42:64:foo:-22"`
}

type Test_E_ARG_BADCMP_2_Validator struct {
	F int `is:"eq:123:&S.G"`
	S struct {
		G string
	}
}

type Test_E_ARG_BADCMP_3_Validator struct {
	F int `is:"eq:0.03"`
}

type Test_E_ARG_BADCMP_4_Validator struct {
	F uint `is:"ne:-345"`
}

type Test_E_ARG_BADCMP_5_Validator struct {
	F bool `is:"eq:1"`
}

type Test_E_ARG_BADCMP_6_Validator struct {
	F int `is:"eq:true"`
}

// -----------------------------------------------------------------------------
// "ordered" test cases

type Test_E_ORDERED_TYPE_1_Validator struct {
	F []string `is:"min:8"`
}

type Test_E_ORDERED_TYPE_2_Validator struct {
	F []int `is:"gt:8"`
}

type Test_E_ORDERED_ARGTYPE_1_Validator struct {
	F int `is:"gte:0.8"`
}

type Test_E_ORDERED_ARGTYPE_2_Validator struct {
	F float64 `is:"gte:foo"`
}

type Test_E_ORDERED_ARGTYPE_3_Validator struct {
	F string `is:"lte:&S.F"`
	S struct{ F float32 }
}

// -----------------------------------------------------------------------------
// "length" test cases

type Test_E_LENGTH_NOLEN_1_Validator struct {
	F uint `is:"len:10"`
}

type Test_E_LENGTH_NOLEN_2_Validator struct {
	F mypkg.Time `is:"len:10"`
}

type Test_E_LENGTH_NORUNE_1_Validator struct {
	F uint `is:"runecount:10"`
}

type Test_E_LENGTH_NORUNE_2_Validator struct {
	F mypkg.Time `is:"runecount:10"`
}

type Test_E_LENGTH_NORUNE_3_Validator struct {
	F rune `is:"runecount:10"`
}

type Test_E_LENGTH_NORUNE_4_Validator struct {
	F []int16 `is:"runecount:10"`
}

type Test_E_LENGTH_NORUNE_5_Validator struct {
	F []uint `is:"runecount:10"`
}

type Test_E_LENGTH_NORUNE_6_Validator struct {
	F [21]byte `is:"runecount:10"`
}

type Test_E_LENGTH_NOARG_1_Validator struct {
	F string `is:"len:"`
}

type Test_E_LENGTH_NOARG_2_Validator struct {
	F string `is:"len::"`
}

type Test_E_LENGTH_NOARG_3_Validator struct {
	F string `is:"runecount::"`
}

type Test_E_LENGTH_BOUNDS_1_Validator struct {
	F string `is:"len:10:5"`
}

type Test_E_LENGTH_BOUNDS_2_Validator struct {
	F string `is:"runecount:100:99"`
}

type Test_E_LENGTH_BOUNDS_3_Validator struct {
	F string `is:"len:42:42"`
}

type Test_E_LENGTH_ARGTYPE_1_Validator struct {
	F string `is:"len:foo"`
}

type Test_E_LENGTH_ARGTYPE_2_Validator struct {
	F string `is:"runecount:123.987"`
}

type Test_E_LENGTH_ARGTYPE_3_Validator struct {
	F string `is:"len:-123"`
}

type Test_E_LENGTH_ARGTYPE_4_Validator struct {
	F string `is:"runecount:&S.F"`
	S struct{ F bool }
}

// -----------------------------------------------------------------------------
// "range" test cases

type Test_E_RANGE_TYPE_1_Validator struct {
	F string `is:"rng:10:20"`
}

type Test_E_RANGE_TYPE_2_Validator struct {
	F []int `is:"rng:10:20"`
}

type Test_E_RANGE_NOARG_1_Validator struct {
	F int `is:"rng::"`
}

type Test_E_RANGE_NOARG_2_Validator struct {
	F int `is:"rng:10:"`
}

type Test_E_RANGE_NOARG_3_Validator struct {
	F int `is:"rng::20"`
}

type Test_E_RANGE_ARGTYPE_1_Validator struct {
	F int `is:"rng:foo:20"`
}

type Test_E_RANGE_ARGTYPE_2_Validator struct {
	F int `is:"rng:3.14:20"`
}

type Test_E_RANGE_ARGTYPE_3_Validator struct {
	F uint16 `is:"rng:-10:20"`
}

type Test_E_RANGE_ARGTYPE_4_Validator struct {
	F int `is:"rng:&S.F:20"`
	S struct{ F string }
}

type Test_E_RANGE_BOUNDS_1_Validator struct {
	F int `is:"rng:20:10"`
}

type Test_E_RANGE_BOUNDS_2_Validator struct {
	F int `is:"rng:10:10"`
}

type Test_E_RANGE_BOUNDS_3_Validator struct {
	F int `is:"rng:10:-20"`
}

// -----------------------------------------------------------------------------
// "enum" test cases

type Test_E_ENUM_NONAME_1_Validator struct {
	F uint `is:"enum"`
}

// cannot be used to declare constants
type enum_kind struct{}

type Test_E_ENUM_KIND_1_Validator struct {
	F enum_kind `is:"enum"`
}

// no known consts declared
type enum_noconst uint

type Test_E_ENUM_NOCONST_1_Validator struct {
	F enum_noconst `is:"enum"`
}

// -----------------------------------------------------------------------------
// "function" test cases

type Test_E_FUNCTION_INTYPE_1_Validator struct {
	F float32 `is:"contains:foo"`
}

type Test_E_FUNCTION_ARGTYPE_1_Validator struct {
	F string `is:"uuid:v6"`
}

type Test_E_FUNCTION_ARGVALUE_1_Validator struct {
	F string `is:"alpha:foo"`
}

type Test_E_FUNCTION_ARGVALUE_2_Validator struct {
	F string `is:"alnum:foo"`
}

type Test_E_FUNCTION_ARGVALUE_3_Validator struct {
	F string `is:"ccy:foo"`
}

type Test_E_FUNCTION_ARGVALUE_4_Validator struct {
	F string `is:"decimal:foo"`
}

type Test_E_FUNCTION_ARGVALUE_5_Validator struct {
	F string `is:"hash:foo"`
}

type Test_E_FUNCTION_ARGVALUE_6_Validator struct {
	F string `is:"ip:5"`
}

type Test_E_FUNCTION_ARGVALUE_7_Validator struct {
	F string `is:"isbn:12"`
}

type Test_E_FUNCTION_ARGVALUE_8_Validator struct {
	F string `is:"iso639:3"`
}

type Test_E_FUNCTION_ARGVALUE_9_Validator struct {
	F string `is:"iso31661a:1"`
}

type Test_E_FUNCTION_ARGVALUE_10_Validator struct {
	F string `is:"mac:7"`
}

type Test_E_FUNCTION_ARGVALUE_11_Validator struct {
	F string `is:"re:[0-9)?"`
}

type Test_E_FUNCTION_ARGVALUE_12_Validator struct {
	F string `is:"uuid:6"`
}

type Test_E_FUNCTION_ARGVALUE_13_Validator struct {
	F string `is:"phone:foo"`
}

type Test_E_FUNCTION_ARGVALUE_14_Validator struct {
	F string `is:"vat:foo"`
}

type Test_E_FUNCTION_ARGVALUE_15_Validator struct {
	F string `is:"zip:foo"`
}

// -----------------------------------------------------------------------------
// "method" test cases

type Test_E_METHOD_TYPE_1_Validator struct {
	F int `is:"isvalid"`
}

// -----------------------------------------------------------------------------
// "optional" test cases

type Test_E_OPTIONAL_CONFLICT_1_Validator struct {
	F *string `is:"optional,required"`
}

type Test_E_OPTIONAL_CONFLICT_2_Validator struct {
	F *string `is:"required,optional"`
}

type Test_E_OPTIONAL_CONFLICT_3_Validator struct {
	F *string `is:"notnil,optional"`
}

type Test_E_OPTIONAL_CONFLICT_4_Validator struct {
	F *string `is:"omitnil,notnil"`
}

// -----------------------------------------------------------------------------
// "preprocessor" test cases

type Test_E_PREPROC_INTYPE_1_Validator struct {
	F string `pre:"p1"`
}

type Test_E_PREPROC_OUTTYPE_1_Validator struct {
	F string `pre:"p2"`
}

type Test_E_PREPROC_ARGTYPE_1_Validator struct {
	F string `pre:"p4:foo"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_checker_Validator struct {
	F1 int
	F2 int
	S1 struct {
		F1 int
		F2 int

		T1 int `is:"gt:.F1"`
		T2 int `is:"gt:&F2"`
		T3 int `is:"gt:.S2.F1"`

		S2 struct {
			F1 int
		}
	}

	T1 int `is:"gt:&F1"`
	T2 int `is:"gt:&S1.F2"`
	T3 int `is:"gt:.S1.F1"`
	T4 int `is:"gt:.S1.S2.F1"`
}

type Test_required_Validator struct {
	F1 string            `is:"required"`
	F2 *string           `is:"notnil"`
	F3 []byte            `is:"notnil"`
	F4 map[string]string `is:"notnil"`
	F5 interface{}       `is:"notnil"`
	F6 int               `is:"required"`
	F7 bool              `is:"required"`
}

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

type Test_ordered_Validator struct {
	F1 int     `is:"gt:8"`
	F2 float64 `is:"lt:0.029"`
	F3 uint16  `is:"gte:128"`
	F4 int16   `is:"lte:-128"`
	F5 int16   `is:"min:8,max:128"`
	F6 string  `is:"lt:foo,min:bar"`
	F7 int     `is:"lt:&S.F1,gt:&S.F2"`
	S  struct {
		F1 int
		F2 float64
	}
}

type Test_length_Validator struct {
	// types
	String1 string              `is:"len:10"`
	String2 string              `is:"runecount:10"`
	String3 mypkg.String        `is:"len:10"`
	String4 mypkg.String        `is:"runecount:10"`
	Bytes1  []byte              `is:"len:10"`
	Bytes2  []byte              `is:"runecount:10"`
	Bytes3  mypkg.Bytes         `is:"len:10"`
	Bytes4  mypkg.Bytes         `is:"runecount:10"`
	Slice1  []struct{}          `is:"len:10"`
	Map1    map[string]struct{} `is:"len:10"`
	// TODO this doesn't really make sense, perhaps "len"
	// should not be allowed to be applied to arrays...
	Array1 [10]struct{} `is:"len:10"`

	// bounds
	Bounds1 string `is:"len:10:"`
	Bounds2 string `is:"len::20"`
	Bounds3 string `is:"len:10:20"`
	Bounds4 string `is:"runecount:10:"`
	Bounds5 string `is:"runecount::20"`
	Bounds6 string `is:"runecount:10:20"`
}

type Test_range_Validator struct {
	Int1   int       `is:"rng:-10:20"`
	Int2   int       `is:"rng:&S.Int:20"`
	MyInt1 mypkg.Int `is:"rng:-10:20"`
	MyInt2 mypkg.Int `is:"rng:&S.MyInt:20"`
	Uint1  int       `is:"rng:10:20"`
	Uint2  int       `is:"rng:&S.Uint:20"`
	Float1 float64   `is:"rng:-1.0:2.0"`
	Float2 float64   `is:"rng:&S.Float:2.0"`
	Float3 float32   `is:"rng:-3:21"`

	S struct {
		Int   int
		MyInt mypkg.Int
		Uint  int
		Float float64
	}
}

type Test_between_Validator struct {
	Int1   int       `is:"between:-10:20"`
	Int2   int       `is:"between:&S.Int:20"`
	MyInt1 mypkg.Int `is:"between:-10:20"`
	MyInt2 mypkg.Int `is:"between:&S.MyInt:20"`
	Uint1  int       `is:"between:10:20"`
	Uint2  int       `is:"between:&S.Uint:20"`
	Float1 float64   `is:"between:-1.0:2.0"`
	Float2 float64   `is:"between:&S.Float:2.0"`
	Float3 float32   `is:"between:-3:21"`

	S struct {
		Int   int
		MyInt mypkg.Int
		Uint  int
		Float float64
	}
}
type enumK uint

const (
	K1 enumK = iota
	K2
	K3
)

type Test_enum_Validator struct {
	F1 enumK `is:"enum"`
}

type Test_function_Validator struct {
	Fn1 string `is:"contains:foo"`
	Fn2 string `is:"contains:foo:bar:baz"`

	UUID1 string `is:"uuid"`
	UUID2 string `is:"uuid:3"`
	UUID3 string `is:"uuid:v3"`
	UUID4 string `is:"uuid:4"`
	UUID5 string `is:"uuid:v4"`
	UUID6 string `is:"uuid:5"`
	UUID7 string `is:"uuid:v5"`

	R8 string `is:"r8:&helper"`
	R9 string `is:"r9:&helper2"`

	helper  *mypkg.CheckHelper
	helper2 *mypkg.CheckWithErrorHelper
}

type Typ struct{}

func (*Typ) IsValid() bool { return true }

type Test_method_Validator struct {
	F1 Typ  `is:"isvalid"`
	F2 *Typ `is:"isvalid"`
}

type Test_optional_Validator struct {
	F1 *string `is:"optional"`
	F2 *string `is:"omitnil"`
	F3 *string `is:"required"`
	F4 *string `is:"notnil"`
	F5 *string `is:"omitnil,required"`
}

type Test_preproc_Validator struct {
	F1 string `pre:"trim"`
}
