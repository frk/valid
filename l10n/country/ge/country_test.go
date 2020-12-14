package ge

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GE", "GEO"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+99550001111",
			"+99551535213",
			"+995798526662",
			"798526662",
			"50001111",
			"798526662",
			"+995799766525",
		},
		Fail: []string{
			"+995500011118",
			"+9957997665250",
			"+995999766525",
			"20000000000",
			"68129485729",
			"6589394827",
			"298RI89572",
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
