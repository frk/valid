package testdata

import (
	"github.com/frk/isvalid/internal/testdata/mypkg"
)

type IsValiderValidator struct {
	F1 mypkg.MyString `is:"required"`
	F2 ***mypkg.MyInt

	// anon
	F3 *interface {
		IsValid() bool
	}
}
