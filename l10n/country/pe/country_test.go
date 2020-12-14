package pe

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PE", "PER"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+51912232764",
			"+51923464567",
			"+51968267382",
			"+51908792973",
			"974980472",
			"908792973",
			"+51974980472",
		},
		Fail: []string{
			"999",
			"+51812232764",
			"+5181223276499",
			"+25589032",
			"123456789",
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
