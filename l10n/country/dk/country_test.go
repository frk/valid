package dk

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"DK", "DNK"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"12345678",
			"12 34 56 78",
			"45 12345678",
			"4512345678",
			"45 12 34 56 78",
			"+45 12 34 56 78",
		},
		Fail: []string{
			"",
			"+45010203",
			"ASDFGJKLmZXJtZtesting123",
			"123456",
			"12 34 56",
			"123 123 12",
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
			"DK31569931",
			"DK31948959",
			"DK38484621",
			"DK31329566",
		},
		Fail: []string{
			"DK31569941",
			"DK31948859",
			"DK38484641",
			"DK31329567",
		},
	}})
}
