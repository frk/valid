package spec

import (
	stdtypes "go/types"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"
)

// Used by tests to convert Error instances to ttError.
//
// Most tests need only to check if the "C ErrorCode" field was set to the
// correct value, the rest of fields' values are less important and so the
// tests need only to check whether the relevant fields were set or not.
type ttError struct {
	C    ErrorCode
	a    *search.AST      `cmp:"+"`
	c    *config.Config   `cmp:"+"`
	ft   *stdtypes.Func   `cmp:"+"`
	rs   *config.RuleSpec //`cmp:"+"`
	rca  *rules.Arg       //`cmp:"+"`
	rcai *int             //`cmp:"+"`
	rcak *string          //`cmp:"+"`
	fp   *types.Var       //`cmp:"+"`
	fpi  *int             //`cmp:"+"`
	err  error            `cmp:"+"`
}

func _ttError(err error) error {
	if e, ok := err.(*Error); ok {
		return (*ttError)(e)
	}
	return err
}

func (e *ttError) Error() string {
	return (*Error)(e).Error()
}
