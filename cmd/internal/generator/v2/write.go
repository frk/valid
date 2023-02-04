package generate

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func node_string(n node) (string, error) {
	buf := bytes.Buffer{}
	if err := write_node(n, &buf); err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func write_node(n node, w io.Writer) error {
	out := &writer{w: w}
	n.writeTo(out)
	return out.err
}

type writer struct {
	w   io.Writer
	err error
}

func (w *writer) write(s string) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write([]byte(s))
}

func (d *method_decl) writeTo(w *writer) {
	w.write("func (")
	if d.recv != nil {
		d.recv.writeTo(w)
		w.write(" ")
	}
	d.typ.writeTo(w)
	w.write(") ")
	w.write(d.name)
	w.write("(")
	d.params.writeTo(w)
	w.write(")")
	if len(d.results) > 0 {
		w.write("(")
		d.results.writeTo(w)
		w.write(")")
	}
	d.body.writeTo(w)
}

func (s *block_stmt) writeTo(w *writer) {
	w.write("{")
	for _, s := range s.list {
		w.write("\n")
		s.writeTo(w)
	}
	if len(s.list) > 0 {
		w.write("\n")
	}
	w.write("}")
}

func (s *if_stmt) writeTo(w *writer) {
	w.write("if ")
	if s.init != nil {
		s.init.writeTo(w)
		w.write("; ")
	}

	s.cond.writeTo(w)
	w.write(" ")
	s.body.writeTo(w)

	if s.els != nil {
		w.write(" else ")
		s.els.writeTo(w)
	}
}

func (s *for_stmt) writeTo(w *writer) {
	w.write("for ")
	switch {
	case s.init == nil && s.cond == nil && s.post == nil:
		s.body.writeTo(w)
	case s.init == nil && s.cond != nil && s.post == nil:
		s.cond.writeTo(w)
		w.write(" ")
		s.body.writeTo(w)
	default:
		if s.init != nil {
			s.init.writeTo(w)
		}
		w.write(";")
		if s.cond != nil {
			w.write(" ")
			s.cond.writeTo(w)
		}
		w.write(";")
		if s.post != nil {
			w.write(" ")
			s.post.writeTo(w)
		}
		w.write(" ")
		s.body.writeTo(w)
	}
}

func (s *range_stmt) writeTo(w *writer) {
	w.write("for ")
	if s.key != nil {
		s.key.writeTo(w)
	} else if s.val != nil {
		w.write("_")
	}

	if s.val != nil {
		w.write(", ")
		s.val.writeTo(w)
	}

	if s.key != nil || s.val != nil {
		w.write(" ")
		w.write(s.tt.String())
		w.write(" ")
	}

	w.write("range ")
	s.x.writeTo(w)
	w.write(" ")
	s.body.writeTo(w)
}

func (s *return_stmt) writeTo(w *writer) {
	w.write("return")
	if s.res != nil {
		w.write(" ")
		(&expr_list{s.res}).writeTo(w)
	}
}

func (s *assign_stmt) writeTo(w *writer) {
	(&expr_list{s.lhs}).writeTo(w)
	w.write(" ")
	w.write(s.tt.String())
	w.write(" ")
	(&expr_list{s.rhs}).writeTo(w)
}

func (s *incdec_stmt) writeTo(w *writer) {
	s.x.writeTo(w)
	w.write(s.tt.String())
}

func (s *expr_stmt) writeTo(w *writer) {
	s.x.writeTo(w)
}

func (c *line_comment) writeTo(w *writer) {
	w.write("//")
	w.write(c.text)
}

func (g *line_comment_group) writeTo(w *writer) {
	if len(g.list) == 0 {
		return
	}

	w.write("//")
	w.write(g.list[0])
	for _, c := range g.list[1:] {
		w.write("\n//")
		w.write(c)
	}
}

func (x *lit_expr) writeTo(w *writer) {
	w.write(x.value)
}

func (x *ident_expr) writeTo(w *writer) {
	w.write(x.name)
}

func (x *unary_expr) writeTo(w *writer) {
	w.write(x.op.String())
	x.x.writeTo(w)
}

func (x *binary_expr) writeTo(w *writer) {
	x.left.writeTo(w)
	w.write(" ")
	w.write(x.op.String())
	w.write(" ")
	x.right.writeTo(w)
}

func (x *index_expr) writeTo(w *writer) {
	x.x.writeTo(w)
	w.write("[")
	x.index.writeTo(w)
	w.write("]")
}

func (x *call_expr) writeTo(w *writer) {
	x.fun.writeTo(w)
	w.write("(")
	(&expr_list{x.args}).writeTo(w)
	if x.ellipsis {
		w.write("...")
	}
	w.write(")")
}

func (x *selector_expr) writeTo(w *writer) {
	x.x.writeTo(w)
	w.write(".")
	x.sel.writeTo(w)
}

func (x *paren_expr) writeTo(w *writer) {
	w.write("(")
	if x.x != nil {
		x.x.writeTo(w)
	}
	w.write(")")
}

func (x *expr_list) writeTo(w *writer) {
	if len(x.items) == 0 {
		return
	}
	x.items[0].writeTo(w)
	for _, x := range x.items[1:] {
		w.write(", ")
		x.writeTo(w)
	}
}

func (ls field_list) writeTo(w *writer) {
	if len(ls) == 1 {
		return
	}

	ls[0].writeTo(w)
	for _, f := range ls[1:] {
		w.write(", ")
		f.writeTo(w)
	}
}

func (f *field_item) writeTo(w *writer) {
	if len(f.names) > 0 {
		f.names[0].writeTo(w)
		for _, n := range f.names[1:] {
			w.write(", ")
			n.writeTo(w)
		}
		w.write(" ")
	}
	if f.variadic {
		w.write("...")
	}
	f.typ.writeTo(w)
}

func (x *param_expr) writeTo(w *writer) {
	// ...
}

////////////////////////////////////////////////////////////////////////////////
//
//
////////////////////////////////////////////////////////////////////////////////

func (g *generator) _write(b []byte) {
	if g.werr != nil {
		return
	}
	_, g.werr = g.buf.Write(b)
}

func (g *generator) _writeb(b byte) {
	if g.werr != nil {
		return
	}
	g.werr = g.buf.WriteByte(b)
}

func (g *generator) _writes(s string) {
	if g.werr != nil {
		return
	}
	_, g.werr = g.buf.WriteString(s)
}

// write the contents of the given generator to g's buffer
func (g *generator) From(r *generator) {
	g._write(r.buf.Bytes())
}

// write given string to the buffer as is
func (g *generator) S(s string) {
	g._writes(s)
}

// write given string to the buffer quoted
func (g *generator) Q(s string) {
	g._writes(strconv.Quote(s))
}

// write given text and args to buffer using fmt.Sprintf
func (g *generator) F(text string, args ...any) {
	g._writes(fmt.Sprintf(text, args...))
}

// interpret the given text with the args, write it to buffer and append a new line
func (g *generator) L(text string, args ...any) {
	if len(text) > 0 {
		g.P(text, args...)
	}
	g._writeb('\n')
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

	// - ":e": when a=Rule, will generate error expression
	// - ":p": when a=Rule, can group expression
	// - ":g": when a=Obj, will generate object code instead of the default objec identifier
	// - ":fk": when a=Obj, will generate the object's field key
	// - ":any": when a=[]Arg, will generate the args expressions as when assigned to any
	// - ":any": when a=Arg, will generate the arg expression as when assigned to any
	case len(x) > 0 && x[0] == ':':
		mod = x[1:]
	}

	switch v := a.(type) {
	default:
		g.A(a)
	case string:
		g.Q(v)

	case rules.List:
		g.gen_rule_list_expr([]*rules.Rule(v), op)

	case []*rules.Rule:
		g.gen_rule_list_expr(v, op)

	case []*rules.Arg:
		switch mod {
		case "any":
			g.gen_arg_list_expr(v, true)
		default:
			g.gen_arg_list_expr(v, false)
		}

	case *rules.Arg:
		switch mod {
		case "any":
			g.gen_arg_expr(v, true)
		default:
			g.gen_arg_expr(v, false)
		}

	case *rules.Rule:
		switch mod {
		case "err", "e":
			g.gen_error_expr(v)
		case "p":
			g.gen_rule_expr(v, true)
		default:
			g.gen_rule_expr(v, false)
		}

	case *types.Obj:
		switch mod {
		case "gen", "g":
			g.genObjCode(v)
		case "fk":
			g.gen_obj_field_key(v)
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
		g.gen_arg_expr(v, false)

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
