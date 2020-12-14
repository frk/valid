package kz

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"KZ", "KAZ"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+77254716212",
			"77254716212",
			"87254716212",
			"7254716212",
		},
		Fail: []string{
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
