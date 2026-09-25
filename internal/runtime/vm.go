package runtime

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

type tableState struct {
	decl cdl.EIRTable
	rows []map[string]any
}

type stage struct {
	states map[string]any
	values map[string]any
	tables map[string]*tableState
	gates  map[string]string
}

type machine struct {
	doc       *cdl.EIRDoc
	fx        Fixture
	caps      map[string]bool
	committed *stage
	work      *stage
	effects   []Effect
	guardDiag *Diagnostic
	bindings  map[string]map[string]any
	trace     []StepRecord
}

// Execute runs every section in order through the native MACHINE VM.
func Execute(doc *cdl.EIRDoc, fx Fixture, caps map[string]bool) (Trace, error) {
	m := &machine{doc: doc, fx: fx, caps: caps, bindings: map[string]map[string]any{}}
	m.committed = m.initialStage()
	for _, sec := range doc.Sections {
		for _, step := range sec.Steps {
			if err := m.execStep(step); err != nil {
				return Trace{}, err
			}
		}
	}
	return Trace{Steps: m.trace}, nil
}

func (m *machine) initialStage() *stage {
	st := &stage{states: map[string]any{}, values: map[string]any{}, tables: map[string]*tableState{}, gates: map[string]string{}}
	for _, d := range m.doc.Declarations.States {
		if d.Initial != "" {
			st.states[d.ID] = stateValue(d.Initial)
		}
	}
	for _, d := range m.doc.Declarations.Tables {
		st.tables[d.ID] = &tableState{decl: d}
	}
	for _, g := range m.doc.Declarations.Gates {
		st.gates[g] = "open"
	}
	return st
}

func (m *machine) cloneStage(src *stage) *stage {
	out := &stage{states: map[string]any{}, values: map[string]any{}, tables: map[string]*tableState{}, gates: map[string]string{}}
	for k, v := range src.states {
		out.states[k] = v
	}
	for k, v := range src.values {
		out.values[k] = v
	}
	for k, v := range src.gates {
		out.gates[k] = v
	}
	for k, t := range src.tables {
		clone := &tableState{decl: t.decl, rows: make([]map[string]any, 0, len(t.rows))}
		for _, row := range t.rows {
			copied := map[string]any{}
			for f, v := range row {
				copied[f] = v
			}
			clone.rows = append(clone.rows, copied)
		}
		out.tables[k] = clone
	}
	return out
}

func (m *machine) execStep(step cdl.EIRStep) error {
	if step.Owner == "AGENT" {
		return m.execAgent(step)
	}
	return m.execMachine(step)
}

func (m *machine) execMachine(step cdl.EIRStep) error {
	rec := StepRecord{
		Step: step.ID, Owner: step.Owner, EIRHash: m.doc.Envelope.EIRHash, Channel: "native",
		Effects: []Effect{}, Diagnostics: []Diagnostic{},
	}
	var missing []string
	for _, c := range step.Requires {
		if !m.caps[c] && !handlesAvailability(step.Ops, c) {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		for _, c := range missing {
			rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CAPABILITY_UNAVAILABLE", Subject: c})
		}
		rec.Outcome = "unsupported"
		rec.Authority = "none"
		m.trace = append(m.trace, rec)
		return nil
	}
	m.work = m.cloneStage(m.committed)
	m.effects = nil
	m.guardDiag = nil
	outcome, err := m.execOps(step.Ops)
	if err != nil {
		rec.Outcome = "failed"
		rec.Authority = "none"
		rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "EIR_EVAL_ERROR", Subject: err.Error()})
		m.trace = append(m.trace, rec)
		return nil
	}
	if outcome == "" {
		outcome = "committed"
	}
	m.committed = m.work
	rec.Outcome = outcome
	rec.Authority = "authoritative"
	rec.Effects = m.effects
	if rec.Effects == nil {
		rec.Effects = []Effect{}
	}
	if m.guardDiag != nil && outcome == "blocked" {
		rec.Diagnostics = append(rec.Diagnostics, *m.guardDiag)
	}
	m.trace = append(m.trace, rec)
	return nil
}

func (m *machine) execOps(ops []map[string]any) (string, error) {
	for _, op := range ops {
		switch {
		case op["requires"] != nil:
		case op["set"] != nil:
			if err := m.execSet(asMap(op["set"])); err != nil {
				return "", err
			}
		case op["let"] != nil:
			return "", fmt.Errorf("LET is not implemented by the pilot VM")
		case op["append"] != nil:
			if err := m.execAppend(asMap(op["append"])); err != nil {
				return "", err
			}
		case op["foreach"] != nil:
			if err := m.execForEach(asMap(op["foreach"])); err != nil {
				return "", err
			}
		case op["guard"] != nil:
			g := asMap(op["guard"])
			val, err := m.eval(g["expr"])
			if err != nil {
				return "", err
			}
			if val != true {
				m.work = m.cloneStage(m.committed)
				m.effects = nil
				m.guardDiag = &Diagnostic{Code: "GUARD_BLOCKED", Subject: exprSubject(g["expr"])}
				if _, err := m.execOps(asMaps(g["otherwise"])); err != nil {
					return "", err
				}
				return "blocked", nil
			}
		case op["if"] != nil:
			cond := asMap(op["if"])
			val, err := m.eval(cond["cond"])
			if err != nil {
				return "", err
			}
			branch := cond["else"]
			if val == true {
				branch = cond["then"]
			}
			if _, err := m.execOps(asMaps(branch)); err != nil {
				return "", err
			}
		case op["block"] != nil:
			b := asMap(op["block"])
			for _, g := range asStrings(b["gates"]) {
				m.work.gates[g] = "blocked"
				m.effects = append(m.effects, Effect{Target: g, Operation: "block", Value: "blocked"})
			}
		default:
			return "", fmt.Errorf("unsupported operation form")
		}
	}
	return "", nil
}

func (m *machine) execSet(set map[string]any) error {
	target := fmt.Sprint(set["target"])
	value, err := m.eval(set["value"])
	if err != nil {
		return err
	}
	if current, ok := m.lookup(target); ok && equal(current, value) {
		return nil
	}
	if _, ok := m.work.states[target]; ok {
		m.work.states[target] = value
	} else {
		m.work.values[target] = value
	}
	m.effects = append(m.effects, Effect{Target: target, Operation: "set", Value: valueString(value)})
	return nil
}

func (m *machine) lookup(name string) (any, bool) {
	if v, ok := m.work.states[name]; ok {
		return v, true
	}
	v, ok := m.work.values[name]
	return v, ok
}

func (m *machine) execAppend(appendOp map[string]any) error {
	tableID := fmt.Sprint(appendOp["to"])
	tbl, ok := m.work.tables[tableID]
	if !ok {
		return fmt.Errorf("append targets unknown table %s", tableID)
	}
	unit := m.bindings[fmt.Sprint(appendOp["for"])]
	if unit == nil {
		return fmt.Errorf("append references unbound variable %v", appendOp["for"])
	}
	row := map[string]any{}
	for _, name := range m.fieldNames(tbl.decl.RowType) {
		if v, ok := unit[name]; ok {
			row[name] = v
		}
	}
	if _, ok := row[tbl.decl.Key]; !ok {
		return fmt.Errorf("appended row is missing key field %s", tbl.decl.Key)
	}
	tbl.rows = append(tbl.rows, row)
	m.effects = append(m.effects, Effect{Target: tableID, Operation: "append", Value: fmt.Sprint(row[tbl.decl.Key])})
	return nil
}

func (m *machine) execForEach(loop map[string]any) error {
	registry := fmt.Sprint(loop["in"])
	records, ok := m.fx.Registries[registry]
	if !ok {
		return fmt.Errorf("registry %s has no fixture data", registry)
	}
	name := fmt.Sprint(loop["var"])
	for _, rec := range records {
		m.bindings[name] = rec
		if _, err := m.execOps(asMaps(loop["do"])); err != nil {
			return err
		}
	}
	delete(m.bindings, name)
	return nil
}

func (m *machine) execAgent(step cdl.EIRStep) error {
	rec := StepRecord{
		Step: step.ID, Owner: step.Owner, EIRHash: m.doc.Envelope.EIRHash, Channel: "native",
		Effects: []Effect{}, Diagnostics: []Diagnostic{},
	}
	work := m.cloneStage(m.committed)
	data := m.fx.Agent[step.ID]
	var effects []Effect
	var failure *Diagnostic
	for _, produced := range step.Produces {
		parts := strings.SplitN(produced, ".", 2)
		if len(parts) != 2 {
			failure = &Diagnostic{Code: "CONTRACT_INVALID", Subject: "malformed produces entry " + produced}
			break
		}
		tbl, ok := work.tables[parts[0]]
		if !ok {
			failure = &Diagnostic{Code: "CONTRACT_INVALID", Subject: "unknown table " + parts[0]}
			break
		}
		field := parts[1]
		fieldType := m.fieldType(tbl.decl.RowType, field)
		enum, isEnum := m.enumValues(fieldType)
		rows := tbl.rows
		for _, row := range rows {
			key := fmt.Sprint(row[tbl.decl.Key])
			entry, ok := data[key]
			if !ok {
				failure = &Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " missing result for " + key}
				break
			}
			value, ok := entry[field]
			if !ok {
				failure = &Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " missing field for " + key}
				break
			}
			if isEnum && !contains(enum, fmt.Sprint(value)) {
				failure = &Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " value outside RESULT domain for " + key}
				break
			}
			row[field] = value
			effects = append(effects, Effect{Target: key + "." + field, Operation: "draft-set", Value: fmt.Sprint(value)})
		}
		if failure != nil {
			break
		}
	}
	if failure != nil {
		rec.Outcome = "failed"
		rec.Authority = "none"
		rec.Diagnostics = append(rec.Diagnostics, *failure)
		m.trace = append(m.trace, rec)
		return nil
	}
	m.committed = work
	rec.Outcome = "committed"
	rec.Authority = "draft"
	rec.Effects = effects
	m.trace = append(m.trace, rec)
	return nil
}

func (m *machine) fieldNames(rowType string) []string {
	for _, f := range m.doc.Declarations.Fields {
		if f.ID == rowType {
			out := make([]string, 0, len(f.Fields))
			for _, spec := range f.Fields {
				out = append(out, spec.Name)
			}
			return out
		}
	}
	return nil
}

func (m *machine) fieldType(rowType, name string) string {
	for _, f := range m.doc.Declarations.Fields {
		if f.ID == rowType {
			for _, spec := range f.Fields {
				if spec.Name == name {
					return spec.Type
				}
			}
		}
	}
	return ""
}

func (m *machine) enumValues(id string) ([]string, bool) {
	for _, e := range m.doc.Declarations.Enums {
		if e.ID == id {
			return e.Values, true
		}
	}
	return nil, false
}

// eval evaluates an EIR expression against the staged view.
func (m *machine) eval(e any) (any, error) {
	switch v := e.(type) {
	case bool, string:
		return v, nil
	case float64:
		if v == math.Trunc(v) {
			return int(v), nil
		}
		return v, nil
	case nil:
		return nil, fmt.Errorf("nil expression")
	case map[string]any:
		switch {
		case v["ref"] != nil:
			return m.evalRef(fmt.Sprint(v["ref"]))
		case v["op"] != nil:
			return m.evalOp(v)
		case v["count"] != nil:
			return m.evalCount(asMap(v["count"]))
		case v["hash"] != nil:
			tbl, ok := m.work.tables[fmt.Sprint(asMap(v["hash"])["table"])]
			if !ok {
				return nil, fmt.Errorf("hash of unknown table")
			}
			b, err := m.tableBytes(tbl)
			if err != nil {
				return nil, err
			}
			return cdl.Fingerprint(b), nil
		case v["available"] != nil:
			return m.caps[fmt.Sprint(v["available"])], nil
		case v["unavailable"] != nil:
			return !m.caps[fmt.Sprint(v["unavailable"])], nil
		}
	}
	return nil, fmt.Errorf("unsupported expression form %T", e)
}

func (m *machine) evalRef(raw string) (any, error) {
	if id, pred, ok := splitRuleRef(raw); ok {
		for _, r := range m.doc.Declarations.Rules {
			if r.ID == id && r.Predicate == pred {
				return m.eval(r.Expr)
			}
		}
		return nil, fmt.Errorf("unknown rule predicate %s", raw)
	}
	if v, ok := m.work.states[raw]; ok {
		return v, nil
	}
	if v, ok := m.work.values[raw]; ok {
		return v, nil
	}
	for _, e := range m.doc.Declarations.Enums {
		for _, value := range e.Values {
			if value == raw {
				return value, nil
			}
		}
	}
	for _, s := range m.doc.Declarations.States {
		for _, value := range s.Values {
			if value == raw {
				return stateValue(value), nil
			}
		}
	}
	return nil, fmt.Errorf("unresolved reference %s", raw)
}

func (m *machine) evalOp(op map[string]any) (any, error) {
	operator := fmt.Sprint(op["op"])
	if operator == "NOT" {
		v, err := m.eval(op["e"])
		if err != nil {
			return nil, err
		}
		return !(v == true), nil
	}
	l, err := m.eval(op["l"])
	if err != nil {
		return nil, err
	}
	r, err := m.eval(op["r"])
	if err != nil {
		return nil, err
	}
	switch operator {
	case "==":
		return equal(l, r), nil
	case "!=":
		return !equal(l, r), nil
	case "AND":
		return l == true && r == true, nil
	case "OR":
		return l == true || r == true, nil
	case "+", "-":
		li, lok := toInt(l)
		ri, rok := toInt(r)
		if !lok || !rok {
			return nil, fmt.Errorf("operator %s requires integers", operator)
		}
		if operator == "+" {
			return li + ri, nil
		}
		return li - ri, nil
	case "<", "<=", ">", ">=":
		li, lok := toInt(l)
		ri, rok := toInt(r)
		if !lok || !rok {
			return nil, fmt.Errorf("operator %s requires integers", operator)
		}
		switch operator {
		case "<":
			return li < ri, nil
		case "<=":
			return li <= ri, nil
		case ">":
			return li > ri, nil
		default:
			return li >= ri, nil
		}
	}
	return nil, fmt.Errorf("unknown operator %s", operator)
}

func (m *machine) evalCount(count map[string]any) (any, error) {
	tbl, ok := m.work.tables[fmt.Sprint(count["table"])]
	if !ok {
		return nil, fmt.Errorf("count of unknown table")
	}
	where, hasWhere := count["where"]
	if !hasWhere {
		return len(tbl.rows), nil
	}
	w := asMap(where)
	field := fmt.Sprint(w["field"])
	want, err := m.eval(w["equals"])
	if err != nil {
		return nil, err
	}
	n := 0
	for _, row := range tbl.rows {
		if equal(row[field], want) {
			n++
		}
	}
	return n, nil
}

func (m *machine) tableBytes(tbl *tableState) ([]byte, error) {
	rows := append([]map[string]any(nil), tbl.rows...)
	sort.SliceStable(rows, func(i, j int) bool {
		return fmt.Sprint(rows[i][tbl.decl.Key]) < fmt.Sprint(rows[j][tbl.decl.Key])
	})
	return cdl.CanonicalBytes(rows)
}

func handlesAvailability(ops []map[string]any, capability string) bool {
	for _, op := range ops {
		if op["if"] != nil {
			ifMap := asMap(op["if"])
			cond := asMap(ifMap["cond"])
			if cond != nil && fmt.Sprint(cond["available"]) == capability {
				return true
			}
			if handlesAvailability(asMaps(ifMap["then"]), capability) || handlesAvailability(asMaps(ifMap["else"]), capability) {
				return true
			}
		}
	}
	return false
}

func stateValue(s string) any {
	switch s {
	case "true":
		return true
	case "false":
		return false
	default:
		return s
	}
}

func valueString(v any) string { return fmt.Sprint(v) }

func equal(a, b any) bool {
	ai, aok := toInt(a)
	bi, bok := toInt(b)
	if aok && bok {
		return ai == bi
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case float64:
		if n == math.Trunc(n) {
			return int(n), true
		}
	}
	return 0, false
}

func exprSubject(e any) string {
	if m, ok := e.(map[string]any); ok {
		if ref := m["ref"]; ref != nil {
			return fmt.Sprint(ref)
		}
	}
	return "condition"
}

func asMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func asMaps(v any) []map[string]any {
	items, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func asStrings(v any) []string {
	if typed, ok := v.([]string); ok {
		return typed
	}
	items, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

func splitRuleRef(raw string) (string, string, bool) {
	idx := strings.LastIndex(raw, ".")
	if idx <= 0 || idx == len(raw)-1 {
		return "", "", false
	}
	return raw[:idx], raw[idx+1:], true
}

func contains(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}
