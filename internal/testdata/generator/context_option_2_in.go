package testdata

import (
	"time"
)

type ContextOption2Validator struct {
	Context string
	F1      *time.Duration `is:"required:@create,enum"`
	F2      *[]string      `is:"required:@create,[]email"`
}
