package de

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"DE", "DEU"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+49015123456789",
			"+4915123456789",
			"+4930405044550",
			"015123456789",
			"15123456789",
			"15623456789",
			"15623456789",
			"1601234567",
			"16012345678",
			"1621234567",
			"1631234567",
			"1701234567",
			"17612345678",
			"15345678910",
			"15412345678",
		},
		Fail: []string{
			"34412345678",
			"14412345678",
			"16212345678",
			"1761234567",
			"16412345678",
			"17012345678",
			"+4912345678910",
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
