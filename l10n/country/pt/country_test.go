package pt

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PT", "PRT"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"4829-489",
			"0294-348",
			"8156-392",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"PT999999990",
			"PT501442600",
		},
		Fail: []string{
			"PT999999999",
		},
	}})
}
