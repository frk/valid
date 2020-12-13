package ar

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AR", "ARG"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"5491143214321",
			"+5491143214321",
			"+5492414321432",
			"5498418432143",
		},
		Fail: []string{
			"1143214321",
			"91143214321",
			"+91143214321",
			"549841004321432",
			"549 11 43214321",
			"549111543214321",
			"5714003425432",
			"549114a214321",
			"54 9 11 4321-4321",
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
