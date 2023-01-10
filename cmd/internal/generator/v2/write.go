package generate

import (
	"fmt"
	"strconv"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

// write the contents of the given generator to g's buffer
func (g *generator) From(r *generator) {
	if g.werr != nil {
		return
	}

	_, g.werr = g.buf.Write(r.buf.Bytes())
}

// write given string to the buffer as is
func (g *generator) S(s string) {
	if g.werr != nil {
		return
	}

	_, g.werr = g.buf.WriteString(s)
}

// write given text and args to buffer using fmt.Sprintf
func (g *generator) F(text string, args ...any) {
	if g.werr != nil {
		return
	}

	_, g.werr = g.buf.WriteString(fmt.Sprintf(text, args...))
}

// interpret the given text with the args, write it to buffer and append a new line
func (g *generator) L(text string, args ...any) {
	if g.werr != nil {
		return
	}
	if len(text) > 0 {
		g.P(text, args...)
	}
	g.werr = g.buf.WriteByte('\n')
}

// replace last line with the provided text and args
func (g *generator) RL(text string, args ...any) {
	if g.werr != nil {
		return
	}

	n := g.buf.Len()
	b := g.buf.Bytes()
	if len(b) > 0 {
		if b[n-1] == '\n' {
			b = b[:n-1]
			n -= 1
		}

		for ; n > 0; n-- {
			if b[n-1] == '\n' {
				break
			}
		}
		g.buf.Truncate(n)
	}

	g.L(text, args...)
}

// interpret the given text with the args and write it to buffer
func (g *generator) P(text string, args ...any) {
	if g.werr != nil {
		return
	}

	out := make([]byte, 0, len(text))
	for i := 0; i < len(text); i++ {
		c := text[i]

		// plain character
		if c != '$' {
			out = append(out, c)
			continue
		}

		// end of text; write literal $
		if len(text) == i+1 {
			out = append(out, '$')
			break
		}

		// when "$$"; escape
		if text[i+1] == '$' {
			out = append(out, '$')
			i += 1
			continue
		}

		if len(out) > 0 {
			_, g.werr = g.buf.Write(out)
			out = out[0:0]
		}

		////////////////////////////////////////////////////////////////

		c2 := text[i+1]

		// when "$<index>"; write args[index]
		if is_digit(c2) {
			j := i + 2
			for ; j < len(text) && is_digit(text[j]); j++ {
			}

			n, err := strconv.Atoi(text[i+1 : j])
			if err != nil {
				g.werr = fmt.Errorf("failed to parse %s: %v", text[i:j], err)
				return
			} else if n >= len(args) {
				g.werr = fmt.Errorf("index %s out of range", text[i:j])
				return
			}

			g.A(args[n])
			i = j - 1
			continue
		}

		// when "${<expansion>}"; write args[index]
		if c2 == '{' {
			j := i + 2
			for ; j < len(text) && text[j] != '}'; j++ {
			}

			exp := text[i+2 : j]
			if err := g.X(exp, args); err != nil {
				g.werr = err
			}
			i = j
			continue
		}
	}

	if g.werr == nil {
		_, g.werr = g.buf.Write(out)
	}
}

type exprOp uint

func (op exprOp) has(ops ...exprOp) bool {
	for _, v := range ops {
		if op&v == 0 {
			return false
		}
	}
	return true
}

const (
	unary_not exprOp = 1 << iota
	bool_and
	bool_or
	func_call
)

// X parses and evaluates the expression x similar to how brace expansions work.
func (g *generator) X(x string, args []any) error {
	var a any
	switch {
	case len(x) > 0 && is_digit(x[0]):
		i := 0
		for ; i < len(x) && is_digit(x[i]); i++ {
		}

		ind, err := strconv.Atoi(x[0:i])
		if err != nil {
			return fmt.Errorf("failed to parse %s: %v", x, err)
		} else if ind >= len(args) {
			return fmt.Errorf("index %s out of range", x)
		}

		a = args[ind]
		x = x[i:]
	}

	var op exprOp
	var mod string
	switch {

	// - "[||]": will expand arg as "expr1 || expr2"
	// - "[&&]": will expand arg as "expr1 && expr2"
	// - "[!|]": will expand arg as "!expr1 || !expr2"
	// - "[!&]": will expand arg as "!expr1 && !expr2"
	// - "[@]": will expand arg as "h(g(f(o)))" (assumes rules are preproc funcs)
	case len(x) > 0 && x[0] == '[':
		i := 1
		for ; i < len(x) && x[i] != ']'; i++ {
		}

		switch x[1:i] {
		case "||":
			op = bool_or
		case "&&":
			op = bool_and
		case "!|":
			op = unary_not | bool_or
		case "!&":
			op = unary_not | bool_and
		case "@":
			op = func_call
		}

	// - ":e": when arg=rule, will generate error expression
	// - ":p": when arg=rule, can group expression
	// - ":g": when arg=obj, will generate object code instead of the default objec identifier
	case len(x) > 0 && x[0] == ':':
		mod = x[1:]
	}

	switch v := a.(type) {
	default:
		g.A(a)

	case rules.List:
		g.gen_rule_list_expr([]*rules.Rule(v), op)

	case []*rules.Rule:
		g.gen_rule_list_expr(v, op)

	case []*rules.Arg:
		g.gen_arg_list_expr(v)

	case *rules.Arg:
		g.genArg(v)

	case *rules.Rule:
		switch mod {
		case "err", "e":
			g.genErrorExpr(v)
		case "p":
			g.gen_rule_expr(v, true)
		default:
			g.gen_rule_expr(v, false)
		}

	case *types.Obj:
		switch mod {
		case "gen", "g":
			g.genObjCode(v)
		default:
			g.S(g.vars[v])
		}

	case types.Ident:
		if pkg := v.GetPkg(); g.file.pkg.Path != pkg.Path {
			pkg := g.file.addImport(pkg)
			g.S(pkg.name + "." + v.GetName())
		} else {
			g.S(v.GetName())
		}
	}

	return nil
}

// write the given arg to the buffer
func (g *generator) A(a any) {
	switch v := a.(type) {
	default:
		g.F("%v", v)

	case types.Ident:
		if pkg := v.GetPkg(); g.file.pkg.Path != pkg.Path {
			pkg := g.file.addImport(pkg)
			g.S(pkg.name + "." + v.GetName())
		} else {
			g.S(v.GetName())
		}

	case *types.Type:
		g.S(v.TypeString(&g.file.pkg))

	case *types.Obj:
		g.S(g.vars[v])

	case *rules.Rule:
		o := g.info.RuleObjMap[v]
		g.genIsRuleExpr(o, v)

	case *rules.Arg:
		g.genArg(v)

	case func():
		v()
	}
}

func is_digit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func is_alnum(c byte) bool {
	return c == '_' ||
		(c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}
