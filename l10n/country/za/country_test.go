package za

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"ZA", "ZAF"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0821231234",
			"+27821231234",
			"27821231234",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"12345",
			"",
			"ASDFGJKLmZXJtZtesting123",
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
