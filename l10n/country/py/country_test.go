package py

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PY", "PRY"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+595991372649",
			"+595992847352",
			"+595993847593",
			"+595994857473",
			"+595995348532",
			"+595996435231",
			"+595981847362",
			"+595982435452",
			"+595983948502",
			"+595984342351",
			"+595985403481",
			"+595986384012",
			"+595971435231",
			"+595972103924",
			"+595973438542",
			"+595974425864",
			"+595975425843",
			"+595976342546",
			"+595961435234",
			"+595963425043",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"65478932",
			"+59599384712",
			"+5959938471234",
			"+595547893218",
			"+591993546843",
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
