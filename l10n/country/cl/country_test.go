package cl

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"CL", "CHL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+56733875615",
			"56928590234",
			"0928590294",
			"0208590294",
		},
		Fail: []string{
			"1234",
			"+5633875615",
			"563875615",
			"56109834567",
			"56069834567",
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
