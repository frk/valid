package mx

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MX", "MEX"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+52019654789321",
			"+52199654789321",
			"+5201965478932",
			"+5219654789321",
			"52019654789321",
			"52199654789321",
			"5201965478932",
			"5219654789321",
			"87654789321",
			"8654789321",
			"0187654789321",
			"18654789321",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+3465478932",
			"65478932",
			"+346547893210",
			"+34704789321",
			"704789321",
			"+34754789321",
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
