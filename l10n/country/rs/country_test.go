package rs

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"RS", "SRB"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0640133338",
			"063333133",
			"0668888878",
			"+381645678912",
			"+381611314000",
			"0655885010",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
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
