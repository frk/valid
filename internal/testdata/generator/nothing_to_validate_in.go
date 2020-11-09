package testdata

type NothingToValidateValidator struct {
	// no rule, no subfields, nothing to do
	F1 string
	F2 interface{}
	F3 ***string

	// nested no-rule fields
	G0 struct {
		F1 string
		F2 interface{}
		GA struct {
			F3 ***string
		}
	}
}
