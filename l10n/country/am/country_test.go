package am

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AM", "ARM"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+37410324123",
			"+37422298765",
			"+37431276521",
			"022698763",
			"37491987654",
			"+37494567890",
		},
		Fail: []string{
			"12345",
			"+37411498855",
			"+37411498123",
			"05614988556",
			"",
			"37456789000",
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
