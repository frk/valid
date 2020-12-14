package ph

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PH", "PHL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+639275149120",
			"+639275142327",
			"+639003002023",
			"09275149116",
			"09194877624",
		},
		Fail: []string{
			"12112-13-345",
			"12345678901",
			"sx23YW11cyBmZxxXJt123123",
			"010-38238383",
			"966684123123-2590",
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
