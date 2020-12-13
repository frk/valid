package us

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"US", "USA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"19876543210",
			"8005552222",
			"+15673628910",
			"+1(567)3628910",
			"+1(567)362-8910",
			"+1(567) 362-8910",
			"1(567)362-8910",
			"1(567)362 8910",
			"223-456-7890",
		},
		Fail: []string{
			"564785",
			"0123456789",
			"1437439210",
			"+10345672645",
			"11435213543",
			"1(067)362-8910",
			"1(167)362-8910",
			"+2(267)362-8910",
			"+3365520145",
		},
	}})
}
