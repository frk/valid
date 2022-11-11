package testdata

type Validator struct {
	G1 struct {
		F4 string `is:"email"`
		GA struct {
			F4 string `is:"email"`
			F5 string `is:"hex,len:8:128"`
			GB struct {
				F4a string `is:"email"`
				GC  struct {
					F6 string `is:"prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
				}
				F4b string `is:"email"`
				F5  string `is:"hex,len:8:128"`
			}
		}
		F6 string `is:"prefix:foo,contains:bar,suffix:baz:quux,len:8:64"`
	}
}
