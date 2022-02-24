package testdata

type Test_ERR_ENUM_NONAME_1_Validator struct {
	F uint `is:"enum"`
}

// cannot be used to declare constants
type enum_kind struct{}

type Test_ERR_ENUM_KIND_1_Validator struct {
	F enum_kind `is:"enum"`
}

// no known consts declared
type enum_noconst uint

type Test_ERR_ENUM_NOCONST_1_Validator struct {
	F enum_noconst `is:"enum"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type enumK uint

const (
	K1 enumK = iota
	K2
	K3
)

type Test_enum_Validator struct {
	F1 enumK `is:"enum"`
}
