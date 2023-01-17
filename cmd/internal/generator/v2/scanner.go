package generate

// The set of tokens supported by the generator. Most of these are the subset
// of lexial tokens of the Go language, in addition to those there are the
// additional parameter tokens.
type token_type uint8

func (t token_type) String() string {
	return _token_types[t]
}

const (
	t_invalid token_type = iota
	t_eof
	t_ws      // ' ', '\t' only, line whitespace is not supported
	t_comment // "//" only

	t_ident
	t_int    // 123
	t_float  // 12.3
	t_char   // 'a'
	t_string // "foo bar"

	t_param  // positional parameter, e.g. $0, $1, ...
	t_paramx // parameter expansion, e.g. ${0:e}

	t_add // +
	t_sub // -
	t_mul // *
	t_quo // /
	t_rem // %

	t_and // &
	t_or  // |

	t_land // &&
	t_lor  // ||
	t_inc  // ++
	t_dec  // --

	t_eql    // ==
	t_lss    // <
	t_gtr    // >
	t_assign // =
	t_not    // !

	t_neq      // !=
	t_leq      // <=
	t_geq      // >=
	t_define   // :=
	t_ellipsis // ...

	t_lparen // (
	t_rparen // )
	t_lbrack // [
	t_rbrack // ]
	t_lbrace // {
	t_rbrace // }

	t_comma     // ,
	t_period    // .
	t_semicolon // ;
	t_colon     // :

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
	t_ident:     "<identifier token>",
	t_int:       "<int token>",
	t_float:     "<float token>",
	t_char:      "<char token>",
	t_string:    "<string token>",
	t_param:     "<param token>",
	t_paramx:    "<paramx token>",
	t_add:       "+",
	t_sub:       "-",
	t_mul:       "*",
	t_quo:       "/",
	t_rem:       "%",
	t_and:       "&",
	t_or:        "|",
	t_land:      "&&",
	t_lor:       "||",
	t_inc:       "++",
	t_dec:       "--",
	t_eql:       "==",
	t_lss:       "<",
	t_gtr:       ">",
	t_assign:    "=",
	t_not:       "!",
	t_neq:       "!=",
	t_leq:       "<=",
	t_geq:       ">=",
	t_define:    ":=",
	t_ellipsis:  "...",
	t_lparen:    "(",
	t_rparen:    ")",
	t_lbrack:    "[",
	t_rbrack:    "]",
	t_lbrace:    "{",
	t_rbrace:    "}",
	t_comma:     ",",
	t_period:    ".",
	t_semicolon: ";",
	t_colon:     ":",
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
	p int        // position of the token
	v string     // literal value; used only when needed
}

// The scanner for the generator's input. Note that the implementation of the
// scanner is quite naive and not at all as sophisticated as Go's own scanner.
type scanner struct {
	in  string
	out chan<- *token

	stop chan struct{}
	exit bool

	ch  byte // current character (ascii only)
	off int  // current offset
	eof bool
}

func (s *scanner) run() {
	for {
		if s.exit {
			return
		}

		s.next()
		if s.eof {
			s.emit(&token{t: t_eof})
			return
		}

		p := s.off - 1
		switch c, c2 := s.ch, s.peek(); {
		case s.is_letter(c):
			p, v := s.scan_ident()
			switch v {
			default:
				s.emit(&token{t: t_ident, p: p, v: v})
			case "if":
				s.emit(&token{t: t_if, p: p})
			case "else":
				s.emit(&token{t: t_else, p: p})
			case "for":
				s.emit(&token{t: t_for, p: p})
			case "range":
				s.emit(&token{t: t_range, p: p})
			case "return":
				s.emit(&token{t: t_return, p: p})
			case "package":
				s.emit(&token{t: t_package, p: p})
			case "import":
				s.emit(&token{t: t_import, p: p})
			case "func":
				s.emit(&token{t: t_func, p: p})
			}

		case s.is_digit(c):
			t, p, v := s.scan_number()
			s.emit(&token{t: t, p: p, v: v})

		case s.is_ws(c):
			p, v := s.scan_ws()
			s.emit(&token{t: t_ws, p: p, v: v})

		case c == '\'':
			t, p, v := s.scan_char()
			s.emit(&token{t: t, p: p, v: v})

		case c == '"':
			t, p, v := s.scan_string()
			s.emit(&token{t: t, p: p, v: v})

		case c == '$':
			t, p, v := s.scan_param()
			s.emit(&token{t: t, p: p, v: v})

		case c == '.':
			t, p, v := s.scan_dots()
			s.emit(&token{t: t, p: p, v: v})

		case c == '/':
			t, p, v := s.scan_quo()
			s.emit(&token{t: t, p: p, v: v})

		case c == '+' && c2 == '+':
			s.next_emit(&token{t: t_inc, p: p})
		case c == '+':
			s.emit(&token{t: t_add, p: p})
		case c == '-' && c2 == '-':
			s.next_emit(&token{t: t_dec, p: p})
		case c == '-':
			s.emit(&token{t: t_sub, p: p})
		case c == '*':
			s.emit(&token{t: t_mul, p: p})
		case c == '%':
			s.emit(&token{t: t_rem, p: p})
		case c == '&' && c2 == '&':
			s.next_emit(&token{t: t_land, p: p})
		case c == '&':
			s.emit(&token{t: t_and, p: p})
		case c == '|' && c2 == '|':
			s.next_emit(&token{t: t_lor, p: p})
		case c == '|':
			s.emit(&token{t: t_or, p: p})
		case c == '=' && c2 == '=':
			s.next_emit(&token{t: t_eql, p: p})
		case c == '=':
			s.emit(&token{t: t_assign, p: p})
		case c == '<' && c2 == '=':
			s.next_emit(&token{t: t_leq, p: p})
		case c == '<':
			s.emit(&token{t: t_lss, p: p})
		case c == '>' && c2 == '=':
			s.next_emit(&token{t: t_geq, p: p})
		case c == '>':
			s.emit(&token{t: t_gtr, p: p})
		case c == '!' && c2 == '=':
			s.next_emit(&token{t: t_neq, p: p})
		case c == '!':
			s.emit(&token{t: t_not, p: p})
		case c == ':' && c2 == '=':
			s.next_emit(&token{t: t_define, p: p})
		case c == ':':
			s.emit(&token{t: t_colon, p: p})
		case c == '(':
			s.emit(&token{t: t_lparen, p: p})
		case c == '[':
			s.emit(&token{t: t_lbrack, p: p})
		case c == '{':
			s.emit(&token{t: t_lbrace, p: p})
		case c == ')':
			s.emit(&token{t: t_rparen, p: p})
		case c == ']':
			s.emit(&token{t: t_rbrack, p: p})
		case c == '}':
			s.emit(&token{t: t_rbrace, p: p})
		case c == ',':
			s.emit(&token{t: t_comma, p: p})
		case c == ';':
			s.emit(&token{t: t_semicolon, p: p})
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

func (s *scanner) emit(t *token) {
	if t.t == t_invalid {
		s.exit = true
	}

	select {
	case <-s.stop:
		s.exit = true

	case s.out <- t:
		// ok;
	}
}

func (s *scanner) next_emit(t *token) {
	s.next()
	s.emit(t)
}

func (s *scanner) scan_ident() (p int, v string) {
	offs := s.off - 1
	for ; s.is_ident_char(s.peek()); s.next() {
	}
	return offs, s.in[offs:s.off]
}

func (s *scanner) scan_ws() (p int, v string) {
	offs := s.off - 1
	for ; s.is_ws(s.peek()); s.next() {
	}
	return offs, s.in[offs:s.off]
}

func (s *scanner) scan_char() (t token_type, p int, v string) {
	offs := s.off
	for ; s.ch != '\'' && !s.eof; s.next() {
	}
	if s.ch != '\'' {
		return t_invalid, s.off, string(s.ch)
	}
	return t_char, offs - 1, s.in[offs : s.off-1]
}

func (s *scanner) scan_string() (t token_type, p int, v string) {
	offs := s.off
	for ; s.ch != '"' && !s.eof; s.next() {
	}
	if s.ch != '"' {
		return t_invalid, s.off, string(s.ch)
	}
	return t_string, offs - 1, s.in[offs : s.off-1]
}

func (s *scanner) scan_dots() (t token_type, p int, v string) {
	if s.peek() != '.' {
		return t_period, s.off, ""
	}
	s.next()
	s.next()
	if s.ch == '.' { // "..."?
		return t_ellipsis, s.off - 2, ""
	}
	return t_invalid, s.off, string(s.ch)

}

func (s *scanner) scan_quo() (t token_type, p int, v string) {
	if s.peek() == '/' { // comment?
		offs := s.off + 1
		s.off = len(s.in)
		return t_comment, offs - 2, s.in[offs:]
	}
	return t_quo, s.off, "" // quotient op

}

func (s *scanner) scan_number() (token_type, int, string) {
	offs := s.off - 1
	with_fp := false

	for {
		prev_ch := s.ch
		s.next()

		switch {
		case s.ch == '.' && !with_fp:
			with_fp = true
		case s.ch == '_' && prev_ch != '_':
			// ok
		case s.is_digit(s.ch):
			// ok
		default:
			return t_invalid, s.off, string(s.ch)
		}
	}

	if with_fp {
		return t_float, offs, s.in[offs:s.off]
	}
	return t_int, offs, s.in[offs:s.off]
}

func (s *scanner) scan_param() (t token_type, p int, v string) {
	switch s.next(); {
	case s.is_digit(s.ch): // index only?
		offs := s.off - 1
		for ; s.is_digit(s.peek()); s.next() {
		}
		return t_param, offs - 1, s.in[offs:s.off]

	case s.ch == '{': // with expansion?
		offs := s.off
		for ; s.ch != '}' && !s.eof; s.next() {
		}
		if s.ch != '}' {
			return t_invalid, s.off, string(s.ch)
		}
		return t_paramx, offs - 2, s.in[offs : s.off-1]
	}

	return t_invalid, s.off, string(s.ch)
}

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
