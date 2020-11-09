// DO NOT EDIT. This file was generated by "github.com/frk/isvalid".

package testdata

func (v ErrorConstructorValidator) Validate() error {
	if len(v.F1) == 0 {
		return v.ec.Error("F1", v.F1, "required")
	} else if v.F1 != "foo" {
		return v.ec.Error("F1", v.F1, "eq", "foo")
	}
	return nil
}