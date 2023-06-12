package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements both the flag.Value and the yaml.Unmarshaler interfaces.
type String struct {
	Value string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (s String) Get() interface{} {
	return s.Value
}

// String implements the flag.Value interface.
func (s String) String() string {
	return s.Value
}

// Set implements the flag.Value interface.
func (s *String) Set(value string) error {
	s.Value = value
	s.IsSet = true
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s *String) UnmarshalYAML(n *yaml.Node) error {
	if !s.IsSet {
		if n.Tag == "!!nil" {
			return nil
		}

		var value string
		if err := n.Decode(&value); err != nil {
			return &Error{C: ERR_YAML_ERROR, tt: s, node: n, err: err}
		}
		s.Value = value
		s.IsSet = true
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Bool implements both the flag.Value and the yaml.Unmarshaler interfaces.
type Bool struct {
	Value bool
	IsSet bool
}

// IsBoolFlag indicates that the Bool type can be used as a boolean flag.
func (b Bool) IsBoolFlag() bool {
	return true
}

// Get implements the flag.Getter interface.
func (b Bool) Get() interface{} {
	return b.String()
}

// String implements the flag.Value interface.
func (b Bool) String() string {
	return strconv.FormatBool(b.Value)
}

// Set implements the flag.Value interface.
func (b *Bool) Set(value string) error {
	if len(value) > 0 {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		b.Value = v
		b.IsSet = true
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (b *Bool) UnmarshalYAML(n *yaml.Node) error {
	if !b.IsSet {
		if n.Tag == "!!nil" {
			return nil
		}

		var value bool
		if err := n.Decode(&value); err != nil {
			return &Error{C: ERR_YAML_ERROR, tt: b, node: n, err: err}
		}
		b.Value = value
		b.IsSet = true
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// StringSlice implements both the flag.Value and the yaml.Unmarshaler interfaces.
type StringSlice struct {
	Value []string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (ss StringSlice) Get() interface{} {
	return ss.String()
}

// String implements the flag.Value interface.
func (ss StringSlice) String() string {
	return strings.Join(ss.Value, ",")
}

// Set implements the flag.Value interface.
func (ss *StringSlice) Set(value string) error {
	if len(value) > 0 {
		ss.Value = append(ss.Value, value)
		ss.IsSet = true
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (ss *StringSlice) UnmarshalYAML(n *yaml.Node) error {
	if !ss.IsSet {
		if n.Tag == "!!nil" {
			return nil
		}

		switch n.Kind {
		case yaml.ScalarNode:
			var value string
			if err := n.Decode(&value); err != nil {
				return &Error{C: ERR_YAML_ERROR, tt: ss, node: n, err: err}
			}
			if len(value) > 0 {
				ss.Value = []string{value}
				ss.IsSet = true
			}
		case yaml.SequenceNode:
			var value []string
			if err := n.Decode(&value); err != nil {
				return &Error{C: ERR_YAML_ERROR, tt: ss, node: n, err: err}
			}
			if len(value) > 0 {
				ss.Value = value
				ss.IsSet = true
			}
		default:
			return &Error{C: ERR_YAML_TYPE, tt: ss, node: n}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// ObjectIdent represents a Go type or function identifier. It implements
// both the flag.Value and the yaml.Unmarshaler interfaces.
type ObjectIdent struct {
	Pkg   string // the package path of the Go object
	Name  string // the name of the Go object
	IsSet bool
}

// Get implements the flag.Getter interface.
func (id ObjectIdent) Get() interface{} {
	return id.String()
}

// String implements the flag.Value interface.
func (id ObjectIdent) String() string {
	v := id.Pkg
	if len(id.Name) > 0 {
		v += "." + id.Name
	}
	return v
}

var rxObjId = regexp.MustCompile(`^\w[\w\.-/]*\w\.[A-Za-z_]\w*$`)

// Set implements the flag.Value interface.
func (id *ObjectIdent) Set(value string) error {
	if len(value) > 0 {
		if !rxObjId.MatchString(value) {
			return fmt.Errorf("config.ObjectIdent.Set: %q invalid object_identifier", value)
		}

		i := strings.LastIndex(value, ".")
		id.Pkg = value[:i]
		id.Name = value[i+1:]
		id.IsSet = true
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (id *ObjectIdent) UnmarshalYAML(n *yaml.Node) error {
	if !id.IsSet {
		if n.Tag == "!!nil" {
			return nil
		}

		var value string
		if err := n.Decode(&value); err != nil {
			return &Error{C: ERR_YAML_ERROR, tt: id, node: n, err: err}
		}
		if err := id.Set(value); err != nil {
			return &Error{C: ERR_OBJECT_IDENT, tt: id, val: value, node: n}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Scalar represents a scalar-only value.
type Scalar struct {
	// The type of the scalar value.
	Type ScalarType `yaml:"-"`
	// The literal string representation of the value.
	Value string `yaml:"-"`
}

type ScalarType uint

const (
	_ ScalarType = iota
	NIL
	BOOL
	STRING
	FLOAT
	INT
)

func (s *Scalar) UnmarshalYAML(n *yaml.Node) error {
	if n.Tag == "!!nil" {
		s.Type = NIL
		return nil
	}

	var val interface{}
	if err := n.Decode(&val); err != nil {
		return err
	}

	switch v := val.(type) {
	case bool:
		s.Type = BOOL
		s.Value = strconv.FormatBool(v)
	case string:
		s.Type = STRING
		s.Value = v
	case float32, float64:
		s.Type = FLOAT
		s.Value = n.Value
	case int, int8, int16, int32, int64:
		s.Type = INT
		s.Value = n.Value
	default:
		return &Error{C: ERR_YAML_TYPE, tt: s, node: n}
	}
	return nil
}

func (s Scalar) MarshalYAML() (interface{}, error) {
	switch s.Type {
	case NIL:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!nil"}, nil
	case BOOL:
		b, err := strconv.ParseBool(s.Value)
		return b, err
	case INT:
		i, err := strconv.ParseInt(s.Value, 10, 64)
		return i, err
	case FLOAT:
		f, err := strconv.ParseFloat(s.Value, 64)
		return f, err
	case STRING:
		return s.Value, nil
	}
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////

// JoinOp represents a logical join operator.
type JoinOp uint

const (
	_ JoinOp = iota
	JOIN_NOT
	JOIN_AND
	JOIN_OR
)

func (op *JoinOp) UnmarshalYAML(n *yaml.Node) error {
	if n.Tag == "!!nil" {
		return nil
	}

	var s string
	if err := n.Decode(&s); err != nil {
		return &Error{C: ERR_YAML_TYPE, tt: op, node: n} //, err: err}
	}

	switch v := strings.ToUpper(s); v {
	case "NOT":
		*op = JOIN_NOT
	case "AND":
		*op = JOIN_AND
	case "OR":
		*op = JOIN_OR
	default:
		return &Error{C: ERR_JOIN_OP, tt: op, val: s, node: n}
	}
	return nil
}

func (op JoinOp) String() string {
	switch op {
	case JOIN_NOT:
		return "NOT"
	case JOIN_AND:
		return "AND"
	case JOIN_OR:
		return "OR"
	}
	return "<invalid>"
}
