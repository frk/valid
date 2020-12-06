package testdata

import (
	"github.com/frk/isvalid"
)

type StrongPasswordValidator struct {
	F1 string `is:"strongpass"`
	F2 string `is:"strongpass:nil"`

	opts *isvalid.StrongPasswordOpts
	F3   string `is:"strongpass:&opts"`

	opts2 isvalid.StrongPasswordOpts
	F4    string `is:"strongpass:&opts2"`
}
