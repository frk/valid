// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

func (v T42Validator) Validate() error {
	if v.F1 == "" {
		v.ea.Error("F1", v.F1, "required")
	} else if v.F1 != "foo" {
		v.ea.Error("F1", v.F1, "eq", "foo")
	}
	return v.ea.Out()
}
