package testdata

type T57Validator struct {
	// nested field with rules & nil guard
	G2 struct {
		F1 *string  `is:"notnil,email"`
		F2 **string `is:"email,notnil"`
		G3 *struct {
			F3 ***string `is:"notnil,hex,len:8:128"`
		} `is:"notnil"`
		G4 ***struct {
			G5 **struct {
				F3 **string `is:"notnil,prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
			} `is:"notnil"`
		} `is:"notnil"`
	}
}
