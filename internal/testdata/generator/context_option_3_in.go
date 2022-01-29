package testdata

import (
	"time"
)

type ContextOption3Validator struct {
	Context string
	F1      *time.Duration `is:"required:@create,enum,lt:10000:@foo_bar"`
	F2      *[]string      `is:"required:@create,[]email:@create,len::128:@update"`
}
