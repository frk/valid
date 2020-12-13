package kw

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"KW", "KWT"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"96550000000",
			"96560000000",
			"96590000000",
			"+96550000000",
			"+96550000220",
			"+96551111220",
		},
		Fail: []string{
			"+96570000220",
			"00962786725261",
			"00962796477263",
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
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
