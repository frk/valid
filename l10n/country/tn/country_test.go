package tn

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"TN", "TUN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"23456789",
			"+21623456789",
			"21623456789",
		},
		Fail: []string{
			"12345",
			"75200123",
			"+216512345678",
			"13520459",
			"85479520",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
}
