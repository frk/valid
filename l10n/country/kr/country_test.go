package kr

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"KR", "KOR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+82-010-1234-5678",
			"+82-10-1234-5678",
			"82-010-1234-5678",
			"82-10-1234-5678",
			"+82 10 1234 5678",
			"010-123-5678",
			"10-1234-5678",
			"+82 10 1234 5678",
			"011 1234 5678",
			"+820112345678",
			"01012345678",
			"+82 016 1234 5678",
			"82 19 1234 5678",
			"+82 010 12345678",
		},
		Fail: []string{
			"abcdefghi",
			"+82 10 1234 567",
			"+82 10o 1234 1234",
			"+82 101 1234 5678",
			"+82 10 12 5678",
			"+011 7766 1234",
			"011_7766_1234",
			"+820 11 7766 1234",
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
