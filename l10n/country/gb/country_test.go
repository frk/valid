package gb

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GB", "GBR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"447789345856",
			"+447861235675",
			"07888814488",
		},
		Fail: []string{
			"67699567",
			"0773894868",
			"077389f8688",
			"+07888814488",
			"0152456999",
			"442073456754",
			"+443003434751",
			"05073456754",
			"08001123123",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"TW8 9GS",
			"BS98 1TL",
			"DE99 3GG",
			"DE55 4SW",
			"DH98 1BT",
			"DH99 1NS",
			"GIR0aa",
			"SA99",
			"W1N 4DJ",
			"AA9A 9AA",
			"AA99 9AA",
			"BS98 1TL",
			"DE993GG",
		},
		Fail: []string{
			//
		},
	}})
}
