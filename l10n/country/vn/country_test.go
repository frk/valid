package vn

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"VN", "VNM"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0336012403",
			"+84586012403",
			"84981577798",
			"0708001240",
			"84813601243",
			"0523803765",
			"0863803732",
			"0883805866",
			"0892405867",
			"+84888696413",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"260976684590",
			"01678912345",
			"+841698765432",
			"841626543219",
			"0533803765",
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
