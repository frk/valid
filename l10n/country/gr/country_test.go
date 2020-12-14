package gr

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GR", "GRC"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+306944848966",
			"6944848966",
			"306944848966",
		},
		Fail: []string{
			"2102323234",
			"+302646041461",
			"120000000",
			"20000000000",
			"68129485729",
			"6589394827",
			"298RI89572",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"022 93",
			"29934",
			"90293",
			"299 42",
			"94944",
		},
		Fail: []string{
			//
		},
	}})
}
