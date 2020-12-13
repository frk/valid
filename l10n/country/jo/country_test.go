package jo

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"JO", "JOR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0796477263",
			"0777866254",
			"0786725261",
			"+962796477263",
			"+962777866254",
			"+962786725261",
			"962796477263",
			"962777866254",
			"962786725261",
		},
		Fail: []string{
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
