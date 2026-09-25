package cdl

import (
	"errors"
	"fmt"
	"strings"
)

// Diagnostic is one stable compiler diagnostic.
type Diagnostic struct {
	Code string
	File string
	Line int
	Msg  string
}

// Diagnostics is an ordered diagnostic list.
type Diagnostics []Diagnostic

func (d Diagnostics) Error() string {
	if len(d) == 0 {
		return "no diagnostics"
	}
	first := fmt.Sprintf("%s: %s", d[0].Code, d[0].Msg)
	if d[0].File != "" {
		first = fmt.Sprintf("%s:%d: %s", d[0].File, d[0].Line, first)
	}
	if len(d) > 1 {
		return fmt.Sprintf("%s (and %d more)", first, len(d)-1)
	}
	return first
}

// Has reports whether any diagnostic carries the code.
func (d Diagnostics) Has(code string) bool {
	for _, diag := range d {
		if diag.Code == code {
			return true
		}
	}
	return false
}

type parser struct {
	file  string
	lines []string
	pos   int
	diags Diagnostics
}

// Parse parses one CDL source file.
func Parse(file string, src []byte) (*Program, Diagnostics) {
	text := strings.ReplaceAll(string(src), "\r\n", "\n")
	p := &parser{file: file, lines: strings.Split(text, "\n")}
	prog := &Program{}
	for {
		line, ok := p.peek()
		if !ok {
			return prog, p.diags
		}
		if line == "" {
			p.pos++
			continue
		}
		switch {
		case line == "GLOBAL DECLARATIONS":
			p.pos++
			if !p.parseGlobals(&prog.Globals) {
				return prog, p.diags
			}
		case strings.HasPrefix(line, "SECTION "):
			p.pos++
			sec, ok := p.parseSection(strings.TrimSpace(strings.TrimPrefix(line, "SECTION ")))
			if !ok {
				return prog, p.diags
			}
			prog.Sections = append(prog.Sections, sec)
		case strings.HasPrefix(line, "PROJECTION "):
			p.pos++
			proj, ok := p.parseProjection(strings.TrimSpace(strings.TrimPrefix(line, "PROJECTION ")))
			if !ok {
				return prog, p.diags
			}
			prog.Projections = append(prog.Projections, proj)
		default:
			p.failHere("CDL_PARSE", "unexpected top-level line %q", line)
			return prog, p.diags
		}
	}
}

func (p *parser) fail(line int, code, format string, args ...any) {
	p.diags = append(p.diags, Diagnostic{Code: code, File: p.file, Line: line, Msg: fmt.Sprintf(format, args...)})
}

func (p *parser) failHere(code, format string, args ...any) {
	p.fail(p.pos+1, code, format, args...)
}

func stripComment(line string) string {
	inQuote := false
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			inQuote = !inQuote
		case '#':
			if !inQuote {
				return strings.TrimRight(line[:i], " \t")
			}
		}
	}
	return strings.TrimRight(line, " \t")
}

func (p *parser) peek() (string, bool) {
	if p.pos >= len(p.lines) {
		return "", false
	}
	return strings.TrimSpace(stripComment(p.lines[p.pos])), true
}

func (p *parser) isTopLevel(line string) bool {
	return line == "GLOBAL DECLARATIONS" || strings.HasPrefix(line, "SECTION ") || strings.HasPrefix(line, "PROJECTION ")
}

func csv(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s[1 : len(s)-1]
	}
	return s
}

func builtinCapabilityKind(id string) string {
	switch id {
	case "hash.sha256", "table.serialize":
		return "deterministic"
	default:
		return "proposing"
	}
}

func (p *parser) readOpaque() string {
	var parts []string
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated text block")
			return ""
		}
		p.pos++
		if line == "END" {
			return strings.Join(parts, " ")
		}
		if line != "" {
			parts = append(parts, line)
		}
	}
}

func (p *parser) readList() []string {
	var out []string
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated list block")
			return out
		}
		p.pos++
		if line == "END" {
			return out
		}
		if line != "" {
			out = append(out, line)
		}
	}
}

func (p *parser) parseGlobals(g *Globals) bool {
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated GLOBAL DECLARATIONS")
			return false
		}
		if line == "END" {
			p.pos++
			return true
		}
		if line == "" {
			p.pos++
			continue
		}
		p.pos++
		switch {
		case strings.HasPrefix(line, "CAPABILITY "):
			id := strings.TrimSpace(strings.TrimPrefix(line, "CAPABILITY "))
			g.Capabilities = append(g.Capabilities, Capability{ID: id, Kind: builtinCapabilityKind(id)})
		case strings.HasPrefix(line, "REGISTRY "):
			g.Registries = append(g.Registries, strings.TrimSpace(strings.TrimPrefix(line, "REGISTRY ")))
		case strings.HasPrefix(line, "ARTIFACT "):
			g.Artifacts = append(g.Artifacts, strings.TrimSpace(strings.TrimPrefix(line, "ARTIFACT ")))
		case strings.HasPrefix(line, "RULE "):
			g.Rules = append(g.Rules, strings.TrimSpace(strings.TrimPrefix(line, "RULE ")))
		case strings.HasPrefix(line, "GATES "):
			g.Gates = append(g.Gates, csv(strings.TrimPrefix(line, "GATES "))...)
		case strings.HasPrefix(line, "WORKFLOW-TARGET "):
			g.WorkflowTargets = append(g.WorkflowTargets, csv(strings.TrimPrefix(line, "WORKFLOW-TARGET "))...)
		default:
			p.failHere("CDL_PARSE", "unexpected global declaration %q", line)
			return false
		}
	}
}

func (p *parser) parseSection(id string) (Section, bool) {
	sec := Section{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			return sec, true
		}
		if line == "" {
			p.pos++
			continue
		}
		if p.isTopLevel(line) {
			return sec, true
		}
		p.pos++
		switch {
		case strings.HasPrefix(line, "NUMBER "):
			sec.Number = strings.TrimSpace(strings.TrimPrefix(line, "NUMBER "))
		case strings.HasPrefix(line, "TITLE "):
			sec.Title = unquote(strings.TrimPrefix(line, "TITLE "))
		case strings.HasPrefix(line, "ARTIFACT "):
			sec.Artifact = strings.TrimSpace(strings.TrimPrefix(line, "ARTIFACT "))
		case line == "GOAL":
			sec.Goal = p.readOpaque()
		case strings.HasPrefix(line, "USES "):
			rest := strings.TrimSpace(strings.TrimPrefix(line, "USES "))
			fields := strings.Fields(rest)
			if len(fields) != 2 {
				p.failHere("CDL_PARSE", "USES expects <KIND> <id>")
				return sec, false
			}
			sec.Uses = append(sec.Uses, Use{Kind: fields[0], ID: fields[1]})
		case strings.HasPrefix(line, "REQUIRES CAPABILITY "):
			sec.Requires = csv(strings.TrimPrefix(line, "REQUIRES CAPABILITY "))
		case strings.HasPrefix(line, "TYPE "):
			rest := strings.TrimSpace(strings.TrimPrefix(line, "TYPE "))
			parts := strings.SplitN(rest, ":=", 2)
			if len(parts) != 2 {
				p.failHere("CDL_PARSE", "TYPE expects <id> := <type>")
				return sec, false
			}
			sec.Types = append(sec.Types, TypeDecl{ID: strings.TrimSpace(parts[0]), Type: strings.TrimSpace(parts[1])})
		case strings.HasPrefix(line, "STATE "):
			st, ok := p.parseState(strings.TrimSpace(strings.TrimPrefix(line, "STATE ")))
			if !ok {
				return sec, false
			}
			sec.States = append(sec.States, st)
		case strings.HasPrefix(line, "VALUE "):
			vd, ok := p.parseValue(strings.TrimSpace(strings.TrimPrefix(line, "VALUE ")))
			if !ok {
				return sec, false
			}
			sec.Values = append(sec.Values, vd)
		case strings.HasPrefix(line, "ENUM "):
			ed, ok := p.parseEnum(strings.TrimSpace(strings.TrimPrefix(line, "ENUM ")))
			if !ok {
				return sec, false
			}
			sec.Enums = append(sec.Enums, ed)
		case strings.HasPrefix(line, "FIELD "):
			fd, ok := p.parseField(strings.TrimSpace(strings.TrimPrefix(line, "FIELD ")))
			if !ok {
				return sec, false
			}
			sec.Fields = append(sec.Fields, fd)
		case strings.HasPrefix(line, "TABLE "):
			td, ok := p.parseTable(strings.TrimSpace(strings.TrimPrefix(line, "TABLE ")))
			if !ok {
				return sec, false
			}
			sec.Tables = append(sec.Tables, td)
		case strings.HasPrefix(line, "RULE "):
			rd, ok := p.parseRule(strings.TrimSpace(strings.TrimPrefix(line, "RULE ")))
			if !ok {
				return sec, false
			}
			sec.Rules = append(sec.Rules, rd)
		case strings.HasPrefix(line, "STEP "):
			st, ok := p.parseStep(strings.TrimSpace(strings.TrimPrefix(line, "STEP ")))
			if !ok {
				return sec, false
			}
			sec.Steps = append(sec.Steps, st)
		default:
			p.failHere("CDL_PARSE", "unexpected section line %q", line)
			return sec, false
		}
	}
}

func (p *parser) parseState(id string) (StateDecl, bool) {
	st := StateDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated STATE %s", id)
			return st, false
		}
		p.pos++
		switch {
		case line == "END":
			return st, true
		case line == "":
		case strings.HasPrefix(line, "VALUES "):
			st.Values = csv(strings.TrimPrefix(line, "VALUES "))
		case strings.HasPrefix(line, "INITIAL "):
			st.Initial = strings.TrimSpace(strings.TrimPrefix(line, "INITIAL "))
		case line == "ALLOWS":
			for {
				inner, ok := p.peek()
				if !ok {
					p.failHere("CDL_PARSE", "unterminated ALLOWS block")
					return st, false
				}
				p.pos++
				if inner == "END" {
					break
				}
				if inner == "" {
					continue
				}
				parts := strings.SplitN(inner, "->", 2)
				if len(parts) != 2 {
					p.failHere("CDL_PARSE", "ALLOWS expects <value> -> <value>")
					return st, false
				}
				st.Allows = append(st.Allows, [2]string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
			}
		default:
			p.failHere("CDL_PARSE", "unexpected STATE line %q", line)
			return st, false
		}
	}
}

func (p *parser) parseValue(id string) (ValueDecl, bool) {
	vd := ValueDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated VALUE %s", id)
			return vd, false
		}
		p.pos++
		switch {
		case line == "END":
			return vd, true
		case line == "":
		case strings.HasPrefix(line, "TYPE "):
			vd.Type = strings.TrimSpace(strings.TrimPrefix(line, "TYPE "))
		case strings.HasPrefix(line, "LIFETIME "):
			vd.Lifetime = strings.TrimSpace(strings.TrimPrefix(line, "LIFETIME "))
		default:
			p.failHere("CDL_PARSE", "unexpected VALUE line %q", line)
			return vd, false
		}
	}
}

func (p *parser) parseEnum(id string) (EnumDecl, bool) {
	ed := EnumDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated ENUM %s", id)
			return ed, false
		}
		p.pos++
		if line == "END" {
			return ed, true
		}
		if line != "" {
			ed.Values = append(ed.Values, line)
		}
	}
}

func (p *parser) parseField(id string) (FieldDecl, bool) {
	fd := FieldDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated FIELD %s", id)
			return fd, false
		}
		p.pos++
		if line == "END" {
			return fd, true
		}
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "->", 2)
		if len(parts) != 2 {
			p.failHere("CDL_PARSE", "FIELD row expects <name> -> <type> required|optional")
			return fd, false
		}
		name := strings.TrimSpace(parts[0])
		rest := strings.Fields(parts[1])
		if len(rest) != 2 {
			p.failHere("CDL_PARSE", "FIELD row expects <type> required|optional")
			return fd, false
		}
		fd.Fields = append(fd.Fields, FieldSpec{Name: name, Type: rest[0], Required: rest[1] == "required"})
	}
}

func (p *parser) parseTable(id string) (TableDecl, bool) {
	td := TableDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated TABLE %s", id)
			return td, false
		}
		p.pos++
		switch {
		case line == "END":
			return td, true
		case line == "":
		case strings.HasPrefix(line, "ROWS "):
			td.Rows = strings.TrimSpace(strings.TrimPrefix(line, "ROWS "))
		case strings.HasPrefix(line, "ROW-TYPE "):
			td.RowType = strings.TrimSpace(strings.TrimPrefix(line, "ROW-TYPE "))
		case strings.HasPrefix(line, "KEY "):
			td.Key = strings.TrimSpace(strings.TrimPrefix(line, "KEY "))
		default:
			p.failHere("CDL_PARSE", "unexpected TABLE line %q", line)
			return td, false
		}
	}
}

func (p *parser) parseRule(id string) (RuleDecl, bool) {
	rd := RuleDecl{ID: id}
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated RULE %s", id)
			return rd, false
		}
		p.pos++
		switch {
		case line == "END":
			return rd, true
		case line == "":
		case line == "GOAL":
			rd.Goal = p.readOpaque()
		case strings.HasPrefix(line, "PREDICATE "):
			rest := strings.TrimSpace(strings.TrimPrefix(line, "PREDICATE "))
			rd.Predicate = strings.TrimSpace(strings.TrimSuffix(rest, ":="))
			var parts []string
			for {
				inner, ok := p.peek()
				if !ok {
					p.failHere("CDL_PARSE", "unterminated PREDICATE %s", id)
					return rd, false
				}
				p.pos++
				if inner == "END" {
					break
				}
				if inner != "" {
					parts = append(parts, inner)
				}
			}
			expr, err := parseExprSource(strings.Join(parts, " "))
			if err != nil {
				p.failHere(diagCode(err), "%s", diagMessage(err))
				return rd, false
			}
			rd.Expr = expr
			return rd, true
		default:
			p.failHere("CDL_PARSE", "unexpected RULE line %q", line)
			return rd, false
		}
	}
}

func (p *parser) parseStep(id string) (Step, bool) {
	step := Step{ID: id}
	line, ok := p.peek()
	if !ok {
		p.failHere("CDL_PARSE", "unterminated STEP %s", id)
		return step, false
	}
	p.pos++
	switch line {
	case "MACHINE:":
		step.Owner = "MACHINE"
		ops, uses, ok := p.parseMachineBody()
		if !ok {
			return step, false
		}
		step.Ops = ops
		step.UsesRules = uses
	case "AGENT:":
		step.Owner = "AGENT"
		if !p.parseAgentBody(&step) {
			return step, false
		}
	default:
		p.failHere("CDL_PARSE", "STEP %s expects MACHINE: or AGENT:", id)
		return step, false
	}
	return step, true
}

func (p *parser) parseMachineBody() ([]Op, []string, bool) {
	var ops []Op
	var uses []string
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated MACHINE block")
			return nil, nil, false
		}
		if line == "END" {
			p.pos++
			return ops, uses, true
		}
		if line == "" {
			p.pos++
			continue
		}
		p.pos++
		if strings.HasPrefix(line, "USES RULE ") {
			uses = append(uses, strings.TrimSpace(strings.TrimPrefix(line, "USES RULE ")))
			continue
		}
		op, err := p.parseOp(line)
		if err != nil {
			if errors.Is(err, errAborted) {
				return nil, nil, false
			}
			p.failHere(diagCode(err), "%s", diagMessage(err))
			return nil, nil, false
		}
		ops = append(ops, op)
	}
}

func (p *parser) parseOpsUntil(terminators ...string) ([]Op, string, bool) {
	var ops []Op
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated block")
			return nil, "", false
		}
		for _, term := range terminators {
			if line == term {
				p.pos++
				return ops, term, true
			}
		}
		if line == "" {
			p.pos++
			continue
		}
		p.pos++
		op, err := p.parseOp(line)
		if err != nil {
			if errors.Is(err, errAborted) {
				return nil, "", false
			}
			p.failHere(diagCode(err), "%s", diagMessage(err))
			return nil, "", false
		}
		ops = append(ops, op)
	}
}

func (p *parser) parseOp(line string) (Op, error) {
	switch {
	case strings.HasPrefix(line, "REQUIRES CAPABILITY "):
		return RequiresOp{Caps: csv(strings.TrimPrefix(line, "REQUIRES CAPABILITY "))}, nil
	case strings.HasPrefix(line, "FOR EACH "):
		rest := strings.TrimSpace(strings.TrimPrefix(line, "FOR EACH "))
		fields := strings.Fields(rest)
		if len(fields) != 3 || fields[1] != "IN" {
			return nil, cdlerr("CDL_PARSE", "FOR EACH expects <var> IN <registry>")
		}
		body, _, ok := p.parseOpsUntil("END")
		if !ok {
			return nil, errAborted
		}
		return ForEachOp{Var: fields[0], In: fields[2], Body: body}, nil
	case strings.HasPrefix(line, "APPEND "):
		fields := strings.Fields(strings.TrimPrefix(line, "APPEND "))
		if len(fields) != 3 || fields[1] != "FOR" {
			return nil, cdlerr("CDL_PARSE", "APPEND expects <row-type> FOR <var>")
		}
		return AppendOp{RowType: fields[0], For: fields[2]}, nil
	case strings.HasPrefix(line, "LET "):
		target, expr, err := parseAssignment(strings.TrimPrefix(line, "LET "))
		if err != nil {
			return nil, err
		}
		return LetOp{Target: target, Value: expr}, nil
	case strings.HasPrefix(line, "SET "):
		target, expr, err := parseAssignment(strings.TrimPrefix(line, "SET "))
		if err != nil {
			return nil, err
		}
		return SetOp{Target: target, Value: expr}, nil
	case strings.HasPrefix(line, "GUARD "):
		expr, err := parseExprSource(strings.TrimPrefix(line, "GUARD "))
		if err != nil {
			return nil, err
		}
		guard := GuardOp{Expr: expr}
		if next, ok := p.peek(); ok && strings.HasPrefix(next, "OTHERWISE BLOCK ") {
			p.pos++
			guard.Otherwise = []Op{BlockOp{Gates: csv(strings.TrimPrefix(next, "OTHERWISE BLOCK "))}}
		}
		return guard, nil
	case strings.HasPrefix(line, "IF "):
		cond, err := parseExprSource(strings.TrimPrefix(line, "IF "))
		if err != nil {
			return nil, err
		}
		then, term, ok := p.parseOpsUntil("ELSE", "END")
		if !ok {
			return nil, errAborted
		}
		ifOp := IfOp{Cond: cond, Then: then}
		if term == "ELSE" {
			els, _, ok := p.parseOpsUntil("END")
			if !ok {
				return nil, errAborted
			}
			ifOp.Else = els
		}
		return ifOp, nil
	case strings.HasPrefix(line, "BLOCK "):
		return BlockOp{Gates: csv(strings.TrimPrefix(line, "BLOCK "))}, nil
	case strings.HasPrefix(line, "COMPUTE "):
		return nil, cdlerr("CDL_UNKNOWN_INSTRUCTION", "COMPUTE is not in the frozen grammar")
	default:
		return nil, cdlerr("CDL_PARSE", "unexpected MACHINE operation %q", line)
	}
}

func parseAssignment(rest string) (string, Expr, error) {
	parts := strings.SplitN(rest, ":=", 2)
	if len(parts) != 2 {
		return "", nil, cdlerr("CDL_PARSE", "assignment expects <target> := <expression>")
	}
	expr, err := parseExprSource(parts[1])
	if err != nil {
		return "", nil, err
	}
	return strings.TrimSpace(parts[0]), expr, nil
}

func (p *parser) parseAgentBody(step *Step) bool {
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated AGENT block")
			return false
		}
		p.pos++
		switch {
		case line == "END":
			return true
		case line == "":
		case line == "EVIDENCE":
			step.Evidence = p.readList()
		case line == "PRODUCES":
			step.Produces = p.readList()
		case line == "RESULT":
			res := p.readList()
			if len(res) > 0 {
				step.Result = res[0]
			}
		case line == "REQUIRE \"\"\"":
			var parts []string
			for {
				inner, ok := p.peek()
				if !ok {
					p.failHere("CDL_PARSE", "unterminated REQUIRE clause")
					return false
				}
				p.pos++
				if inner == "\"\"\"" {
					break
				}
				if inner != "" {
					parts = append(parts, inner)
				}
			}
			step.Require = strings.Join(parts, " ")
		default:
			p.failHere("CDL_PARSE", "unexpected AGENT line %q", line)
			return false
		}
	}
}

func (p *parser) parseProjection(id string) (Projection, bool) {
	proj := Projection{ID: id, StepTexts: map[string]string{}, RuleTexts: map[string]string{}}
	for {
		line, ok := p.peek()
		if !ok {
			return proj, true
		}
		if line == "" {
			p.pos++
			continue
		}
		if p.isTopLevel(line) {
			return proj, true
		}
		p.pos++
		switch {
		case line == "END":
			return proj, true
		case strings.HasPrefix(line, "CHANNELS "):
			proj.Channels = csv(strings.TrimPrefix(line, "CHANNELS "))
		case strings.HasPrefix(line, "CHANNEL "):
			proj.Channels = []string{strings.TrimSpace(strings.TrimPrefix(line, "CHANNEL "))}
		case strings.HasPrefix(line, "COVERS SECTION "):
			proj.Covers = strings.TrimSpace(strings.TrimPrefix(line, "COVERS SECTION "))
		case strings.HasPrefix(line, "STEP "):
			stepID := strings.TrimSpace(strings.TrimPrefix(line, "STEP "))
			text, ok := p.readTextBlock()
			if !ok {
				return proj, false
			}
			if closing, ok := p.peek(); !ok || closing != "END" {
				p.failHere("CDL_PARSE", "PROJECTION STEP %s missing closing END", stepID)
				return proj, false
			}
			p.pos++
			proj.StepTexts[stepID] = text
		case strings.HasPrefix(line, "RULE "):
			ruleID := strings.TrimSpace(strings.TrimPrefix(line, "RULE "))
			text, ok := p.readTextBlock()
			if !ok {
				return proj, false
			}
			if closing, ok := p.peek(); !ok || closing != "END" {
				p.failHere("CDL_PARSE", "PROJECTION RULE %s missing closing END", ruleID)
				return proj, false
			}
			p.pos++
			proj.RuleTexts[ruleID] = text
		case strings.HasPrefix(line, "OMIT STEP "):
			om := Omission{Step: strings.TrimSpace(strings.TrimPrefix(line, "OMIT STEP "))}
			done := false
			for !done {
				inner, ok := p.peek()
				if !ok {
					p.failHere("CDL_PARSE", "unterminated OMIT block")
					return proj, false
				}
				p.pos++
				switch {
				case inner == "END":
					proj.Omissions = append(proj.Omissions, om)
					done = true
				case inner == "":
				case strings.HasPrefix(inner, "CHANNELS "):
					om.Channels = csv(strings.TrimPrefix(inner, "CHANNELS "))
				case strings.HasPrefix(inner, "CHANNEL "):
					om.Channels = []string{strings.TrimSpace(strings.TrimPrefix(inner, "CHANNEL "))}
				case strings.HasPrefix(inner, "REASON "):
					om.Reason = unquote(strings.TrimPrefix(inner, "REASON "))
				case strings.HasPrefix(inner, "SUPPLIED-BY "):
					om.SuppliedBy = strings.TrimSpace(strings.TrimPrefix(inner, "SUPPLIED-BY "))
				default:
					p.failHere("CDL_PARSE", "unexpected OMIT line %q", inner)
					return proj, false
				}
			}
		default:
			p.failHere("CDL_PARSE", "unexpected projection line %q", line)
			return proj, false
		}
	}
}

func (p *parser) readTextBlock() (string, bool) {
	header, ok := p.peek()
	if !ok || header != "TEXT" {
		p.failHere("CDL_PARSE", "expected TEXT block")
		return "", false
	}
	p.pos++
	var parts []string
	for {
		line, ok := p.peek()
		if !ok {
			p.failHere("CDL_PARSE", "unterminated TEXT block")
			return "", false
		}
		p.pos++
		if line == "END" {
			break
		}
		if line != "" {
			parts = append(parts, line)
		}
	}
	text := strings.Join(parts, " ")
	if len(text) >= 2 && strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\"") {
		text = text[1 : len(text)-1]
	}
	return text, true
}
