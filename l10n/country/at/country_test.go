package at

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AT", "AUT"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+436761234567",
			"06761234567",
			"00436123456789",
			"+436123456789",
			"01999",
			"+4372876",
			"06434908989562345",
		},
		Fail: []string{
			"167612345678",
			"1234",
			"064349089895623459",
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
