package rules

import (
	"go/types"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"
)

// Used by tests to convert Error instances to ttError.
//
// Most tests need only to check if the "C ErrorCode" field was set to the
// correct value, the rest of fields' values are less important and so the
// tests need only to check whether the relevant fields were set or not.
type ttError struct {
	C    ErrorCode
	a    *search.AST         `cmp:"+"`
	c    *config.Config      `cmp:"+"`
	rc   *config.RuleConfig  //`cmp:"+"`
	rs   *config.RuleSpec    //`cmp:"+"`
	rca  *Arg                //`cmp:"+"`
	rcai *int                //`cmp:"+"`
	rcak *string             //`cmp:"+"`
	ft   *types.Func         `cmp:"+"`
	sf   *gotype.StructField //`cmp:"+"`
	sfv  *types.Var          `cmp:"+"`
	ty   *gotype.Type        //`cmp:"+"`
	tag  *Tag                //`cmp:"+"`
	r    *Rule               //`cmp:"+"`
	r2   *Rule               //`cmp:"+"`
	ra   *Arg                //`cmp:"+"`
	raf  *gotype.StructField `cmp:"+"`
	fp   *gotype.Var         //`cmp:"+"`
	fpi  *int                //`cmp:"+"`
	err  error               `cmp:"+"`
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
