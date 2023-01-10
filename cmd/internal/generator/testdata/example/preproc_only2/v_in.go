package testdata

type Validator struct {
	F1 string    `pre:"trim"`
	F2 ***string `pre:"trim"`

	S1 *struct {
		S2 struct {
			S3 ***struct {
				F **string `pre:"trim"`
			}
		}
	}
}
