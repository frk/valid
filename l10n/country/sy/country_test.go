package sy

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SY", "SYR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0944549710",
			"+963944549710",
			"956654379",
			"0944549710",
			"0962655597",
		},
		Fail: []string{
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
