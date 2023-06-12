package config

import (
	"gopkg.in/yaml.v3"
)

type TypeConfig struct {
	// The identifier of the custom type that the tool should support.
	Type ObjectIdent `yaml:"type"`
	// The accessor(s) to the type's actual value.
	Value TypeValue `yaml:"value"`
	// The name of the type's field or method that should be used when
	// generating code for the "required" check on values of this type.
	//
	// If the name identifes a field that field MUST be of type bool.
	// If it identifes a method that method MUST take zero arguments
	// and return a single value of type bool.
	RequiredCheck String `yaml:"required_check"`
	// The name of the type's field or method that should be used when
	// generating code for the "notnil" check on values of this type.
	//
	// If the name identifes a field that field MUST be of type bool.
	// If it identifes a method that method MUST take zero arguments
	// and return a single value of type bool.
	NotnilCheck String `yaml:"notnil_check"`
	// The name of the type's field or method that should be used when
	// generating code for the "optional" check on values of this type.
	//
	// If the name identifes a field that field MUST be of type bool.
	// If it identifes a method that method MUST take zero arguments
	// and return a single value of type bool.
	OptionalCheck String `yaml:"optional_check"`
	// The name of the type's field or method that should be used when
	// generating code for the "omitnil" check on values of this type.
	//
	// If the name identifes a field that field MUST be of type bool.
	// If it identifes a method that method MUST take zero arguments
	// and return a single value of type bool.
	OmitnilCheck String `yaml:"omitnil_check"`
}

type TypeValue struct {
	// The name of the type's field or method that should be used by
	// the generated code to get the actual value to be validated.
	//
	// If the name identifes a method that method MUST take zero
	// arguments and return a single value.
	Get String `yaml:"get"`
	// The name of the type's field or method that should be used by
	// the generated code to set the actual value.
	//
	// If the name identifes a method that method MUST take
	// a single argument and return zero values.
	Set String `yaml:"set"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (v *TypeValue) UnmarshalYAML(n *yaml.Node) error {
	if n.Tag == "!!nil" {
		return nil
	}

	switch n.Kind {
	case yaml.ScalarNode:
		var s string
		if err := n.Decode(&s); err != nil {
			return &Error{C: ERR_YAML_ERROR, tt: v, node: n, err: err}
		}
		if len(s) > 0 {
			_ = v.Get.Set(s)
			_ = v.Set.Set(s)
		}

	case yaml.MappingNode:
		type V TypeValue
		if err := n.Decode((*V)(v)); err != nil {
			return &Error{C: ERR_YAML_ERROR, tt: v, node: n, err: err}
		}

	default:
		return &Error{C: ERR_YAML_TYPE, tt: v, node: n}
	}
	return nil
}
