// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
	"fmt"

	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

func (v T02Validator) Validate() error {
	if !mypkg.HasUniqueInts(v.F1) {
		return errors.New("F1 is not valid")
	}
	if !mypkg.HasUniqueInts(v.F2, v.F1) {
		return errors.New("F2 is not valid")
	}
	for _, e1 := range v.F3 {
		if e1 != nil {
			for _, e2 := range e1.X1 {
				if e2 <= e1.X2 {
					return fmt.Errorf("F3.X1 must be greater than: %v", e1.X2)
				}
			}
		}
	}
	return nil
}
