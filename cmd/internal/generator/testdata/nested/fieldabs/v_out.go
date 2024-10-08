// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"fmt"
)

func (v Validator) Validate() error {
	if v.F2 <= v.F1 {
		return fmt.Errorf("F2 must be greater than: %v", v.F1)
	}
	if v.S1.S2.F2 <= v.S1.S2.F1 {
		return fmt.Errorf("S1.S2.F2 must be greater than: %v", v.S1.S2.F1)
	}
	if v.S1.S2.F3 <= v.F1 {
		return fmt.Errorf("S1.S2.F3 must be greater than: %v", v.F1)
	}
	if v.S1.F2 >= v.S1.S2.F1 {
		return fmt.Errorf("S1.F2 must be less than: %v", v.S1.S2.F1)
	}
	for _, e1 := range v.S2 {
		for _, e2 := range e1.S2 {
			if e2.X1 <= v.S1.F2 {
				return fmt.Errorf("S2.S2.X1 must be greater than: %v", v.S1.F2)
			}
			if e2.X2 <= v.F1 {
				return fmt.Errorf("S2.S2.X2 must be greater than: %v", v.F1)
			}
		}
		if e1.S3.Y2 <= v.S1.F2 {
			return fmt.Errorf("S2.S3.Y2 must be greater than: %v", v.S1.F2)
		}
		for _, e2 := range e1.Z1 {
			for _, e3 := range e2 {
				if e3 <= v.F1 {
					return fmt.Errorf("S2.Z1 must be greater than: %v", v.F1)
				}
			}
		}
	}
	return nil
}
