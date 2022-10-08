package testdata

import (
	"time"
)

type Validator struct {
	F1 string    `is:"required"`
	F2 float64   `is:"required"`
	F3 any       `is:"required"`
	F4 time.Time `is:"required"`

	P1 *string    `is:"required"`
	P2 **float64  `is:"required"`
	P3 *any       `is:"required"`
	P4 *time.Time `is:"required"`
}
