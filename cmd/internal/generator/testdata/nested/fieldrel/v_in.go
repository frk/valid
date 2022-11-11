package testdata

type Validator struct {
	F1 int
	F2 int `is:"gt:.F1"`

	S1 struct {
		S2 struct {
			F1 int
			F2 int `is:"gt:.F1"`
		}
		S3 struct {
			F1 int
			F2 int `is:"gt:.F1"`
		}

		F1 int `is:"gt:.S2.F1"`
		F2 int `is:"gt:.S3.F1"`
	}
}
