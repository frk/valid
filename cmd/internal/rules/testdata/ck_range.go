package testdata

import (
	"github.com/frk/valid/cmd/internal/rules/testdata/mypkg"
)

type Test_ERR_RANGE_TYPE_1_Validator struct {
	F string `is:"rng:10:20"`
}

type Test_ERR_RANGE_TYPE_2_Validator struct {
	F []int `is:"rng:10:20"`
}

type Test_ERR_RANGE_NOARG_1_Validator struct {
	F int `is:"rng::"`
}

type Test_ERR_RANGE_NOARG_2_Validator struct {
	F int `is:"rng:10:"`
}

type Test_ERR_RANGE_NOARG_3_Validator struct {
	F int `is:"rng::20"`
}

type Test_ERR_RANGE_ARGTYPE_1_Validator struct {
	F int `is:"rng:foo:20"`
}

type Test_ERR_RANGE_ARGTYPE_2_Validator struct {
	F int `is:"rng:3.14:20"`
}

type Test_ERR_RANGE_ARGTYPE_3_Validator struct {
	F uint16 `is:"rng:-10:20"`
}

type Test_ERR_RANGE_ARGTYPE_4_Validator struct {
	F int `is:"rng:&S.F:20"`
	S struct{ F string }
}

type Test_ERR_RANGE_BOUNDS_1_Validator struct {
	F int `is:"rng:20:10"`
}

type Test_ERR_RANGE_BOUNDS_2_Validator struct {
	F int `is:"rng:10:10"`
}

type Test_ERR_RANGE_BOUNDS_3_Validator struct {
	F int `is:"rng:10:-20"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

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
