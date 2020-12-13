package ma

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MA", "MAR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0522714782",
			"0690851123",
			"0708186135",
			"+212522714782",
			"+212690851123",
			"+212708186135",
			"00212522714782",
			"00212690851123",
			"00212708186135",
		},
		Fail: []string{
			"522714782",
			"690851123",
			"708186135",
			"212522714782",
			"212690851123",
			"212708186135",
			"0212522714782",
			"0212690851123",
			"0212708186135",
			"",
			"12345",
			"0922714782",
			"+212190851123",
			"00212408186135",
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
