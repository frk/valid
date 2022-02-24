package testdata

type T58Validator struct {
	// nested field with rules & nil guard
	G2 struct {
		F1 *string  `is:"required,email"`
		F2 **string `is:"email,required"`
		G3 *struct {
			F3 ***string `is:"required,hex,len:8:128"`
		} `is:"required"`
		G4 ***struct {
			G5 **struct {
				F3 **string `is:"required,prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
			} `is:"required"`
		} `is:"required"`
	}
}
