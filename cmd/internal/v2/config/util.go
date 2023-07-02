package config

import (
	"gopkg.in/yaml.v3"
)

func (v RuleErrMesg) ToYAML() ([]byte, error) {
	return yaml.Marshal(v)
}
