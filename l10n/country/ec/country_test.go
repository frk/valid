package ec

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"EC", "ECU"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+593987654321",
			"593987654321",
			"0987654321",
			"027332615",
			"+59323456789",
		},
		Fail: []string{
			"03321321",
			"+593387561",
			"59312345677",
			"02344635",
			"593123456789",
			"081234567",
			"+593912345678",
			"+593902345678",
			"+593287654321",
			"593287654321",
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
