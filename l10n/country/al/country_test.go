package al

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AL", "ALB"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"067123456",
			"+35567123456",
		},
		Fail: []string{
			"67123456",
			"06712345",
			"06712345678",
			"065123456",
			"057123456",
			"NotANumber",
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
