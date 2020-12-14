package ir

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"IR", "IRN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+989123456789",
			"989223456789",
			"09323456789",
			"09021456789",
			"+98-990-345-6789",
			"+98 938 345 6789",
			"0938 345 6789",
		},
		Fail: []string{
			"",
			"+989623456789",
			"+981123456789",
			"01234567890",
			"09423456789",
			"09823456789",
			"9123456789",
			"091234567890",
			"0912345678",
			"+98 912 3456 6789",
			"0912 345 678",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"4351666456",
			"5614736867",
		},
		Fail: []string{
			"43516 6456",
			"123443516 6456",
			"891123",
		},
	}})
}
