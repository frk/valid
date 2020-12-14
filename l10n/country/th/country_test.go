package th

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"TH", "THA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0912345678",
			"+66912345678",
			"66912345678",
		},
		Fail: []string{
			"99123456789",
			"12345",
			"67812345623",
			"081234567891",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"10250",
			"72170",
			"12140",
		},
		Fail: []string{
			"T1025",
			"T72170",
			"12140TH",
		},
	}})
}
