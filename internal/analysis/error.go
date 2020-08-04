package analysis

import (
	"fmt"
)

type anError struct {
	Code        errorCode
	StructField *StructField
	Rule        *Rule
	RuleParam   *RuleParam
	Err         error
}

func (e *anError) Error() string {
	return fmt.Sprintf("errorCode: #%d", e.Code)
}

type errorCode uint8

const (
	_ errorCode = iota
	errRuleUnknown
	errRuleContext
	errRuleParamKind
	errRuleParamTypeUint
	errRuleParamTypeNint
	errRuleParamTypeFloat
	errRuleParamTypeString
	errRuleParamValueRegexp
	errRuleParamValueUUID
	errRuleParamValueIP
	errRuleParamValueMAC
	errRuleParamValueCountryCode
	errTypeLength
	errTypeNumeric
	errTypeString
	errFieldKeyUnknown
	errFieldKeyConflict
)
