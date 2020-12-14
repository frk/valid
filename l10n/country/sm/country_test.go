package sm

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SM", "SMR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"612345",
			"05496123456",
			"+37861234567",
			"+390549612345678",
			"+37805496123456789",
		},
		Fail: []string{
			"61234567890",
			"6123",
			"1234567",
			"+49123456",
			"NotANumber",
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
