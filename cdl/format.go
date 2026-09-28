package cdl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Format renders a canonical, deterministic CDL source from a parsed program.
// Opaque text blocks (GOAL, TEXT, REASON, REQUIRE) keep their content; layout
// is normalized to the canonical indentation, alignment, and declaration order.
func Format(prog *Program) []byte {
	f := &formatter{comments: map[string][]string{}}
	for _, c := range prog.Comments {
		f.comments[c.Before] = append(f.comments[c.Before], c.Text)
	}
	f.emitGlobals(prog)
	f.emitBaseReads(prog)
	for _, mode := range prog.Globals.Modes {
		f.flush("mode:" + mode.ID)
		f.line(0, "MODE %s", mode.ID)
		if mode.Title != "" {
			f.line(2, "TITLE %s", quote(mode.Title))
		}
		if mode.Strict {
			f.line(2, "STRICT")
		}
		if len(mode.Reads) > 0 {
			f.line(2, "READS %s", strings.Join(mode.Reads, ", "))
		}
		f.line(0, "END")
	}
	f.emitSections(prog)
	for _, proj := range prog.Projections {
		f.flush("projection:" + proj.ID)
		f.emitProjection(proj)
	}
	f.flush("")
	return []byte(f.b.String())
}

// FormatSource parses src and returns its canonical rendering. Sources with
// inline or nested comments are reported unsupported rather than reformatted.
func FormatSource(file string, src []byte) ([]byte, Diagnostics) {
	prog, diags := Parse(file, src)
	if len(diags) > 0 {
		return nil, diags
	}
	if prog.UnpreservedComments {
		return nil, Diagnostics{{
			Code: "CDL_FORMAT_UNSUPPORTED",
			File: file,
			Msg:  "source contains inline or nested comments; cdl fmt cannot preserve them",
		}}
	}
	return Format(prog), nil
}

type formatter struct {
	b        strings.Builder
	comments map[string][]string
}

func (f *formatter) flush(key string) {
	for _, text := range f.comments[key] {
		f.b.WriteString(text)
		f.b.WriteByte('\n')
	}
	delete(f.comments, key)
}

func (f *formatter) line(indent int, format string, args ...any) {
	if format == "" {
		f.b.WriteByte('\n')
		return
	}
	f.b.WriteString(strings.Repeat(" ", indent))
	fmt.Fprintf(&f.b, format, args...)
	f.b.WriteByte('\n')
}

func (f *formatter) blank() { f.b.WriteByte('\n') }

// opaque writes a text block body with a uniform indentation. The block's
// relative line structure and content are preserved exactly.
func (f *formatter) opaque(text string, indent int) {
	if text == "" {
		return
	}
	prefix := strings.Repeat(" ", indent)
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			f.b.WriteByte('\n')
			continue
		}
		f.b.WriteString(prefix)
		f.b.WriteString(line)
		f.b.WriteByte('\n')
	}
}

func (f *formatter) emitGlobals(prog *Program) {
	g := prog.Globals
	if len(g.Capabilities) == 0 && len(g.Registries) == 0 && len(g.Artifacts) == 0 &&
		len(g.Rules) == 0 && len(g.Gates) == 0 && len(g.WorkflowTargets) == 0 {
		return
	}
	f.flush("globals")
	f.line(0, "GLOBAL DECLARATIONS")
	for _, c := range g.Capabilities {
		f.line(2, "CAPABILITY %s", c.ID)
	}
	for _, id := range g.Registries {
		f.line(2, "REGISTRY %s", id)
	}
	for _, id := range g.Artifacts {
		f.line(2, "ARTIFACT %s", id)
	}
	for _, id := range g.Rules {
		f.line(2, "RULE %s", id)
	}
	if len(g.Gates) > 0 {
		f.line(2, "GATES %s", strings.Join(g.Gates, ", "))
	}
	if len(g.WorkflowTargets) > 0 {
		f.line(2, "WORKFLOW-TARGET %s", strings.Join(g.WorkflowTargets, ", "))
	}
	f.line(0, "END")
	f.blank()
}

func (f *formatter) emitBaseReads(prog *Program) {
	if len(prog.Globals.BaseReads) == 0 {
		return
	}
	f.flush("base-reads")
	f.line(0, "BASE-READS")
	for _, id := range prog.Globals.BaseReads {
		f.line(2, "%s", id)
	}
	f.line(0, "END")
	f.blank()
}

func uniquePartList(parts []PartDecl) []PartDecl {
	seen := map[string]bool{}
	var out []PartDecl
	for _, part := range parts {
		if seen[part.ID] {
			continue
		}
		seen[part.ID] = true
		out = append(out, part)
	}
	return out
}

func (f *formatter) emitSections(prog *Program) {
	parts := map[string]string{}
	for _, part := range uniquePartList(prog.Globals.Parts) {
		parts[part.ID] = part.Title
	}
	used := map[string]bool{}
	last := ""
	first := true
	for _, sec := range prog.Sections {
		if sec.Part != "" && sec.Part != last {
			if !first {
				f.blank()
			}
			f.flush("part:" + sec.Part)
			f.line(0, "PART %s %s", sec.Part, quote(parts[sec.Part]))
			f.blank()
			used[sec.Part] = true
			last = sec.Part
		}
		first = false
		f.flush("section:" + sec.ID)
		f.emitSection(sec, parts)
	}
	for _, part := range uniquePartList(prog.Globals.Parts) {
		if used[part.ID] {
			continue
		}
		f.flush("part:" + part.ID)
		f.line(0, "PART %s %s", part.ID, quote(part.Title))
		f.blank()
	}
}

func (f *formatter) emitSection(sec Section, parts map[string]string) {
	f.line(0, "SECTION %s", sec.ID)
	if sec.Parent != "" {
		f.line(0, "SUBSECTION OF %s", sec.Parent)
	}
	if sec.Title != "" {
		f.line(0, "TITLE %s", quote(sec.Title))
	}
	if sec.Artifact != "" {
		f.line(0, "ARTIFACT %s", sec.Artifact)
	}
	for _, use := range sec.Uses {
		f.line(0, "USES %s %s", use.Kind, use.ID)
	}
	if len(sec.Requires) > 0 {
		f.line(0, "REQUIRES CAPABILITY %s", strings.Join(sec.Requires, ", "))
	}
	if sec.Goal != "" {
		f.blank()
		f.line(0, "GOAL")
		f.opaque(sec.Goal, 0)
		f.line(0, "END")
		f.blank()
	}
	for _, t := range sec.Types {
		f.line(0, "TYPE %s := %s", t.ID, t.Type)
		f.blank()
	}
	for _, v := range sec.Values {
		f.line(0, "VALUE %s", v.ID)
		if v.Type != "" {
			f.line(2, "TYPE %s", v.Type)
		}
		if v.Lifetime != "" {
			f.line(2, "LIFETIME %s", v.Lifetime)
		}
		f.line(0, "END")
		f.blank()
	}
	for _, s := range sec.States {
		f.line(0, "STATE %s", s.ID)
		if len(s.Values) > 0 {
			f.line(2, "VALUES %s", strings.Join(s.Values, ", "))
		}
		if s.Initial != "" {
			f.line(2, "INITIAL %s", s.Initial)
		}
		if len(s.Allows) > 0 {
			f.line(2, "ALLOWS")
			for _, pair := range s.Allows {
				f.line(4, "%s -> %s", pair[0], pair[1])
			}
			f.line(2, "END")
		}
		f.line(0, "END")
		f.blank()
	}
	for _, e := range sec.Enums {
		f.line(0, "ENUM %s", e.ID)
		for _, value := range e.Values {
			f.line(2, "%s", value)
		}
		f.line(0, "END")
		f.blank()
	}
	for _, field := range sec.Fields {
		f.emitField(field)
	}
	for _, t := range sec.Tables {
		f.line(0, "TABLE %s", t.ID)
		if t.Rows != "" {
			f.line(2, "ROWS %s", t.Rows)
		}
		if t.RowType != "" {
			f.line(2, "ROW-TYPE %s", t.RowType)
		}
		if t.Key != "" {
			f.line(2, "KEY %s", t.Key)
		}
		f.line(0, "END")
		f.blank()
	}
	for _, r := range sec.Rules {
		f.emitRule(r)
	}
	for _, step := range sec.Steps {
		f.emitStep(step)
	}
}

func (f *formatter) emitField(field FieldDecl) {
	f.line(0, "FIELD %s", field.ID)
	nameWidth, typeWidth := 0, 0
	for _, spec := range field.Fields {
		if len(spec.Name) > nameWidth {
			nameWidth = len(spec.Name)
		}
		if len(spec.Type) > typeWidth {
			typeWidth = len(spec.Type)
		}
	}
	for _, spec := range field.Fields {
		req := "optional"
		if spec.Required {
			req = "required"
		}
		f.line(2, "%-*s -> %-*s %s", nameWidth, spec.Name, typeWidth, spec.Type, req)
	}
	f.line(0, "END")
	f.blank()
}

func (f *formatter) emitRule(rule RuleDecl) {
	f.line(0, "RULE %s", rule.ID)
	if rule.Goal != "" {
		f.line(2, "GOAL")
		f.opaque(rule.Goal, 4)
		f.line(2, "END")
	}
	if rule.Predicate != "" {
		f.line(2, "PREDICATE %s :=", rule.Predicate)
		for _, line := range wrapExpr(rule.Expr) {
			f.line(4, "%s", line)
		}
		f.line(0, "END")
		return
	}
	f.line(0, "END")
	f.blank()
}

func (f *formatter) emitStep(step Step) {
	f.line(0, "STEP %s", step.ID)
	if step.Owner == "MACHINE" {
		f.line(0, "MACHINE:")
		for _, id := range step.UsesRules {
			f.line(2, "USES RULE %s", id)
		}
		f.emitOps(step.Ops, 2)
		f.line(0, "END")
		f.blank()
		return
	}
	f.line(0, "AGENT:")
	emitList := func(key string, values []string) {
		if len(values) == 0 {
			return
		}
		f.line(2, "%s", key)
		for _, value := range values {
			f.line(4, "%s", value)
		}
		f.line(2, "END")
	}
	emitList("EVIDENCE", step.Evidence)
	emitList("PRODUCES", step.Produces)
	if step.Result != "" {
		f.line(2, "RESULT")
		f.line(4, "%s", step.Result)
		f.line(2, "END")
	}
	if step.Require != "" {
		f.line(2, `REQUIRE """`)
		f.line(4, "%s", step.Require)
		f.line(2, `"""`)
	}
	f.line(0, "END")
	f.blank()
}

func (f *formatter) emitOps(ops []Op, indent int) {
	for _, op := range ops {
		switch o := op.(type) {
		case RequiresOp:
			f.line(indent, "REQUIRES CAPABILITY %s", strings.Join(o.Caps, ", "))
		case SetOp:
			f.line(indent, "SET %s := %s", o.Target, exprString(o.Value))
		case LetOp:
			f.line(indent, "LET %s := %s", o.Target, exprString(o.Value))
		case AppendOp:
			f.line(indent, "APPEND %s FOR %s", o.RowType, o.For)
		case BlockOp:
			f.line(indent, "BLOCK %s", strings.Join(o.Gates, ", "))
		case ForEachOp:
			f.line(indent, "FOR EACH %s IN %s", o.Var, o.In)
			f.emitOps(o.Body, indent+2)
			f.line(indent, "END")
		case GuardOp:
			f.line(indent, "GUARD %s", exprString(o.Expr))
			if len(o.Otherwise) == 1 {
				if block, ok := o.Otherwise[0].(BlockOp); ok {
					f.line(indent, "OTHERWISE BLOCK %s", strings.Join(block.Gates, ", "))
					continue
				}
			}
			if len(o.Otherwise) > 0 {
				f.line(indent, "OTHERWISE")
				f.emitOps(o.Otherwise, indent+2)
				f.line(indent, "END")
			}
		case IfOp:
			f.line(indent, "IF %s", exprString(o.Cond))
			f.emitOps(o.Then, indent+2)
			if len(o.Else) > 0 {
				f.line(indent, "ELSE")
				f.emitOps(o.Else, indent+2)
			}
			f.line(indent, "END")
		}
	}
}

func (f *formatter) emitProjection(proj Projection) {
	f.line(0, "PROJECTION %s", proj.ID)
	if len(proj.Channels) > 0 {
		f.line(2, "CHANNELS %s", strings.Join(proj.Channels, ", "))
	}
	if proj.Covers != "" {
		f.line(2, "COVERS SECTION %s", proj.Covers)
	}
	stepIDs := sortedStringKeys(proj.StepTexts)
	for _, id := range stepIDs {
		f.blank()
		f.line(2, "STEP %s", id)
		f.line(4, "TEXT")
		f.line(6, "%s", quote(proj.StepTexts[id]))
		f.line(4, "END")
		f.line(2, "END")
	}
	ruleIDs := sortedStringKeys(proj.RuleTexts)
	for _, id := range ruleIDs {
		f.blank()
		f.line(2, "RULE %s", id)
		f.line(4, "TEXT")
		f.line(6, "%s", quote(proj.RuleTexts[id]))
		f.line(4, "END")
		f.line(2, "END")
	}
	for _, om := range proj.Omissions {
		f.blank()
		f.line(2, "OMIT STEP %s", om.Step)
		if len(om.Channels) > 0 {
			f.line(4, "CHANNELS %s", strings.Join(om.Channels, ", "))
		}
		if om.Reason != "" {
			f.line(4, "REASON %s", quote(om.Reason))
		}
		if om.SuppliedBy != "" {
			f.line(4, "SUPPLIED-BY %s", om.SuppliedBy)
		}
		f.line(2, "END")
	}
	f.line(0, "END")
	f.blank()
}

func sortedStringKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func quote(s string) string {
	if strings.ContainsAny(s, "\n#") {
		return strconv.Quote(s)
	}
	return `"` + s + `"`
}

const (
	precOr      = 1
	precAnd     = 2
	precNot     = 3
	precCompare = 4
	precAdd     = 5
	precPrimary = 6
)

func exprPrec(e Expr) int {
	switch v := e.(type) {
	case Binary:
		switch v.Op {
		case "OR":
			return precOr
		case "AND":
			return precAnd
		case "==", "!=", "<", "<=", ">", ">=":
			return precCompare
		case "+", "-":
			return precAdd
		}
		return precPrimary
	case Unary:
		return precNot
	default:
		return precPrimary
	}
}

func exprString(e Expr) string {
	if e == nil {
		return ""
	}
	switch v := e.(type) {
	case Ident:
		return v.Raw
	case Literal:
		switch val := v.Value.(type) {
		case int:
			return strconv.Itoa(val)
		case string:
			return `"` + val + `"`
		case bool:
			if val {
				return "true"
			}
			return "false"
		}
		return fmt.Sprint(v.Value)
	case Unary:
		operand := exprString(v.E)
		if exprPrec(v.E) < precNot {
			operand = "(" + operand + ")"
		}
		return v.Op + " " + operand
	case Binary:
		prec := exprPrec(v)
		left := exprString(v.L)
		if exprPrec(v.L) < prec {
			left = "(" + left + ")"
		}
		right := exprString(v.R)
		if exprPrec(v.R) <= prec {
			right = "(" + right + ")"
		}
		return left + " " + v.Op + " " + right
	case Call:
		if v.Name == "COUNT" && len(v.Args) > 0 {
			inner := exprString(v.Args[0])
			if v.Where != nil {
				inner += " WHERE " + v.Where.Field + " == " + exprString(v.Where.Value)
			}
			return "COUNT(" + inner + ")"
		}
		args := make([]string, 0, len(v.Args))
		for _, arg := range v.Args {
			args = append(args, exprString(arg))
		}
		return v.Name + "(" + strings.Join(args, ", ") + ")"
	}
	return ""
}

// wrapExpr renders an expression on one line; long predicates remain a single
// canonical line, which the parser accepts.
func wrapExpr(e Expr) []string {
	text := exprString(e)
	if text == "" {
		return nil
	}
	return []string{text}
}
