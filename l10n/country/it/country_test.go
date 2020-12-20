package it

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"IT", "ITA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"370 3175423",
			"333202925",
			"+39 310 7688449",
			"+39 3339847632",
		},
		Fail: []string{
			"011 7387545",
			"12345",
			"+45 345 6782395",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"IT07643520567",
		},
		Fail: []string{
			//
		},
	}})
}
