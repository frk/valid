package ad

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AD", "AND"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+376312345",
			"312345",
		},
		Fail: []string{
			"31234",
			"31234567",
			"512345",
			"NotANumber",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"AD100",
			"AD200",
			"AD300",
			"AD400",
			"AD500",
			"AD600",
			"AD700",
		},
	}})
}
