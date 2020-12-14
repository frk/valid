package mt

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MT", "MLT"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+35699000000",
			"+35679000000",
			"99000000",
		},
		Fail: []string{
			"356",
			"+35699000",
			"+35610000000",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"VLT2345",
			"VLT 2345",
			"ATD1234",
			"MSK8723",
		},
		Fail: []string{
			//
		},
	}})
}
