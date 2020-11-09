package testdata

type NestedFieldsWithRulesAndNilGuardValidator struct {
	// nested field with rules & nil guard
	G2 struct {
		F1 *string  `is:"email"`
		F2 **string `is:"email"`
		G3 *struct {
			F3 ***string `is:"hex,len:8:128"`
		}
		G4 ***struct {
			G5 **struct {
				F3 **string `is:"required,prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
			}
		}
	}
}
