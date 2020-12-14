package in

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"IN", "IND"}, testutil.List{{
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
			"364240",
			"360005",
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
