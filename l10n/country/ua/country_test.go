package ua

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"UA", "UKR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+380982345679",
			"380982345679",
			"80982345679",
			"0982345679",
		},
		Fail: []string{
			"+30982345679",
			"982345679",
			"+380 98 234 5679",
			"+380-98-234-5679",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"740123456",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"65000",
			"65080",
			"01000",
		},
		Fail: []string{
			//
		},
	}})
}
