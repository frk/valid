package co

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"CO", "COL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+573003321235",
			"573003321235",
			"579871235",
			"3003321235",
			"3213321235",
			"3103321235",
			"3253321235",
			"3321235",
			"574321235",
			"5784321235",
			"5784321235",
			"9821235",
			"573011140876",
			"0698345",
		},
		Fail: []string{
			"1234",
			"+57443875615",
			"57309875615",
			"57109834567",
			"5792434567",
			"5702345689",
			"5714003425432",
			"5703013347567",
			"069834567",
			"969834567",
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
