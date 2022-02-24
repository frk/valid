package testdata

import (
	"github.com/frk/valid/cmd/internal/rules/testdata/mypkg"
)

type Test_ERR_LENGTH_NOLEN_1_Validator struct {
	F uint `is:"len:10"`
}

type Test_ERR_LENGTH_NOLEN_2_Validator struct {
	F mypkg.Time `is:"len:10"`
}

type Test_ERR_LENGTH_NORUNE_1_Validator struct {
	F uint `is:"runecount:10"`
}

type Test_ERR_LENGTH_NORUNE_2_Validator struct {
	F mypkg.Time `is:"runecount:10"`
}

type Test_ERR_LENGTH_NORUNE_3_Validator struct {
	F []rune `is:"runecount:10"`
}

type Test_ERR_LENGTH_NORUNE_4_Validator struct {
	F []int32 `is:"runecount:10"`
}

type Test_ERR_LENGTH_NORUNE_5_Validator struct {
	F []uint `is:"runecount:10"`
}

type Test_ERR_LENGTH_NORUNE_6_Validator struct {
	F [21]byte `is:"runecount:10"`
}

type Test_ERR_LENGTH_NOARG_1_Validator struct {
	F string `is:"len:"`
}

type Test_ERR_LENGTH_NOARG_2_Validator struct {
	F string `is:"len::"`
}

type Test_ERR_LENGTH_NOARG_3_Validator struct {
	F string `is:"runecount::"`
}

type Test_ERR_LENGTH_BOUNDS_1_Validator struct {
	F string `is:"len:10:5"`
}

type Test_ERR_LENGTH_BOUNDS_2_Validator struct {
	F string `is:"runecount:100:99"`
}

type Test_ERR_LENGTH_BOUNDS_3_Validator struct {
	F string `is:"len:42:42"`
}

type Test_ERR_LENGTH_ARGTYPE_1_Validator struct {
	F string `is:"len:foo"`
}

type Test_ERR_LENGTH_ARGTYPE_2_Validator struct {
	F string `is:"runecount:123.987"`
}

type Test_ERR_LENGTH_ARGTYPE_3_Validator struct {
	F string `is:"len:-123"`
}

type Test_ERR_LENGTH_ARGTYPE_4_Validator struct {
	F string `is:"runecount:&S.F"`
	S struct{ F bool }
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

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
