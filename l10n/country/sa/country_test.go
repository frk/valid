package sa

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SA", "SAU"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0556578654",
			"+966556578654",
			"966556578654",
			"596578654",
			"572655597",
		},
		Fail: []string{
			"12345",
			"",
			"+9665626626262",
			"+96633221097",
			"0114152198",
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
