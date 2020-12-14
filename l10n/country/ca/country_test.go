package ca

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"CA", "CAN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"19876543210",
			"8005552222",
			"+15673628910",
		},
		Fail: []string{
			"564785",
			"0123456789",
			"1437439210",
			"+10345672645",
			"11435213543",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"L4T 0A5",
			"G1A-0A2",
			"A1A 1A1",
			"X0A-0H0",
			"V5K 0A1",
		},
		Fail: []string{
			//
		},
	}})
}
