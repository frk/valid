package rules

import (
	"fmt"
	"go/token"
	"strconv"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/errors"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"

	"github.com/frk/tagutil"
)

var _ = fmt.Println

// Info holds the result type information for rule-checking a Validator.
type Info struct {
	// The Validator being rule-checked.
	Validator *gotype.Validator
	// The Node representation of the Validator
	RootNode *Node
	// KeyMap maps field keys to the corresponding type nodes.
	KeyMap map[string]*FieldNode
	// EnumMap maps types to a slice of
	// constants declared with that type.
	EnumMap map[*gotype.Type][]gotype.Const
}

// Checker maintains the state of the rule checker.
type Checker struct {
	*Info
	// The set of go/ast packages associated with the type being checked.
	ast *search.AST
	// The package for which the code will be generated.
	pkg gotype.Pkg
	// The file set with which the type being checked is associated.
	fs *token.FileSet
	// The analyzer used to analyze the matched type.
	an *gotype.Analyzer
	// The Validtor struct being rule-checked.
	vs *gotype.Validator
	// The function used for generating unique field keys.
	fieldKey FieldKeyFunc
}

// NewChecker returns a new Checker instance.
// The optional info argument, will be populated during rule-checking.
func NewChecker(ast *search.AST, pkg search.Pkg, fkCfg *config.FieldKeyConfig, info *Info) (c *Checker) {
	c = &Checker{ast: ast, pkg: gotype.Pkg(pkg)}
	c.fieldKey = fkFunc(fkCfg)

	c.Info = info
	if c.Info == nil {
		c.Info = new(Info)
	}
	c.Info.KeyMap = make(map[string]*FieldNode)
	c.Info.EnumMap = make(map[*gotype.Type][]gotype.Const)
	return c
}

// Check rule-checks the validator struct represented by the
// given *search.Match and returns the first error encountered.
func (c *Checker) Check(match *search.Match) error {
	// 1. analyze
	c.an = gotype.NewAnalyzer(match.Named.Obj().Pkg())
	c.vs = c.an.Validator(match.Named)

	// 2. convert to a Node tree
	rootNode, err := c.makeNode(c.vs.Type, nil, nil, nil)
	if err != nil {
		return c.err(err, errOpts{a: c.ast})
	}

	// 3. rule-check the Node
	if err := c.check(rootNode); err != nil {
		return c.err(err, errOpts{a: c.ast})
	}

	// 4. populate c.Info (if no error)
	c.Info.Validator = c.vs
	c.Info.RootNode = rootNode
	return nil
}

func (c *Checker) check(n *Node) error {
	c.fixImplicitRules(n)

	if len(n.IsRules) > 0 {
		if err := c.checkRules(n); err != nil {
			return err
		}
	}
	if len(n.PreRules) > 0 {
		if err := c.checkPreproc(n); err != nil {
			return err
		}
	}

	switch n.Type.Kind {
	case gotype.K_PTR:
		return c.check(n.Elem)
	case gotype.K_ARRAY, gotype.K_SLICE:
		return c.check(n.Elem)
	case gotype.K_MAP:
		if err := c.check(n.Key); err != nil {
			return err
		}
		return c.check(n.Elem)
	case gotype.K_STRUCT:
		for _, f := range n.Fields {
			if err := c.check(f.Type); err != nil {
				return c.err(err, errOpts{sf: f.Field})
			}
		}
	}
	return nil
}

func (c *Checker) checkRules(n *Node) error {
	for _, r := range n.IsRules {
		// Ensure that the Value of a Arg of type ARG_FIELD
		// references a valid field key which will be indicated
		// by a presence of a selector in the analyzer's KeyMap.
		for _, a := range r.Args {
			if a.Type == ARG_FIELD {
				if _, ok := c.Info.KeyMap[a.Value]; !ok {
					return errors.TODO("checkRules: field key matches no known field")
				}
			}
		}

		c.fixRuleArgs(r)

		// Check that the number of arguments provided
		// to the rule is allowed by the spec.
		if r.Spec.ArgMin > -1 && r.Spec.ArgMin > len(r.Args) {
			return errors.TODO("checkRules: rule has not enough arguments")
		}
		if r.Spec.ArgMax > -1 && r.Spec.ArgMax < len(r.Args) {
			return errors.TODO("checkRules: rule has too many arguments")
		}

		// run type specific rule-check
		switch r.Spec.Kind {
		case REQUIRED:
			if err := c.requiredCheck(n, r); err != nil {
				return err
			}
		case COMPARABLE:
			if err := c.comparableCheck(n, r); err != nil {
				return err
			}
		case ORDERED:
			if err := c.orderedCheck(n, r); err != nil {
				return err
			}
		case LENGTH:
			if err := c.lengthCheck(n, r); err != nil {
				return err
			}
		case RANGE:
			if err := c.rangeCheck(n, r); err != nil {
				return err
			}
		case ENUM:
			if err := c.enumCheck(n, r); err != nil {
				return err
			}
		case FUNCTION:
			if err := c.functionCheck(n, r); err != nil {
				return err
			}
		case METHOD:
			if err := c.methodCheck(n, r); err != nil {
				return err
			}
		case OPTIONAL:
			if err := c.optionalCheck(n, r); err != nil {
				return err
			}
		case REMOVE:
			// if err := c.checkRemove(n, r, spec); err != nil {
			// 	return err
			// }
		}
	}
	return nil
}

func (c *Checker) checkPreproc(n *Node) error {
	for _, r := range n.PreRules {
		if r.Spec.Kind != PREPROC {
			return errors.TODO("checkPreproc: rule kind is not preprocessor")
		}

		// Ensure that the Value of a Arg of kind AFIELD
		// references a valid field key which will be indicated
		// by a presence of a selector in the analyzer's KeyMap.
		for _, arg := range r.Args {
			if arg.Type == ARG_FIELD {
				if _, ok := c.Info.KeyMap[arg.Value]; !ok {
					return errors.TODO("checkPreproc: field key matches no known field")
				}
			}
		}

		c.fixRuleArgs(r)

		// Check that the number of arguments provided
		// to the rule is allowed by the spec.
		if r.Spec.ArgMin > -1 && r.Spec.ArgMin > len(r.Args) {
			return errors.TODO("checkPreproc: rule has not enough arguments")
		}
		if r.Spec.ArgMax > -1 && r.Spec.ArgMax < len(r.Args) {
			return errors.TODO("checkPreproc: rule has too many arguments")
		}

		// run spec specific rule-check
		return c.preprocessorCheck(n, r)
	}
	return nil
}

type errOpts Error

func (c *Checker) err(err error, opts errOpts) error {
	e, ok := err.(*Error)
	if !ok {
		return err
	}

	if opts.C > 0 {
		e.C = opts.C
	}
	if opts.a != nil {
		e.a = opts.a
	}
	if opts.c != nil {
		e.c = opts.c
	}
	if opts.rc != nil {
		e.rc = opts.rc
	}
	if opts.rs != nil {
		e.rs = opts.rs
	}
	if opts.rca != nil {
		e.rca = opts.rca
	}
	if opts.rcai != nil {
		e.rcai = opts.rcai
	}
	if opts.rcak != nil {
		e.rcak = opts.rcak
	}
	if opts.ft != nil {
		e.ft = opts.ft
	}
	if e.sf == nil && opts.sf != nil {
		e.sf = opts.sf
	}
	if opts.ty != nil {
		e.ty = opts.ty
	}
	if opts.tag != nil {
		e.tag = opts.tag
	}
	if opts.r != nil {
		e.r = opts.r
	}
	if opts.ra != nil {
		e.ra = opts.ra
	}
	if opts.fp != nil {
		e.fp = opts.fp
	}
	if opts.fpi != nil {
		e.fpi = opts.fpi
	}
	if opts.err != nil {
		e.err = opts.err
	}

	if e.sf != nil {
		e.sfv = e.sf.Var
	}
	if e.ra != nil && e.ra.Type == ARG_FIELD {
		e.raf = c.KeyMap[e.ra.Value].Field
	}
	return e
}

// FieldKeyFunc is the type of the function called by the Checker
// for each field to generate a unique key from the FieldSelector.
type FieldKeyFunc func(gotype.FieldSelector) (key string)

// fkFunc returns a function that, based on the given configuration,
// generates a unique field key for a given FieldSelector.
func fkFunc(c *config.FieldKeyConfig) FieldKeyFunc {
	// keyset & unique are used by the returned function
	// to ensure that the generated key is unique.
	keyset := make(map[string]uint)
	unique := func(key string) string {
		if num, ok := keyset[key]; ok {
			keyset[key] = num + 1
			key += "-" + strconv.FormatUint(uint64(num), 10)
		} else {
			keyset[key] = 1
		}
		return key
	}

	if c != nil && len(c.Tag.Value) > 0 {
		if c.Join.Value {
			// Returns the joined tag values of the fields in the given slice.
			// If one of the fields does not have a tag value set, their name
			// will be used in the join as default.
			return func(fs gotype.FieldSelector) (key string) {
				tag := c.Tag.Value
				sep := c.Separator.Value

				for _, f := range fs {
					t := tagutil.New(f.Tag)
					if t.Contains("is", "omitkey") || f.IsEmbedded {
						continue
					}

					v := t.First(tag)
					if len(v) == 0 {
						v = f.Name
					}
					key += v + sep
				}
				if len(sep) > 0 && len(key) > len(sep) {
					return unique(key[:len(key)-len(sep)])
				}
				return unique(key)
			}
		}

		// Returns the tag value of the last field, if no value was
		// set the field's name will be returned instead.
		return func(fs gotype.FieldSelector) string {
			t := tagutil.New(fs[len(fs)-1].Tag)
			if key := t.First(c.Tag.Value); len(key) > 0 {
				return unique(key)
			}
			return unique(fs[len(fs)-1].Name)
		}
	}

	if c != nil && c.Join.Value {
		sep := c.Separator.Value

		// Returns the joined names of the fields in the given slice.
		return func(fs gotype.FieldSelector) (key string) {
			for _, f := range fs {
				t := tagutil.New(f.Tag)
				if t.Contains("is", "omitkey") || f.IsEmbedded {
					continue
				}
				key += f.Name + sep
			}
			if len(sep) > 0 && len(key) > len(sep) {
				return unique(key[:len(key)-len(sep)])
			}
			return unique(key)
		}
	}

	// Returns the name of the last field.
	return func(fs gotype.FieldSelector) string {
		return unique(fs[len(fs)-1].Name)
	}
}

////////////////////////////////////////////////////////////////////////////////
// helpers

// canConvertRuleArg reports whether or not the arg's literal
// value can be converted to the Go type represented by t.
func (c *Checker) canConvertRuleArg(t *gotype.Type, arg *Arg) bool {
	if arg.Type == ARG_FIELD {
		typ := c.KeyMap[arg.Value].Type

		// can use the addr, accept
		if t.PtrOf(typ.Type) {
			return true
		}

		return t.CanAssign(typ.Type) != gotype.ASSIGNMENT_INVALID
	}

	// t is interface{} or string, accept
	if t.IsEmptyInterface() || t.Kind == gotype.K_STRING {
		return true
	}

	// arg is unknown, accept
	if arg.Type == ARG_UNKNOWN {
		return true
	}

	// both are booleans, accept
	if t.Kind == gotype.K_BOOL && arg.Type == ARG_BOOL {
		return true
	}

	// t is float and option is numeric, accept
	if t.Kind.IsFloat() && (arg.Type == ARG_INT || arg.Type == ARG_FLOAT) {
		return true
	}

	// both are integers, accept
	if t.Kind.IsInteger() && arg.Type == ARG_INT {
		return true
	}

	// t is unsigned and option is not negative, accept
	if t.Kind.IsUnsigned() && arg.Type == ARG_INT && arg.Value[0] != '-' {
		return true
	}

	// arg is string & string can be converted to t, accept
	if arg.Type == ARG_STRING && (t.Kind == gotype.K_STRING || (t.Kind == gotype.K_SLICE &&
		t.Elem.Name == "" && (t.Elem.Kind == gotype.K_UINT8 || t.Elem.Kind == gotype.K_INT32))) {
		return true
	}

	return false
}

// Add implicit rules that weren't provided explicitly and remove
// implicit rules that were explicitly specified to be omitted.
func (c *Checker) fixImplicitRules(n *Node) {
	// add
	if n.Type.HasIsValid() && !n.IsRules.Contains("isvalid") {
		r := &Rule{Name: "isvalid"}
		r.Spec = GetSpec(r.Name)
		n.IsRules.Add(r)
	}

	// remove
	if n.IsRules.Contains("-isvalid") {
		n.IsRules.Remove("isvalid")
		n.IsRules.Remove("-isvalid") // not needed anymore
	}
}

// fixRuleArgs updates the rule's Args based on ArgOpts specified by the spec.
func (c *Checker) fixRuleArgs(r *Rule) {
	for i, argOpts := range r.Spec.ArgOpts {
		if len(r.Args) <= i {
			// If no rule arg was provided at the ith index
			// then initialize it to an "unknown" and see if
			// the argOpts contains a default (key="") entry.
			arg := &Arg{Type: ARG_UNKNOWN}
			if opt, ok := argOpts[""]; ok {
				*arg = opt
			}
			r.Args = append(r.Args, arg)
			continue
		}

		arg := r.Args[i]

		// If the the rule's argument is "unknown" and
		// the argOpts contain a default (key="") entry,
		// then update the argument with the default.
		if arg.Value == "" && arg.Type == ARG_UNKNOWN {
			if opt, ok := argOpts[""]; ok {
				*arg = opt
			}
			continue
		}

		if arg.Type != ARG_FIELD {
			// If a Arg's non-field Value matches an entry
			// in the argOpts map, then update the Arg.
			if opt, ok := argOpts[arg.Value]; ok {
				*arg = opt
			}
		}
	}
}
