package es

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"ES", "ESP"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+34654789321",
			"654789321",
			"+34714789321",
			"714789321",
			"+34744789321",
			"744789321",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+3465478932",
			"65478932",
			"+346547893210",
			"6547893210",
			"+3470478932",
			"7047893210",
			"+34854789321",
			"7547893219",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"01001",
			"52999",
			"27880",
		},
		Fail: []string{
			"123",
			"1234",
			"53000",
			"052999",
			"0123",
			"abcde",
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"ESX9999999R",
			"ESX99999999",
			"ES99999999R",
		},
		Fail: []string{
			"ESX9999999",
			"ES9999999R",
			"ESXR9999999",
			"ES9999999XR",
		},
	}})
}
