package testdata

import (
	"github.com/frk/valid"
)

type T36Validator struct {
	F1 string `is:"strongpass"`
	F2 string `is:"strongpass:nil"`

	opts *valid.StrongPasswordOpts
	F3   string `is:"strongpass:&opts"`

	opts2 valid.StrongPasswordOpts
	F4    string `is:"strongpass:&opts2"`
}
