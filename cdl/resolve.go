package cdl

import (
	"fmt"
	"sort"
	"strings"
)

type symbols struct {
	caps           map[string]Capability
	registries     map[string]bool
	artifacts      map[string]bool
	rules          map[string]RuleDecl
	gates          map[string]bool
	workflows      map[string]bool
	types          map[string]string
	states         map[string]StateDecl
	values         map[string]ValueDecl
	enums          map[string][]string
	enumOwner      map[string]string
	fields         map[string]FieldDecl
	tables         map[string]TableDecl
	tableByRowType map[string]string
	sections       map[string]bool
	steps          map[string]bool
	projections    map[string]bool
	projectionByID map[string]Projection
}

func newSymbols() *symbols {
	return &symbols{
		caps:           map[string]Capability{},
		registries:     map[string]bool{},
		artifacts:      map[string]bool{},
		rules:          map[string]RuleDecl{},
		gates:          map[string]bool{},
		workflows:      map[string]bool{},
		types:          map[string]string{},
		states:         map[string]StateDecl{},
		values:         map[string]ValueDecl{},
		enums:          map[string][]string{},
		enumOwner:      map[string]string{},
		fields:         map[string]FieldDecl{},
		tables:         map[string]TableDecl{},
		tableByRowType: map[string]string{},
		sections:       map[string]bool{},
		steps:          map[string]bool{},
		projections:    map[string]bool{},
		projectionByID: map[string]Projection{},
	}
}

type resolver struct {
	prog  *Program
	sym   *symbols
	diags Diagnostics
}

func (r *resolver) fail(code, format string, args ...any) {
	r.diags = append(r.diags, Diagnostic{Code: code, Msg: fmt.Sprintf(format, args...)})
}

// resolve validates the merged program and returns its symbol table.
func resolve(prog *Program) (Diagnostics, *symbols) {
	r := &resolver{prog: prog, sym: newSymbols()}
	r.buildSymbols()
	r.checkSections()
	r.checkProjections()
	return r.diags, r.sym
}

func (r *resolver) buildSymbols() {
	add := func(id, kind string, exists bool) bool {
		if id == "" {
			r.fail("CDL_PARSE", "empty %s identifier", kind)
			return false
		}
		if exists {
			r.fail("CDL_DUPLICATE_ID", "duplicate %s %q", kind, id)
			return false
		}
		return true
	}
	for _, c := range r.prog.Globals.Capabilities {
		if add(c.ID, "capability", r.existsCap(c.ID)) {
			r.sym.caps[c.ID] = c
		}
	}
	for _, id := range r.prog.Globals.Registries {
		if add(id, "registry", r.sym.registries[id]) {
			r.sym.registries[id] = true
		}
	}
	for _, id := range r.prog.Globals.Artifacts {
		if add(id, "artifact", r.sym.artifacts[id]) {
			r.sym.artifacts[id] = true
		}
	}
	for _, id := range r.prog.Globals.Rules {
		if add(id, "rule", r.existsRule(id)) {
			r.sym.rules[id] = RuleDecl{ID: id}
		}
	}
	for _, id := range r.prog.Globals.Gates {
		if add(id, "gate", r.sym.gates[id]) {
			r.sym.gates[id] = true
		}
	}
	for _, id := range r.prog.Globals.WorkflowTargets {
		if add(id, "workflow target", r.sym.workflows[id]) {
			r.sym.workflows[id] = true
		}
	}
	for _, sec := range r.prog.Sections {
		if add(sec.ID, "section", r.sym.sections[sec.ID]) {
			r.sym.sections[sec.ID] = true
		}
		for _, t := range sec.Types {
			if add(t.ID, "type", r.sym.types[t.ID] != "") {
				r.sym.types[t.ID] = t.Type
			}
		}
		for _, s := range sec.States {
			if add(s.ID, "state", r.existsState(s.ID)) {
				r.sym.states[s.ID] = s
			}
		}
		for _, v := range sec.Values {
			if add(v.ID, "value", r.existsValue(v.ID)) {
				r.sym.values[v.ID] = v
			}
		}
		for _, e := range sec.Enums {
			if add(e.ID, "enum", r.existsEnum(e.ID)) {
				r.sym.enums[e.ID] = e.Values
				for _, val := range e.Values {
					r.sym.enumOwner[val] = e.ID
				}
			}
		}
		for _, f := range sec.Fields {
			if add(f.ID, "field", r.existsField(f.ID)) {
				r.sym.fields[f.ID] = f
			}
		}
		for _, t := range sec.Tables {
			if !add(t.ID, "table", r.existsTable(t.ID)) {
				continue
			}
			r.sym.tables[t.ID] = t
			r.sym.tableByRowType[t.RowType] = t.ID
		}
		for _, ru := range sec.Rules {
			if add(ru.ID, "rule", r.existsRule(ru.ID)) {
				r.sym.rules[ru.ID] = ru
			}
		}
		for _, st := range sec.Steps {
			if add(st.ID, "step", r.sym.steps[st.ID]) {
				r.sym.steps[st.ID] = true
			}
		}
	}
	for _, p := range r.prog.Projections {
		if add(p.ID, "projection", r.sym.projections[p.ID]) {
			r.sym.projections[p.ID] = true
			r.sym.projectionByID[p.ID] = p
		}
	}
}

func (r *resolver) existsCap(id string) bool  { _, ok := r.sym.caps[id]; return ok }
func (r *resolver) existsRule(id string) bool { _, ok := r.sym.rules[id]; return ok }
func (r *resolver) existsState(id string) bool {
	_, ok := r.sym.states[id]
	return ok
}
func (r *resolver) existsValue(id string) bool {
	_, ok := r.sym.values[id]
	return ok
}
func (r *resolver) existsEnum(id string) bool {
	_, ok := r.sym.enums[id]
	return ok
}
func (r *resolver) existsField(id string) bool {
	_, ok := r.sym.fields[id]
	return ok
}
func (r *resolver) existsTable(id string) bool {
	_, ok := r.sym.tables[id]
	return ok
}

func (r *resolver) checkSections() {
	for si := range r.prog.Sections {
		sec := &r.prog.Sections[si]
		r.checkSectionDeclarations(sec)
		r.checkSectionUses(sec)
		r.checkSectionSteps(sec)
	}
}

func (r *resolver) checkSectionDeclarations(sec *Section) {
	for _, t := range sec.Types {
		if !r.typeResolves(t.Type) {
			r.fail("CDL_UNRESOLVED_REFERENCE", "type %s refers to unknown type %q", t.ID, t.Type)
		}
	}
	for _, s := range sec.States {
		if len(s.Values) == 0 {
			r.fail("CDL_PARSE", "state %s declares no VALUES", s.ID)
			continue
		}
		if s.Initial != "" && !contains(s.Values, s.Initial) {
			r.fail("CDL_STATE_UNKNOWN_VALUE", "state %s initial value %q is not declared", s.ID, s.Initial)
		}
		for _, tr := range s.Allows {
			if !contains(s.Values, tr[0]) || !contains(s.Values, tr[1]) {
				r.fail("CDL_STATE_UNKNOWN_VALUE", "state %s transition %s -> %s uses an undeclared value", s.ID, tr[0], tr[1])
			}
		}
	}
	for _, v := range sec.Values {
		if !r.typeResolves(v.Type) {
			r.fail("CDL_UNRESOLVED_REFERENCE", "value %s refers to unknown type %q", v.ID, v.Type)
		}
		switch v.Lifetime {
		case "step", "section", "run", "checkpoint":
		default:
			r.fail("CDL_PARSE", "value %s has invalid LIFETIME %q", v.ID, v.Lifetime)
		}
	}
	for _, f := range sec.Fields {
		for _, spec := range f.Fields {
			if !r.typeResolves(spec.Type) {
				r.fail("CDL_UNRESOLVED_REFERENCE", "field %s.%s refers to unknown type %q", f.ID, spec.Name, spec.Type)
			}
		}
	}
	for _, t := range sec.Tables {
		if !r.sym.registries[t.Rows] {
			r.fail("CDL_UNRESOLVED_REFERENCE", "table %s ROWS unknown registry %q", t.ID, t.Rows)
		}
		fd, ok := r.sym.fields[t.RowType]
		if !ok {
			r.fail("CDL_UNRESOLVED_REFERENCE", "table %s ROW-TYPE unknown field %q", t.ID, t.RowType)
			continue
		}
		if !fieldDeclared(fd, t.Key) {
			r.fail("CDL_UNRESOLVED_REFERENCE", "table %s KEY %q is not a %s field", t.ID, t.Key, t.RowType)
		}
	}
	for _, ru := range sec.Rules {
		if ru.Expr == nil {
			continue
		}
		if ru.Predicate == "" {
			r.fail("CDL_PARSE", "rule %s has an expression but no PREDICATE name", ru.ID)
			continue
		}
		if typ := r.inferExpr(ru.Expr, nil); typ != "bool" {
			r.fail("CDL_GUARD_NOT_BOOLEAN", "rule %s predicate %s must be boolean", ru.ID, ru.Predicate)
		}
	}
}

func (r *resolver) checkSectionUses(sec *Section) {
	for _, u := range sec.Uses {
		switch u.Kind {
		case "REGISTRY":
			if !r.sym.registries[u.ID] {
				r.fail("CDL_UNRESOLVED_REFERENCE", "USES REGISTRY %s does not resolve", u.ID)
			}
		case "RULE":
			if !r.existsRule(u.ID) {
				r.fail("CDL_UNRESOLVED_REFERENCE", "USES RULE %s does not resolve", u.ID)
			}
		case "STATE":
			if !r.existsState(u.ID) {
				r.fail("CDL_UNRESOLVED_REFERENCE", "USES STATE %s does not resolve", u.ID)
			}
		default:
			r.fail("CDL_PARSE", "USES kind %q is not recognized", u.Kind)
		}
	}
}

func (r *resolver) checkSectionSteps(sec *Section) {
	union := map[string]bool{}
	declared := map[string]bool{}
	for _, c := range sec.Requires {
		declared[c] = true
	}
	for si := range sec.Steps {
		step := &sec.Steps[si]
		switch step.Owner {
		case "MACHINE", "AGENT":
		default:
			r.fail("CDL_PARSE", "step %s has invalid owner %q", step.ID, step.Owner)
		}
		for _, id := range step.UsesRules {
			if !r.existsRule(id) {
				r.fail("CDL_UNRESOLVED_REFERENCE", "step %s USES RULE %s does not resolve", step.ID, id)
			}
		}
		step.Ops = r.resolveOps(step, step.Ops)
		step.Requires = collectRequires(step.Ops)
		for _, c := range step.Requires {
			union[c] = true
		}
		r.checkAgentStep(sec, step)
	}
	if !sameSet(union, declared) {
		r.fail("CDL_REQUIRES_SUMMARY_MISMATCH",
			"section %s REQUIRES %v does not equal the union of step requirements %v",
			sec.ID, sortedKeys(declared), sortedKeys(union))
	}
}

func (r *resolver) checkAgentStep(sec *Section, step *Step) {
	if step.Owner != "AGENT" {
		return
	}
	for _, ev := range step.Evidence {
		if !r.sym.artifacts[ev] && !r.sym.registries[ev] {
			r.fail("CDL_UNRESOLVED_REFERENCE", "step %s EVIDENCE %q does not resolve", step.ID, ev)
		}
	}
	if step.Result != "" {
		if !r.existsEnum(step.Result) {
			r.fail("CDL_UNRESOLVED_REFERENCE", "step %s RESULT %q does not resolve", step.ID, step.Result)
		}
	}
	for _, p := range step.Produces {
		parts := strings.SplitN(p, ".", 2)
		if len(parts) != 2 {
			r.fail("CDL_PARSE", "step %s PRODUCES %q is not <table>.<field>", step.ID, p)
			continue
		}
		table, ok := r.sym.tables[parts[0]]
		if !ok {
			r.fail("CDL_UNRESOLVED_REFERENCE", "step %s PRODUCES unknown table %q", step.ID, parts[0])
			continue
		}
		fd, ok := r.sym.fields[table.RowType]
		if !ok || !fieldDeclared(fd, parts[1]) {
			r.fail("CDL_UNRESOLVED_REFERENCE", "step %s PRODUCES unknown field %q on %s", step.ID, parts[1], table.RowType)
		}
	}
}

func (r *resolver) resolveOps(step *Step, ops []Op) []Op {
	out := make([]Op, 0, len(ops))
	for _, op := range ops {
		switch v := op.(type) {
		case RequiresOp:
			r.resolveCapList(step, v.Caps)
			out = append(out, v)
		case SetOp:
			r.resolveSet(step, v.Target, v.Value)
			out = append(out, v)
		case LetOp:
			r.inferExpr(v.Value, step)
			out = append(out, v)
		case AppendOp:
			table, ok := r.sym.tableByRowType[v.RowType]
			if !ok {
				r.fail("CDL_UNRESOLVED_REFERENCE", "step %s APPEND refers to unknown row type %q", step.ID, v.RowType)
				out = append(out, v)
				continue
			}
			v.To = table
			out = append(out, v)
		case ForEachOp:
			if !r.sym.registries[v.In] {
				r.fail("CDL_UNRESOLVED_REFERENCE", "step %s FOR EACH refers to unknown registry %q", step.ID, v.In)
			}
			v.Body = r.resolveOps(step, v.Body)
			out = append(out, v)
		case GuardOp:
			if typ := r.inferExpr(v.Expr, step); typ != "bool" {
				r.fail("CDL_GUARD_NOT_BOOLEAN", "step %s GUARD must be boolean (use an explicit comparison)", step.ID)
			}
			v.Otherwise = r.resolveOps(step, v.Otherwise)
			out = append(out, v)
		case IfOp:
			if typ := r.inferExpr(v.Cond, step); typ != "bool" {
				r.fail("CDL_GUARD_NOT_BOOLEAN", "step %s IF condition must be boolean (use an explicit comparison)", step.ID)
			}
			v.Then = r.resolveOps(step, v.Then)
			v.Else = r.resolveOps(step, v.Else)
			out = append(out, v)
		case BlockOp:
			r.resolveGates(step, v.Gates)
			out = append(out, v)
		default:
			r.fail("CDL_UNKNOWN_INSTRUCTION", "step %s contains an unsupported operation", step.ID)
			out = append(out, op)
		}
	}
	return out
}

// collectRequires returns the ordered, de-duplicated operation-level capability
// requirements of a step body.
func collectRequires(ops []Op) []string {
	var out []string
	seen := map[string]bool{}
	var walk func([]Op)
	walk = func(list []Op) {
		for _, op := range list {
			switch v := op.(type) {
			case RequiresOp:
				for _, c := range v.Caps {
					if !seen[c] {
						seen[c] = true
						out = append(out, c)
					}
				}
			case ForEachOp:
				walk(v.Body)
			case GuardOp:
				walk(v.Otherwise)
			case IfOp:
				walk(v.Then)
				walk(v.Else)
			}
		}
	}
	walk(ops)
	return out
}

func (r *resolver) resolveCapList(step *Step, caps []string) {
	for _, c := range caps {
		cap, ok := r.sym.caps[c]
		if !ok {
			r.fail("CDL_UNRESOLVED_REFERENCE", "step %s requires unknown capability %q", step.ID, c)
			continue
		}
		if cap.Kind != "deterministic" {
			r.fail("CDL_CAPABILITY_KIND", "step %s is MACHINE but capability %s is %s", step.ID, c, cap.Kind)
		}
	}
}

func (r *resolver) resolveGates(step *Step, gates []string) {
	for _, g := range gates {
		if !r.sym.gates[g] {
			r.fail("CDL_BLOCK_TARGET_UNRESOLVED", "step %s BLOCK target %q does not resolve to a declared gate", step.ID, g)
		}
	}
}

func (r *resolver) resolveSet(step *Step, target string, value Expr) {
	if st, ok := r.sym.states[target]; ok {
		if id, isIdent := value.(Ident); isIdent && contains(st.Values, id.Raw) {
			r.checkStateAssignment(step, st, []string{id.Raw})
			return
		}
		if vals := r.possibleValues(value); vals != nil {
			r.inferExpr(value, step) // reference validation only
			r.checkStateAssignment(step, st, vals)
			return
		}
		typ := r.inferExpr(value, step)
		if !stateAcceptsType(st, typ) {
			r.fail("CDL_STATE_UNKNOWN_VALUE", "step %s SET %s value type %s is not in state VALUES", step.ID, target, typ)
		}
		return
	}
	if vd, ok := r.sym.values[target]; ok {
		typ := r.inferExpr(value, step)
		if typ != vd.Type {
			r.fail("CDL_SET_TYPE_MISMATCH", "step %s SET %s expects %s but expression is %s", step.ID, target, vd.Type, typ)
		}
		return
	}
	r.fail("CDL_SET_TARGET_INVALID", "step %s SET target %q is not a declared STATE or VALUE", step.ID, target)
}

// inferExpr returns a coarse type for an expression and records unresolved
// references. Types: int, bool, string, hash, state:<id>, value:<type>,
// enum:<id>, rule:<id>, unknown.
func (r *resolver) inferExpr(e Expr, step *Step) string {
	switch v := e.(type) {
	case nil:
		return "unknown"
	case Literal:
		switch v.Value.(type) {
		case bool:
			return "bool"
		case int:
			return "int"
		case string:
			return "string"
		}
		return "unknown"
	case Ident:
		if id, pred, ok := splitRuleRef(v.Raw); ok {
			rule, exists := r.sym.rules[id]
			if !exists || rule.Predicate != pred {
				r.failStep(step, "CDL_UNRESOLVED_REFERENCE", "reference %q does not resolve to a rule predicate", v.Raw)
				return "unknown"
			}
			return "bool"
		}
		if st, ok := r.sym.states[v.Raw]; ok {
			if contains(st.Values, "unknown") {
				return "state:" + v.Raw
			}
			return "state:" + v.Raw
		}
		if vd, ok := r.sym.values[v.Raw]; ok {
			return "value:" + vd.Type
		}
		if _, ok := r.sym.enumOwner[v.Raw]; ok {
			return "enum:" + r.sym.enumOwner[v.Raw]
		}
		r.failStep(step, "CDL_UNRESOLVED_REFERENCE", "expression reference %q does not resolve", v.Raw)
		return "unknown"
	case Binary:
		switch v.Op {
		case ">=", "<=", ">", "<":
			r.requireType(v.L, "int", step)
			r.requireType(v.R, "int", step)
			return "bool"
		case "+", "-":
			r.requireType(v.L, "int", step)
			r.requireType(v.R, "int", step)
			return "int"
		case "==", "!=":
			r.inferExpr(v.L, step)
			r.inferExpr(v.R, step)
			return "bool"
		case "AND", "OR":
			r.requireType(v.L, "bool", step)
			r.requireType(v.R, "bool", step)
			return "bool"
		}
		return "unknown"
	case Unary:
		r.requireType(v.E, "bool", step)
		return "bool"
	case Call:
		switch v.Name {
		case "COUNT":
			if v.Where != nil {
				r.inferExpr(v.Where.Value, step)
			}
			return "int"
		case "HASH":
			return "hash"
		case "AVAILABLE", "UNAVAILABLE", "ALL", "ANY":
			return "bool"
		case "NEXT":
			return "string"
		}
		for _, a := range v.Args {
			r.inferExpr(a, step)
		}
		return "unknown"
	}
	return "unknown"
}

func (r *resolver) requireType(e Expr, want string, step *Step) {
	got := r.inferExpr(e, step)
	if got != want {
		r.failStep(step, "CDL_GUARD_NOT_BOOLEAN", "expression has type %s where %s is required", got, want)
	}
}

func (r *resolver) failStep(step *Step, code, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if step != nil {
		msg = "step " + step.ID + ": " + msg
	}
	r.fail(code, "%s", msg)
}

func (r *resolver) checkStateAssignment(step *Step, st StateDecl, vals []string) {
	for _, v := range vals {
		if !contains(st.Values, v) {
			r.fail("CDL_STATE_UNKNOWN_VALUE", "step %s SET %s assigns undeclared value %q", step.ID, st.ID, v)
			continue
		}
		for _, from := range st.Values {
			if from == v {
				continue
			}
			if !allowsTransition(st, from, v) {
				r.fail("CDL_ALLOWS_VIOLATION", "step %s SET %s transition %s -> %s is not in ALLOWS", step.ID, st.ID, from, v)
			}
		}
	}
}

func (r *resolver) possibleValues(e Expr) []string {
	switch v := e.(type) {
	case Literal:
		switch val := v.Value.(type) {
		case bool:
			if val {
				return []string{"true"}
			}
			return []string{"false"}
		case string:
			return []string{val}
		}
		return nil
	case Ident:
		if _, ok := r.sym.enumOwner[v.Raw]; ok {
			return []string{v.Raw}
		}
		if _, _, ok := splitRuleRef(v.Raw); ok {
			return []string{"true", "false"}
		}
		return nil
	case Binary:
		if v.Op == "==" || v.Op == "!=" || v.Op == ">" || v.Op == ">=" || v.Op == "<" || v.Op == "<=" ||
			v.Op == "AND" || v.Op == "OR" {
			return []string{"true", "false"}
		}
		return nil
	case Unary:
		return []string{"true", "false"}
	case Call:
		switch v.Name {
		case "AVAILABLE", "UNAVAILABLE", "ALL", "ANY":
			return []string{"true", "false"}
		}
		return nil
	}
	return nil
}

func (r *resolver) typeResolves(t string) bool {
	switch t {
	case "string", "path", "id", "hash", "int", "bool":
		return true
	}
	if strings.HasPrefix(t, "set<") && strings.HasSuffix(t, ">") {
		inner := strings.TrimSuffix(strings.TrimPrefix(t, "set<"), ">")
		if _, ok := r.sym.types[inner]; ok {
			return true
		}
		if _, ok := r.sym.enums[inner]; ok {
			return true
		}
		switch inner {
		case "string", "path", "id", "hash", "int", "bool":
			return true
		}
	}
	if _, ok := r.sym.types[t]; ok {
		return true
	}
	_, ok := r.sym.enums[t]
	return ok
}

func (r *resolver) checkProjections() {
	for _, p := range r.prog.Projections {
		for _, ch := range p.Channels {
			if ch != "prompt" && ch != "native" {
				r.fail("CDL_PARSE", "projection %s has unknown channel %q", p.ID, ch)
			}
		}
		if p.Covers != "" && !r.sym.sections[p.Covers] {
			r.fail("CDL_UNRESOLVED_REFERENCE", "projection %s COVERS unknown section %q", p.ID, p.Covers)
		}
		for stepID := range p.StepTexts {
			if !r.sym.steps[stepID] {
				r.fail("CDL_UNRESOLVED_REFERENCE", "projection %s overrides unknown step %q", p.ID, stepID)
			}
		}
		for ruleID := range p.RuleTexts {
			if !r.existsRule(ruleID) {
				r.fail("CDL_UNRESOLVED_REFERENCE", "projection %s overrides unknown rule %q", p.ID, ruleID)
			}
		}
		channels := p.Channels
		if len(channels) == 0 {
			channels = []string{"prompt"}
		}
		for _, om := range p.Omissions {
			step, ok := r.findStep(om.Step)
			if !ok {
				r.fail("CDL_UNRESOLVED_REFERENCE", "projection %s omits unknown step %q", p.ID, om.Step)
				continue
			}
			omChannels := om.Channels
			if len(omChannels) == 0 {
				omChannels = channels
			}
			if om.SuppliedBy == "" {
				r.fail("CDL_OMISSION_MISSING_SUPPLIED_BY", "projection %s omission of %s has no SUPPLIED-BY", p.ID, om.Step)
				continue
			}
			switch {
			case om.SuppliedBy == "native-harness":
				for _, ch := range omChannels {
					if ch == "prompt" {
						r.fail("CDL_OMISSION_NATIVE_IN_PROMPT", "projection %s omits %s for prompt without an available supply", p.ID, om.Step)
					}
				}
			case strings.HasPrefix(om.SuppliedBy, "RULE "):
				ruleID := strings.TrimSpace(strings.TrimPrefix(om.SuppliedBy, "RULE "))
				rule, exists := r.sym.rules[ruleID]
				if !exists || rule.Predicate == "" {
					r.fail("CDL_OMISSION_SUPPLIED_BY_UNRESOLVED", "projection %s omission of %s cites rule %q that does not define a predicate", p.ID, om.Step, ruleID)
					continue
				}
				if hasSetOp(step.Ops) {
					r.fail("CDL_OMISSION_RULE_SUPPLIES_EFFECT", "projection %s omits %s which performs a SET effect; a rule cannot supply that", p.ID, om.Step)
				}
			case strings.HasPrefix(om.SuppliedBy, "CAPABILITY "):
				capID := strings.TrimSpace(strings.TrimPrefix(om.SuppliedBy, "CAPABILITY "))
				if !r.existsCap(capID) {
					r.fail("CDL_OMISSION_SUPPLIED_BY_UNRESOLVED", "projection %s omission of %s cites unknown capability %q", p.ID, om.Step, capID)
				}
			default:
				r.fail("CDL_OMISSION_SUPPLIED_BY_UNRESOLVED", "projection %s omission of %s has unrecognized SUPPLIED-BY %q", p.ID, om.Step, om.SuppliedBy)
			}
		}
	}
}

func (r *resolver) findStep(id string) (Step, bool) {
	for _, sec := range r.prog.Sections {
		for _, step := range sec.Steps {
			if step.ID == id {
				return step, true
			}
		}
	}
	return Step{}, false
}

func splitRuleRef(raw string) (string, string, bool) {
	idx := strings.LastIndex(raw, ".")
	if idx <= 0 || idx == len(raw)-1 {
		return "", "", false
	}
	return raw[:idx], raw[idx+1:], true
}

func stateAcceptsType(st StateDecl, typ string) bool {
	switch typ {
	case "bool":
		return contains(st.Values, "true") && contains(st.Values, "false")
	case "int", "string", "hash":
		return false
	}
	return false
}

func allowsTransition(st StateDecl, from, to string) bool {
	for _, tr := range st.Allows {
		if tr[0] == from && tr[1] == to {
			return true
		}
	}
	return false
}

func contains(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

func fieldDeclared(fd FieldDecl, name string) bool {
	for _, f := range fd.Fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

func hasSetOp(ops []Op) bool {
	for _, op := range ops {
		switch v := op.(type) {
		case SetOp:
			return true
		case ForEachOp:
			if hasSetOp(v.Body) {
				return true
			}
		case GuardOp:
			if hasSetOp(v.Otherwise) {
				return true
			}
		case IfOp:
			if hasSetOp(v.Then) || hasSetOp(v.Else) {
				return true
			}
		}
	}
	return false
}

func sameSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// collectIDs returns every stable identity in the program for ledger checking.
func collectIDs(prog *Program) []IDRecord {
	var out []IDRecord
	add := func(id, kind string) {
		if id != "" {
			out = append(out, IDRecord{ID: id, Kind: kind, Status: "active"})
		}
	}
	for _, c := range prog.Globals.Capabilities {
		add(c.ID, "CAPABILITY")
	}
	for _, id := range prog.Globals.Registries {
		add(id, "REGISTRY")
	}
	for _, id := range prog.Globals.Artifacts {
		add(id, "ARTIFACT")
	}
	for _, id := range prog.Globals.Rules {
		add(id, "RULE")
	}
	for _, id := range prog.Globals.Gates {
		add(id, "GATE")
	}
	for _, id := range prog.Globals.WorkflowTargets {
		add(id, "WORKFLOW-TARGET")
	}
	for _, sec := range prog.Sections {
		add(sec.ID, "SECTION")
		for _, t := range sec.Types {
			add(t.ID, "TYPE")
		}
		for _, s := range sec.States {
			add(s.ID, "STATE")
		}
		for _, v := range sec.Values {
			add(v.ID, "VALUE")
		}
		for _, e := range sec.Enums {
			add(e.ID, "ENUM")
		}
		for _, f := range sec.Fields {
			add(f.ID, "FIELD")
		}
		for _, t := range sec.Tables {
			add(t.ID, "TABLE")
		}
		for _, ru := range sec.Rules {
			add(ru.ID, "RULE")
		}
		for _, st := range sec.Steps {
			add(st.ID, "STEP")
		}
	}
	for _, p := range prog.Projections {
		add(p.ID, "PROJECTION")
	}
	return out
}

// checkLedger verifies the committed identity ledger.
func checkLedger(ids []IDRecord, ledger map[string]IDRecord) Diagnostics {
	var diags Diagnostics
	if ledger == nil {
		return Diagnostics{{Code: "CDL_ID_LEDGER_MISSING", Msg: "no committed identity ledger supplied"}}
	}
	for _, id := range ids {
		rec, ok := ledger[id.ID]
		if !ok {
			diags = append(diags, Diagnostic{Code: "CDL_ID_LEDGER_UNKNOWN_ID", Msg: "identity " + id.ID + " is not in the committed ledger"})
			continue
		}
		if rec.Status != "active" {
			diags = append(diags, Diagnostic{Code: "CDL_ID_LEDGER_RETIRED", Msg: "identity " + id.ID + " is " + rec.Status + " in the committed ledger"})
			continue
		}
		if rec.Kind != id.Kind {
			diags = append(diags, Diagnostic{Code: "CDL_ID_LEDGER_KIND", Msg: "identity " + id.ID + " is " + rec.Kind + " in the ledger but " + id.Kind + " in source"})
		}
	}
	return diags
}
