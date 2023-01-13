package generate

import (
	"fmt"
	"strconv"

	"github.com/frk/valid/cmd/internal/types"
)

type scope struct {
	parent   *scope
	children []*scope
	exprs    map[*types.Obj]string
}

type parser struct {
	tokens chan *token
	scopes map[any]*scope
	// the validate method's body
	body *block_stmt
	// the current statement
	cur_stmt stmt
	// the current expression
	cur_expr expr
	// stack of paren_expr
	parens []*paren_expr

	// remove after testing
	debug bool
}

func (p *parser) _debug(mesg string, vv ...any) {
	if p.debug {
		fmt.Printf(mesg, vv...)
		fmt.Println()
	}
}

func (p *parser) init() *parser {
	p.tokens = make(chan *token)
	p.scopes = make(map[any]*scope)
	p.body = new(block_stmt)
	p.cur_stmt = p.body
	return p
}

func (p *parser) parse(in string, args ...any) {
	go (&scanner{in: in, out: p.tokens}).run()
	for tok := range p.tokens {
		if tok.t == t_ws { // ignore ws
			continue
		}

		if p.debug {
			v := any(tok.v)
			if v == "" {
				v = tok.t
			}
			fmt.Printf("- %q\n", v)
		}

		switch tok.t {
		case t_eof:
			return
		case t_ws: // ignore
		case t_comment:
			p.add_line_comment(tok)
		case t_int, t_float:
			p.add_lit_expr(tok)
		case t_ident:
			p.add_ident_expr(tok)
		case t_param, t_paramx:
			p.add_param_expr(tok, args)
		case t_inc, t_dec:
			p.add_incdec_stmt(tok)
		case t_ptr, t_amp, t_not:
			p.add_unary_expr(tok)
		case t_land, t_lor, t_eql, t_lss, t_gtr, t_neq, t_leq, t_geq:
			p.add_binary_expr(tok)
		case t_assign, t_define:
			p.add_assign_stmt(tok)
		case t_lparen:
			p.add_paren_expr()
		case t_rparen:
			p.pop_paren_expr()
		case t_lbrack:
			p.add_index_expr()
		case t_rbrack:
			p.pop_index_expr()
		case t_lbrace:
			p.add_block_stmt()
		case t_rbrace:
			p.pop_block_stmt()
		case t_comma:
			p.handle_comma()
		case t_period:
		case t_semicolon:
			p.handle_semicolon()
		case t_if:
			p.add_if_stmt()
		case t_else:
			p.add_else_stmt()
		case t_for:
			p.add_for_stmt()
		case t_range:
			p.add_range_stmt()
		case t_return:
			p.add_return_stmt()
		case t_package:
		case t_import:
		case t_func:
			p.add_func_decl()
		}
	}
}

func (p *parser) fail(mesg string, vv ...any) {
	panic(fmt.Sprintf("bad token stream: "+mesg, vv...))
}

func (p *parser) end() {
	if p.cur_expr != nil {
		s := &expr_stmt{p.cur_expr}
		p.body.list = append(p.body.list, s)
		p.cur_expr = nil
	}
}

func (p *parser) handle_comma() {
	if p.cur_expr == nil {
		p.fail("unexpected ',' in %T", p.cur_stmt)
	}

	switch cx := p.cur_expr.(type) {
	default:
		x := &expr_list{}
		x.items = append(x.items, cx)
		p.cur_expr = x

	case *expr_list:
		// nothing to do

	}
}

func (p *parser) handle_semicolon() {
	s := &expr_stmt{p.cur_expr}

	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected ';' in %T", p.cur_stmt)

	case *block_stmt:
		cs.list = append(cs.list, s)

	case *for_stmt:
		cs.cond = s.x

	case *assign_stmt:
		cs.rhs = append(cs.rhs, s.x)
		p.cur_stmt = cs.outer
	}

	p.cur_expr = nil
}

func (p *parser) add_line_comment(tok *token) {
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected '%s' in %T", tok.t, p.cur_stmt)

	case *block_stmt:
		cs.list = append(cs.list, &line_comment{tok.v})

	}
}

func (p *parser) add_block_stmt() {
	block := new(block_stmt)
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected '{' in %T", p.cur_stmt)

	case *block_stmt:
		block.outer = cs
		cs.list = append(cs.list, block)

	case *if_stmt:
		switch {
		case cs.cond == nil && p.cur_expr == nil:
			p.fail("missing conditional in 'if'")

		case cs.cond == nil && p.cur_expr != nil:
			cs.cond = p.cur_expr
			p.cur_expr = nil

		}
		block.outer = cs
		cs.body = block

	case *for_stmt:
		switch {
		case p.cur_expr != nil && cs.cond == nil:
			cs.cond = p.cur_expr
			p.cur_expr = nil

		case p.cur_expr != nil && cs.post == nil:
			cs.post = &expr_stmt{p.cur_expr}
			p.cur_expr = nil
		}
		block.outer = cs
		cs.body = block

	case *range_stmt:
		switch {
		case cs.x == nil && p.cur_expr == nil:
			p.fail("missing expression in 'range'")

		case cs.x == nil && p.cur_expr != nil:
			cs.x = p.cur_expr
			p.cur_expr = nil
		}
		block.outer = cs
		cs.body = block

	}
	p.cur_stmt = block
}

func (p *parser) pop_block_stmt() {
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected '}' in %T", p.cur_stmt)

	case *block_stmt:
		if p.cur_expr != nil {
			cs.list = append(cs.list, &expr_stmt{p.cur_expr})
			p.cur_expr = nil
		}
		p.cur_stmt = cs.outer
	}
}

func (p *parser) add_if_stmt() {
	s := &if_stmt{outer: p.cur_stmt}
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected 'if' in %T", p.cur_stmt)

	case *if_stmt:
		cs.els = s

	case *block_stmt:
		cs.list = append(cs.list, s)
	}
	p.cur_stmt = s
}

func (p *parser) add_else_stmt() {
	if _, ok := p.cur_stmt.(*if_stmt); !ok {
		p.fail("unexpected 'else' in %T", p.cur_stmt)
	}

	for tok := range p.tokens {
		switch tok.t {
		default:
			p.fail("unexpected token %s after 'else'", tok.t)
		case t_ws: // ignore
		case t_if:
			p.add_if_stmt()
		case t_lparen:
			p.add_block_stmt()
		}
	}
}

func (p *parser) add_for_stmt() {
	s := &for_stmt{}
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected 'for' in %T", p.cur_stmt)

	case *block_stmt:
		s.outer = cs
		cs.list = append(cs.list, s)
	}
	p.cur_stmt = s
}

func (p *parser) add_range_stmt() {
	s := &range_stmt{}
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected 'range' in %T", cs)

	case *for_stmt:
		// replace for-stmt with range-stmt in outer block
		s.outer = cs.outer
		s.outer.list[len(cs.outer.list)-1] = s

	case *assign_stmt:
		// replace for-stmt with range-stmt in outer block
		fs, ok := cs.outer.(*for_stmt)
		if !ok {
			p.fail("unexpected 'range' assignment in %T", cs)
		}
		s.outer = fs.outer
		s.outer.list[len(fs.outer.list)-1] = s

		if len(cs.lhs) > 0 {
			s.key = cs.lhs[0]
		}
		if len(cs.lhs) > 1 {
			s.val = cs.lhs[1]
		}
		s.tt = cs.tt
	}
	p.cur_stmt = s
}

func (p *parser) add_return_stmt() {
	ret := new(return_stmt)
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected 'return' in %T", p.cur_stmt)

	case *block_stmt:
		ret.outer = cs
		cs.list = append(cs.list, ret)
	}
	p.cur_stmt = ret
}

func (p *parser) add_assign_stmt(t *token) {
	asn := &assign_stmt{outer: p.cur_stmt, tt: t.t}
	switch cs := p.cur_stmt.(type) {
	default:
		p.fail("unexpected '%s' in %T", t.t, p.cur_stmt)

	case *block_stmt:
		cs.list = append(cs.list, asn)

	case *if_stmt:
		cs.init = asn

	case *for_stmt:
		if cs.init == nil {
			cs.init = asn
		} else {
			cs.post = asn
		}
	}

	switch cx := p.cur_expr.(type) {
	default:
		asn.lhs = append(asn.lhs, cx)
	case *expr_list:
		asn.lhs = cx.items
	}

	p.cur_stmt = asn
	p.cur_expr = nil
}

func (p *parser) add_incdec_stmt(t *token) {
	switch cx := p.cur_expr.(type) {
	case nil:
		p.fail("missing expression for '%s' in %T", t.t, p.cur_stmt)

	default:
		s := &incdec_stmt{outer: p.cur_stmt, x: cx, tt: t.t}
		switch cs := p.cur_stmt.(type) {
		default:
			p.fail("unsupported '%s' in %T", t.t, p.cur_stmt)

		case *for_stmt:
			cs.post = s

		case *block_stmt:
			cs.list = append(cs.list, s)
		}
		p.cur_expr = nil
	}
}

func (p *parser) add_func_decl() {
	// ...
}

func (p *parser) add_paren_expr() {
	x := &paren_expr{}
	p.parens = append(p.parens, x)

	switch cx := p.cur_expr.(type) {
	default:
		p.fail("unexpected '(' in %T", p.cur_stmt)

	case nil: // nothing to do

	case *paren_expr:
		cx.x = x

	case *binary_expr:
		cx.right = x
		x.outer = cx

	case *expr_list:
		cx.items = append(cx.items, x)
		x.outer = cx
	}

	p.cur_expr = x
}

func (p *parser) pop_paren_expr() {
	l := len(p.parens)
	x := p.parens[l-1]
	p.parens = p.parens[:l-1]

	if x != p.cur_expr {
		x.x = p.cur_expr
		p.cur_expr = x
	}
	if x.outer != nil {
		p.cur_expr = x.outer
	}
}

func (p *parser) add_index_expr() {
	switch cx := p.cur_expr.(type) {
	default:
		p.fail("unexpected '%s' in %T", t_lbrack, p.cur_stmt)

	case *ident_expr, *param_expr, *index_expr, *call_expr, *selector_expr:
		p.cur_expr = &index_expr{outer: p.cur_stmt, x: cx}
	}
}

func (p *parser) pop_index_expr() {
	// ...
}

func (p *parser) add_lit_expr(t *token) {
	p.add_expr(&lit_expr{outer: p.cur_stmt, tt: t.t, value: t.v})
}

func (p *parser) add_ident_expr(t *token) {
	p.add_expr(&ident_expr{outer: p.cur_stmt, name: t.v})
}

func (p *parser) add_param_expr(t *token, args []any) {
	a, _ := p.parse_param(t, args)
	p.add_expr(&param_expr{outer: p.cur_stmt, arg: a})
}

func (p *parser) add_unary_expr(t *token) {
	p.add_expr(&unary_expr{outer: p.cur_stmt, op: t.t})
}

func (p *parser) add_binary_expr(t *token) {
	left := p.cur_expr
	if left == nil {
		p.fail("unexpected %s in %T", t, p.cur_stmt)
	}
	p.cur_expr = &binary_expr{op: t.t, left: left}
}

func (p *parser) add_expr(x expr) {
	switch cx := p.cur_expr.(type) {
	default:
		p.fail("unexpected %T in %T", x, p.cur_stmt)
	case nil:
		p.cur_expr = x
	case *unary_expr:
		cx.x = x
		p.cur_expr = x
	case *binary_expr:
		cx.right = x
	case *index_expr:
		cx.index = x
	case *call_expr:
		cx.args = append(cx.args, x)
	case *selector_expr:
		cx.sel = x
	case *paren_expr:
		cx.x = x
		p.cur_expr = x
	case *expr_list:
		cx.items = append(cx.items, x)
	}
}

func (p *parser) parse_param(t *token, args []any) (a any, x string) {
	i := 0
	for ; i < len(t.v) && is_digit(t.v[i]); i++ {
	}

	ind, err := strconv.Atoi(t.v[0:i])
	if err != nil {
		p.fail("failed to parse %q: %v", t.v, err)
	} else if ind >= len(args) {
		p.fail("index %q out of range", t.v)
	}

	a = args[ind]
	x = x[i:]
	return a, x
}
