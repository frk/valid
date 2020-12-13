package sk

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SK", "SVK"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+421 123 456 789",
			"+421 123456789",
			"+421123456789",
			"123 456 789",
			"123456789",
		},
		Fail: []string{
			"",
			"+42112345678",
			"+422 123 456 789",
			"+421 023456789",
			"+4211234567892",
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
