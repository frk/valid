package testdata

type Validator struct {
	F1 int
	F2 int `is:"gt:&F1"`

	S1 struct {
		S2 struct {
			F1 int
			F2 int `is:"gt:&S1.S2.F1"`
			F3 int `is:"gt:&F1"`
		}
		F2 int `is:"lt:&S1.S2.F1"`
	}

	S2 []struct {
		S2 []struct {
			X1 int `is:"gt:&S1.F2"`
			X2 int `is:"gt:&F1"`
		}
		S3 struct {
			Y1 int
			Y2 int `is:"gt:&S1.F2"`
		}

		Z1 [][]int `is:"[][]gt:&F1"`
	}
}
