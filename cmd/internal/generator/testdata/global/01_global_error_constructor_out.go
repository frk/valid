// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

func (v T01Validator) Validate() error {
	if v.F1 == "" {
		return mypkg.NewError("F1", v.F1, "required")
	}
	if len(v.F2) == 0 {
		return mypkg.NewError("F2", v.F2, "required")
	} else if len(v.F2) > 9 {
		return mypkg.NewError("F2", v.F2, "len", "", 9)
	}
	return nil
}