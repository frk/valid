package mypkg

import (
	"time"
)

type Int int
type String string
type Bytes []byte
type Time time.Time

type CheckHelper struct{}

func (h *CheckHelper) Check(v string) bool {
	return true
}
