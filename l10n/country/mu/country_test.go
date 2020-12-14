package mu

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MU", "MUS"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+23012341234",
			"12341234",
			"012341234",
		},
		Fail: []string{
			"41234",
			"",
			"+230",
			"+2301",
			"+23012",
			"+230123",
			"+2301234",
			"+23012341",
			"+230123412",
			"+2301234123",
			"+230123412341",
			"+2301234123412",
			"+23012341234123",
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
