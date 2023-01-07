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
		if isDigit(c2) {
			j := i + 2
			for ; j < len(text) && isDigit(text[j]); j++ {
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

		// when "$<key>"; write vars[key]
		if isAlnum(c2) {
			j := i + 2
			for ; j < len(text) && isAlnum(text[j]); j++ {
			}

			key := text[i+1 : j]
			g.S(g.vars[key])

			i = j - 1
			continue
		}

		// when ${<key>}; write vars[key]
		if c2 == '{' {
			j := i + 2
			for ; j < len(text) && text[j] != '}'; j++ {
			}

			key := text[i+2 : j]
			g.S(g.vars[key])

			i = j
			continue
		}
	}

	if g.werr == nil {
		_, g.werr = g.buf.Write(out)
	}
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

	case *rules.Rule:
		o := g.info.RuleObjMap[v]
		f := g.info.ObjFieldMap[o]
		g.genIsRuleExpr(f, o, v)

	case *rules.Arg:
		g.genArg(v)

	case func():
		v()
	}
}

func isDigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isAlnum(c byte) bool {
	return c == '_' ||
		(c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}
