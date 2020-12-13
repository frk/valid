package ae

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AE", "ARE"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+971502674453",
			"+971521247658",
			"+971541255684",
			"+971555454458",
			"+971561498855",
			"+971585215778",
			"971585215778",
			"0585215778",
			"585215778",
		},
		Fail: []string{
			"12345",
			"+971511498855",
			"+9715614988556",
			"+9745614988556",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
		},
	}})
}
