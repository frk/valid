package generate

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type writer struct {
	buf  bytes.Buffer
	vars map[string]string
	err  error
}

func (w *writer) from(r *writer) {
	if w.err != nil {
		return
	}

	_, w.err = io.Copy(&w.buf, &r.buf)
}

func (w *writer) ln(line string, args ...any) {
	if w.err != nil {
		return
	}
	w.p(line, args...)
	w.nl()
}

func (w *writer) nl() {
	if w.err != nil {
		return
	}
	w.err = w.buf.WriteByte('\n')
}

func (w *writer) p(text string, args ...any) {
	if w.err != nil {
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

		c2 := text[i+1]

		// when "$<index>"; write args[index]
		if isDigit(c2) {
			j := i + 2
			for ; j < len(text) && isDigit(text[j]); j++ {
			}

			n, err := strconv.Atoi(text[i+1 : j])
			if err != nil {
				w.err = fmt.Errorf("failed to parse %s: %v", text[i:j], err)
				return
			} else if n >= len(args) {
				w.err = fmt.Errorf("index %s out of range", text[i:j])
				return
			}

			out = append(out, fmt.Sprintf("%v", args[n])...)
			i = j - 1
			continue
		}

		// when "$<key>"; write vars[key]
		if isAlnum(c2) {
			j := i + 2
			for ; j < len(text) && isAlnum(text[j]); j++ {
			}

			key := text[i+1 : j]
			out = append(out, fmt.Sprintf("%v", w.vars[key])...)
			i = j - 1
			continue
		}

		// when ${<key>}; write vars[key]
		if c2 == '{' {
			j := i + 2
			for ; j < len(text) && text[j] != '}'; j++ {
			}

			key := text[i+2 : j]
			out = append(out, fmt.Sprintf("%v", w.vars[key])...)
			i = j
			continue
		}
	}

	if w.err == nil {
		_, w.err = w.buf.Write(out)
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
