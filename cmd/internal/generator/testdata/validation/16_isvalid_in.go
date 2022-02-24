package testdata

import (
	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

type T16Validator struct {
	F1 mypkg.MyString `is:"required"`
	F2 *mypkg.MyInt
	F3 **mypkg.MyInt

	// anon
	F4 *interface {
		IsValid() bool
	}
}
