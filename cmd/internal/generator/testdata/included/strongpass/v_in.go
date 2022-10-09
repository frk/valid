package testdata

import (
	"github.com/frk/valid"
)

type Validator struct {
	F1 string  `is:"strongpass"`
	F2 *string `is:"strongpass"`
	F3 string  `is:"strongpass:&pwOpts"`

	pwOpts *valid.StrongPasswordOpts
}
