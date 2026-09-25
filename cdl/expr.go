package cdl

import (
	"fmt"
	"strconv"
	"strings"
)

// exprParser parses expressions from the frozen pilot grammar. The only failure
// codes it produces are CDL_PARSE and CDL_UNKNOWN_INSTRUCTION.
type exprParser struct {
	toks []token
	pos  int
}

type token struct {
	kind string // ident, int, str, op, lparen, rparen, comma
	val  string
}

func parseExprSource(src string) (Expr, error) {
	toks, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	p := &exprParser{toks: toks}
	e, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.toks) {
		return nil, fmt.Errorf("CDL_PARSE: unexpected token %q", p.toks[p.pos].val)
	}
	return e, nil
}

func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || c == '.' || c == '-'
}

func tokenize(s string) ([]token, error) {
	var toks []token
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == ' ' || c == '\t':
			i++
		case c == '"':
			j := strings.IndexByte(s[i+1:], '"')
			if j < 0 {
				return nil, fmt.Errorf("CDL_PARSE: unterminated string")
			}
			toks = append(toks, token{kind: "str", val: s[i+1 : i+1+j]})
			i += j + 2
		case strings.HasPrefix(s[i:], "=="), strings.HasPrefix(s[i:], "!="),
			strings.HasPrefix(s[i:], "<="), strings.HasPrefix(s[i:], ">="):
			toks = append(toks, token{kind: "op", val: s[i : i+2]})
			i += 2
		case c == '<' || c == '>' || c == '+' || c == '-' || c == '(' || c == ')' || c == ',':
			kind := "op"
			if c == '(' {
				kind = "lparen"
			} else if c == ')' {
				kind = "rparen"
			} else if c == ',' {
				kind = "comma"
			}
			toks = append(toks, token{kind: kind, val: string(c)})
			i++
		case isIdentStart(c):
			j := i
			for j < len(s) && isIdentPart(s[j]) {
				if s[j] == '-' && (j == i || s[j-1] == ' ') {
					break
				}
				j++
			}
			word := s[i:j]
			if isAllDigits(word) {
				toks = append(toks, token{kind: "int", val: word})
			} else {
				toks = append(toks, token{kind: "ident", val: word})
			}
			i = j
		default:
			return nil, fmt.Errorf("CDL_PARSE: unexpected character %q", string(c))
		}
	}
	return toks, nil
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > 0
}

func (p *exprParser) peek() (token, bool) {
	if p.pos >= len(p.toks) {
		return token{}, false
	}
	return p.toks[p.pos], true
}

func (p *exprParser) next() (token, bool) {
	t, ok := p.peek()
	if ok {
		p.pos++
	}
	return t, ok
}

func (p *exprParser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != "ident" || t.val != "OR" {
			return left, nil
		}
		p.pos++
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = Binary{Op: "OR", L: left, R: right}
	}
}

func (p *exprParser) parseAnd() (Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != "ident" || t.val != "AND" {
			return left, nil
		}
		p.pos++
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = Binary{Op: "AND", L: left, R: right}
	}
}

func (p *exprParser) parseNot() (Expr, error) {
	t, ok := p.peek()
	if ok && t.kind == "ident" && t.val == "NOT" {
		p.pos++
		e, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return Unary{Op: "NOT", E: e}, nil
	}
	return p.parseCompare()
}

func (p *exprParser) parseCompare() (Expr, error) {
	left, err := p.parseAdd()
	if err != nil {
		return nil, err
	}
	t, ok := p.peek()
	if !ok || t.kind != "op" {
		return left, nil
	}
	switch t.val {
	case "==", "!=", "<", "<=", ">", ">=":
		p.pos++
		right, err := p.parseAdd()
		if err != nil {
			return nil, err
		}
		return Binary{Op: t.val, L: left, R: right}, nil
	}
	return left, nil
}

func (p *exprParser) parseAdd() (Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != "op" || (t.val != "+" && t.val != "-") {
			return left, nil
		}
		p.pos++
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = Binary{Op: t.val, L: left, R: right}
	}
}

func (p *exprParser) parsePrimary() (Expr, error) {
	t, ok := p.next()
	if !ok {
		return nil, fmt.Errorf("CDL_PARSE: unexpected end of expression")
	}
	switch t.kind {
	case "int":
		n, err := strconv.Atoi(t.val)
		if err != nil {
			return nil, fmt.Errorf("CDL_PARSE: invalid integer %q", t.val)
		}
		return Literal{Value: n}, nil
	case "str":
		return Literal{Value: t.val}, nil
	case "lparen":
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		close, ok := p.next()
		if !ok || close.kind != "rparen" {
			return nil, fmt.Errorf("CDL_PARSE: missing closing parenthesis")
		}
		return e, nil
	case "ident":
		switch t.val {
		case "true":
			return Literal{Value: true}, nil
		case "false":
			return Literal{Value: false}, nil
		case "COMPUTE":
			return nil, cdlerr("CDL_UNKNOWN_INSTRUCTION", "COMPUTE is not in the frozen grammar")
		}
		if next, ok := p.peek(); ok && next.kind == "lparen" {
			return p.parseCall(t.val)
		}
		return Ident{Raw: t.val}, nil
	default:
		return nil, fmt.Errorf("CDL_PARSE: unexpected token %q", t.val)
	}
}

func (p *exprParser) parseCall(name string) (Expr, error) {
	p.pos++ // consume lparen
	call := Call{Name: name}
	if name == "COUNT" {
		table, ok := p.next()
		if !ok || table.kind != "ident" {
			return nil, fmt.Errorf("CDL_PARSE: COUNT expects a table identifier")
		}
		call.Args = []Expr{Ident{Raw: table.val}}
		if w, ok := p.peek(); ok && w.kind == "ident" && w.val == "WHERE" {
			p.pos++
			field, ok := p.next()
			if !ok || field.kind != "ident" {
				return nil, fmt.Errorf("CDL_PARSE: WHERE expects a field identifier")
			}
			eq, ok := p.next()
			if !ok || eq.kind != "op" || eq.val != "==" {
				return nil, fmt.Errorf("CDL_PARSE: WHERE expects '=='")
			}
			value, err := p.parseOr()
			if err != nil {
				return nil, err
			}
			call.Where = &Where{Field: field.val, Value: value}
		}
		close, ok := p.next()
		if !ok || close.kind != "rparen" {
			return nil, fmt.Errorf("CDL_PARSE: missing closing parenthesis")
		}
		return call, nil
	}
	if next, ok := p.peek(); ok && next.kind == "rparen" {
		p.pos++
		return call, nil
	}
	for {
		arg, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, arg)
		t, ok := p.next()
		if !ok {
			return nil, fmt.Errorf("CDL_PARSE: missing closing parenthesis")
		}
		if t.kind == "rparen" {
			return call, nil
		}
		if t.kind != "comma" {
			return nil, fmt.Errorf("CDL_PARSE: unexpected token %q in call", t.val)
		}
	}
}
