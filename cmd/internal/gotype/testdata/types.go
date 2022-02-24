package testdata

import (
	"unsafe"

	"github.com/frk/valid/cmd/internal/gotype/testdata/mypkg"
)

// types of basic kinds
type (
	BoolType       bool
	IntType        int
	int8Type       int8
	int16Type      int16
	int32Type      int32
	int64Type      int64
	UintType       uint
	uint8Type      uint8
	uint16Type     uint16
	uint32Type     uint32
	uint64Type     uint64
	UintptrType    uintptr
	Float32Type    float32
	float64Type    float64
	Complex64Type  complex64
	complex128Type complex128
	StringType     string
)

// slice types
type (
	IntSlice           []IntType
	strSlice           []string
	byteSlice          []byte
	runeSlice          []rune
	unsafePointerSlice []unsafe.Pointer
)

// array types
type (
	IntArray [10]IntType
	strArray [4]string
)

// map types
type (
	Str2BoolMap map[StringType]bool
	str2BoolMap map[string]BoolType
)

// pointer types
type (
	StrPointer *StringType
	strPointer *string
)

// interface types
type (
	IfaceType interface {
		Method1()
		method2()
	}
	ifaceType interface {
		Method(string) bool
	}
)

// func types
type (
	FuncType func(string) bool
	funcType func(int, string, StringType) (int, error)
)

// chan types
type (
	ChanType chan StringType
	chanType <-chan string
)

// struct types
type (
	StructType struct {
		F1 string `is:"some_rule"`
		f2 bool   `foo bar baz`
	}

	structType struct {
		f1 mypkg.ConstType1 `is:"some_rule"`
		f2 mypkg.Struct     `is:"some_rule"`
	}
)

// recursive types
type (
	recursiveSlice []recursiveSlice

	RecursiveStruct struct {
		F1 *RecursiveStruct
	}
)
