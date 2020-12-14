package lt

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"LT", "LTU"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+37051234567",
			"851234567",
		},
		Fail: []string{
			"+65740 123 456",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"740123456",
			"+65640123456",
			"+64210123456",
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
