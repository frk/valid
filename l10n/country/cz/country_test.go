package cz

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"CZ", "CZE"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+420 123 456 789",
			"+420 123456789",
			"+420123456789",
			"123 456 789",
			"123456789",
		},
		Fail: []string{
			"",
			"+42012345678",
			"+421 123 456 789",
			"+420 023456789",
			"+4201234567892",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"20134",
			"392 90",
			"39919",
			"938 29",
			"39949",
		},
		Fail: []string{
			//
		},
	}})
}
