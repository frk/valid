package analysis

import (
	"github.com/frk/isvalid/internal/testdata"
	"github.com/frk/isvalid/internal/testdata/mypkg"
)

type AnalysisTestOK_Validator struct {
	UserInput *testdata.UserInput `isvalid:"omitkey"`
	Context   string
}

func (v *AnalysisTestOK_Validator) beforevalidate() error {
	return nil
}

func (v AnalysisTestOK_Validator) AfterValidate() error {
	return nil
}

type AnalysisTestOK_ErrorConstructorValidator struct {
	F string `is:"required"`
	mypkg.MyErrorConstructor
}

type AnalysisTestOK_ErrorAggregatorValidator struct {
	F      string `is:"required"`
	erragg mypkg.MyErrorAggregator
}

type AnalysisTestOK_ContextValidator struct {
	F       string `is:"required"`
	context string
}

type AnalysisTestOK_Context2Validator struct {
	F       string `is:"required"`
	Context string
}
