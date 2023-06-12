package types

// Kind describes the specific kind of a Go type.
type Kind uint

const (
	// basic
	INVALID Kind = iota

	_basic_start
	BOOL

	_numeric_start // int/uint/float
	_integer_start // int
	INT
	INT8
	INT16
	INT32
	INT64
	_integer_end

	_unsigned_start // uint
	UINT
	UINT8
	UINT16
	UINT32
	UINT64
	UINTPTR
	_unsigned_end

	FLOAT32
	FLOAT64
	_numeric_end

	COMPLEX64
	COMPLEX128
	STRING
	UNSAFEPOINTER
	_basic_end

	// non-basic
	ARRAY
	INTERFACE
	MAP
	PTR
	SLICE
	STRUCT
	CHAN // don't validate
	FUNC // don't validate

	UNION

	// alisases (basic)
	BYTE = UINT8
	RUNE = INT32
)

// Reports whether or not k is of a basic kind.
func (k Kind) IsBasic() bool { return _basic_start < k && k < _basic_end }

// Reports whether or not k is of the numeric kind, note that this
// does not include the complex64 and complex128 kinds.
func (k Kind) IsNumeric() bool { return _numeric_start < k && k < _numeric_end }

// Reports whether or not k is one of the int types. (does not include unsigned integers)
func (k Kind) IsInteger() bool { return _integer_start < k && k < _integer_end }

// Reports whether or not k is one of the uint types.
func (k Kind) IsUnsigned() bool { return _unsigned_start < k && k < _unsigned_end }

// Reports whether or not k is one of the float types.
func (k Kind) IsFloat() bool { return k == FLOAT32 || k == FLOAT64 }

// Reports whether or not k is one of the complex types.
func (k Kind) IsComplex() bool { return k == COMPLEX64 || k == COMPLEX128 }

// BasicString returns a string representation of the basic kind k.
func (k Kind) BasicString() string {
	if k.IsBasic() {
		return _kindstring[k]
	}
	return "<unknown>"
}

func (k Kind) String() string {
	if int(k) < len(_kindstring) {
		return _kindstring[k]
	}
	return "<unknown>"
}

var _kindstring = [...]string{
	INVALID:    "<invalid>",
	BOOL:       "bool",
	INT:        "int",
	INT8:       "int8",
	INT16:      "int16",
	INT32:      "int32",
	INT64:      "int64",
	UINT:       "uint",
	UINT8:      "uint8",
	UINT16:     "uint16",
	UINT32:     "uint32",
	UINT64:     "uint64",
	UINTPTR:    "uintptr",
	FLOAT32:    "float32",
	FLOAT64:    "float64",
	COMPLEX64:  "complex64",
	COMPLEX128: "complex128",
	STRING:     "string",

	// ...
	ARRAY:     "array",
	INTERFACE: "interface",
	MAP:       "map",
	PTR:       "ptr",
	SLICE:     "slice",
	STRUCT:    "struct",
	CHAN:      "chan",
	FUNC:      "func",
}
