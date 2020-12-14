package bo

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BO", "BOL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+59175553635",
			"+59162223685",
			"+59179783890",
			"+59160081890",
			"79783890",
			"60081890",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+59199783890",
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
