package bh

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BH", "BHR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+97335078110",
			"+97339534385",
			"+97366331055",
			"+97333146000",
			"97335078110",
			"35078110",
			"66331055",
		},
		Fail: []string{
			"12345",
			"+973350781101",
			"+97379534385",
			"+973035078110",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
			"035078110",
			"16331055",
			"hello",
			"+9733507811a",
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
