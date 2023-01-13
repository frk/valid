package generate

// The supported subset of lexial tokens of the Go
// language, plus the additional parameter tokens.
type token_type uint8

func (t token_type) String() string {
	return _token_types[t]
}

const (
	t_invalid token_type = iota
	t_eof
	t_ws      // ' ', '\t' only, line whitespace is not supported
	t_comment // "//" only

	t_int   // 123
	t_float // 12.3

	t_param  // positional parameter, e.g. $0, $1, ...
	t_paramx // parameter expansion, e.g. ${0:e}

	t_ptr // *
	t_amp // &

	t_land // &&
	t_lor  // ||

	t_eql    // ==
	t_lss    // <
	t_gtr    // >
	t_assign // =
	t_not    // !

	t_neq    // !=
	t_leq    // <=
	t_geq    // >=
	t_define // :=

	t_inc // ++
	t_dec // --

	t_lparen // (
	t_rparen // )
	t_lbrack // [
	t_rbrack // ]
	t_lbrace // {
	t_rbrace // }

	t_comma     // ,
	t_period    // .
	t_semicolon // ;

	t_ident

	// supported keywords
	t_if
	t_else
	t_for
	t_range
	t_return
	t_package
	t_import
	t_func
)

var _token_types = [...]string{
	t_invalid:   "<invalid token>",
	t_eof:       "<eof token>",
	t_ws:        "<whitespace token>",
	t_comment:   "<comment token>",
	t_int:       "<int token>",
	t_float:     "<float token>",
	t_param:     "<param token>",
	t_paramx:    "<paramx token>",
	t_ptr:       "*",
	t_amp:       "&",
	t_land:      "&&",
	t_lor:       "||",
	t_eql:       "==",
	t_lss:       "<",
	t_gtr:       ">",
	t_assign:    "=",
	t_not:       "!",
	t_neq:       "!=",
	t_leq:       "<=",
	t_geq:       ">=",
	t_define:    ":=",
	t_inc:       "++",
	t_dec:       "--",
	t_lparen:    "(",
	t_rparen:    ")",
	t_lbrack:    "[",
	t_rbrack:    "]",
	t_lbrace:    "{",
	t_rbrace:    "}",
	t_comma:     ",",
	t_period:    ".",
	t_semicolon: ";",
	t_ident:     "<identifier token>",
	t_if:        "if",
	t_else:      "else",
	t_for:       "for",
	t_range:     "range",
	t_return:    "return",
	t_package:   "package",
	t_import:    "import",
	t_func:      "func",
}

type token struct {
	t token_type // the type of the token
	v string     // literal value; used only when needed
	p int        // position; used only when needed
}

type scanner struct {
	in  string
	out chan<- *token

	ch  byte // current character (ascii only)
	off int
	eof bool
}

func (s *scanner) run() {
loop:
	for {
		s.next()
		if s.eof {
			s.emit_eof()
			return
		}

		if s.is_letter(s.ch) {
			lit := s.scan_ident()
			switch lit {
			default:
				s.out <- &token{t: t_ident, v: lit}
			case "if":
				s.out <- &token{t: t_if}
			case "else":
				s.out <- &token{t: t_else}
			case "for":
				s.out <- &token{t: t_for}
			case "range":
				s.out <- &token{t: t_range}
			case "return":
				s.out <- &token{t: t_return}
			case "package":
				s.out <- &token{t: t_package}
			case "import":
				s.out <- &token{t: t_import}
			case "func":
				s.out <- &token{t: t_func}
			}
			continue
		}

		if s.is_digit(s.ch) {
			lit, is_fp := s.scan_number()
			switch {
			case lit != "" && !is_fp:
				s.out <- &token{t: t_int, v: lit}
			case lit != "" && is_fp:
				s.out <- &token{t: t_float, v: lit}
			case lit == "":
				s.out <- &token{t: t_invalid, v: string(s.ch), p: s.off}
				return
			}
			continue
		}

		switch s.ch {
		case ' ', '\t':
			s.emit_ws()
		case '$':
			switch c := s.peek(); true {
			case s.is_digit(c): // param?
				s.emit_param()
			case c == '{': // paramx?
				s.emit_paramx()
			default:
				break loop
			}
		case '/':
			if s.peek() == '/' {
				s.emit_comment()
			} else {
				break loop
			}
		case '*':
			s.emit_ptr()
		case '&':
			if s.peek() == '&' {
				s.emit_land()
			} else {
				s.emit_amp()
			}
		case '|':
			if s.peek() == '|' {
				s.emit_lor()
			} else {
				break loop
			}
		case '=':
			if s.peek() == '=' {
				s.emit_eql()
			} else {
				s.emit_assign()
			}
		case '>':
			if s.peek() == '=' {
				s.emit_geq()
			} else {
				s.emit_gtr()
			}
		case '<':
			if s.peek() == '=' {
				s.emit_leq()
			} else {
				s.emit_lss()
			}
		case '!':
			if s.peek() == '=' {
				s.emit_neq()
			} else {
				s.emit_not()
			}
		case ':':
			if s.peek() == '=' {
				s.emit_define()
			} else {
				break loop
			}
		case '+':
			if s.peek() == '+' {
				s.emit_inc()
			} else {
				break loop
			}
		case '-':
			if s.peek() == '-' {
				s.emit_dec()
			} else {
				break loop
			}
		case '(':
			s.emit_lparen()
		case '[':
			s.emit_lbrack()
		case '{':
			s.emit_lbrace()
		case ')':
			s.emit_rparen()
		case ']':
			s.emit_rbrack()
		case '}':
			s.emit_rbrace()
		case ',':
			s.emit_comma()
		case '.':
			s.emit_period()
		case ';':
			s.emit_semicolon()
		}
	}
}

func (s *scanner) next() {
	if s.off < len(s.in) {
		s.ch = s.in[s.off]
		s.off += 1
	} else {
		s.off = len(s.in)
		s.eof = true
	}
}

func (s *scanner) peek() byte {
	if s.off < len(s.in) {
		return s.in[s.off]
	}
	return 0
}

func (s *scanner) scan_ident() string {
	offs := s.off - 1
	for ; s.is_ident_char(s.peek()); s.next() {
	}
	return s.in[offs:s.off]
}

func (s *scanner) scan_number() (_ string, is_fp bool) {
	offs := s.off - 1
	is_digsep := false

loop:
	for {
		switch c := s.peek(); {
		default:
			break loop
		case c == '.' && !is_fp:
			is_fp = true
		case c == '_' && !is_digsep:
			is_digsep = true
		case s.is_digit(c):
			// ok
		}

		s.next()
		is_digsep = (s.ch == '_')
	}

	return s.in[offs:s.off], is_fp
}

func (s *scanner) emit_eof() {
	s.out <- &token{t: t_eof}
}

func (s *scanner) emit_comment() {
	s.out <- &token{t: t_comment, v: s.in[s.off+1:]}
	s.off = len(s.in)
}

func (s *scanner) emit_ws() {
	for ; s.is_ws(s.peek()); s.next() {
	}
	s.out <- &token{t: t_ws}
}

func (s *scanner) emit_param() {
	offs := s.off
	for ; s.is_digit(s.peek()); s.next() {
	}
	s.out <- &token{t: t_param, v: s.in[offs:s.off]}
}

func (s *scanner) emit_paramx() {
	s.next() // read '{'

	offs := s.off
	for ; s.ch != '}' && !s.eof; s.next() {
	}
	if s.ch == '}' {
		s.out <- &token{t: t_paramx, v: s.in[offs : s.off-1]}
		return
	}
	s.out <- &token{t: t_invalid, v: string(s.ch), p: s.off}
}

func (s *scanner) emit_ptr() {
	s.out <- &token{t: t_ptr}
}

func (s *scanner) emit_amp() {
	s.out <- &token{t: t_amp}
}

func (s *scanner) emit_land() {
	s.next()
	s.out <- &token{t: t_land}
}

func (s *scanner) emit_lor() {
	s.next()
	s.out <- &token{t: t_lor}
}

func (s *scanner) emit_eql() {
	s.next()
	s.out <- &token{t: t_eql}
}

func (s *scanner) emit_assign() {
	s.out <- &token{t: t_assign}
}

func (s *scanner) emit_geq() {
	s.next()
	s.out <- &token{t: t_geq}
}

func (s *scanner) emit_gtr() {
	s.out <- &token{t: t_gtr}
}

func (s *scanner) emit_leq() {
	s.next()
	s.out <- &token{t: t_leq}
}

func (s *scanner) emit_lss() {
	s.out <- &token{t: t_lss}
}

func (s *scanner) emit_neq() {
	s.next()
	s.out <- &token{t: t_neq}
}

func (s *scanner) emit_not() {
	s.out <- &token{t: t_not}
}

func (s *scanner) emit_define() {
	s.next()
	s.out <- &token{t: t_define}
}

func (s *scanner) emit_inc() {
	s.next()
	s.out <- &token{t: t_inc}
}

func (s *scanner) emit_dec() {
	s.next()
	s.out <- &token{t: t_dec}
}

func (s *scanner) emit_lparen()    { s.out <- &token{t: t_lparen} }
func (s *scanner) emit_lbrack()    { s.out <- &token{t: t_lbrack} }
func (s *scanner) emit_lbrace()    { s.out <- &token{t: t_lbrace} }
func (s *scanner) emit_rparen()    { s.out <- &token{t: t_rparen} }
func (s *scanner) emit_rbrack()    { s.out <- &token{t: t_rbrack} }
func (s *scanner) emit_rbrace()    { s.out <- &token{t: t_rbrace} }
func (s *scanner) emit_comma()     { s.out <- &token{t: t_comma} }
func (s *scanner) emit_period()    { s.out <- &token{t: t_period} }
func (s *scanner) emit_semicolon() { s.out <- &token{t: t_semicolon} }

func (s *scanner) is_ws(c byte) bool {
	return c == ' ' || c == '\t'
}

func (s *scanner) is_digit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (s *scanner) is_letter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *scanner) is_ident_char(c byte) bool {
	return s.is_letter(c) || s.is_digit(c)
}
