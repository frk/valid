package il

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"IL", "ISR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"10200",
			"10292",
			"10300",
			"10329",
			"3885500",
			"4290500",
			"4286000",
			"7080000",
		},
		Fail: []string{
			"123",
			"012345",
			"011111",
			"101123",
			"291123",
			"351123",
			"541123",
			"551123",
			"651123",
			"661123",
			"861123",
			"871123",
			"881123",
			"891123",
		},
	}})
}
