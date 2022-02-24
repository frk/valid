package testdata

import (
	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

type T17Validator struct {
	F1 mypkg.MyString `is:"required,-isvalid"`
	F2 ***mypkg.MyInt `is:"-isvalid"`

	// anon
	F3 *interface {
		IsValid() bool
	} `is:"-isvalid"`
}
