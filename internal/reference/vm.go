// Package reference is a deliberately small second implementation of the frozen
// pilot execution semantics. It consumes the same EIR and fixtures as the native
// runtime and exists to prove that execution semantics do not live accidentally
// inside Legacy Autopsy. It shares no execution code with internal/runtime.
package reference

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// Fixture mirrors the execution fixture contract independently.
type Fixture struct {
	Capabilities []string                             `json:"capabilities"`
	Registries   map[string][]map[string]any          `json:"registries"`
	Agent        map[string]map[string]map[string]any `json:"agent"`
}

// LoadFixture reads a synthetic execution fixture.
func LoadFixture(path string) (Fixture, error) {
	var fx Fixture
	b, err := os.ReadFile(path)
	if err != nil {
		return fx, err
	}
	if err := json.Unmarshal(b, &fx); err != nil {
		return fx, err
	}
	return fx, nil
}

// Effect mirrors the canonical effect record.
type Effect struct {
	Target    string `json:"target"`
	Operation string `json:"operation"`
	Value     string `json:"value,omitempty"`
}

// Diagnostic mirrors the canonical diagnostic record.
type Diagnostic struct {
	Code    string `json:"code"`
	Subject string `json:"subject"`
}

// StepRecord mirrors the canonical step record.
type StepRecord struct {
	Step        string       `json:"step"`
	Owner       string       `json:"owner"`
	Outcome     string       `json:"outcome"`
	Authority   string       `json:"authority"`
	Effects     []Effect     `json:"effects"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	EIRHash     string       `json:"eir-hash"`
	Channel     string       `json:"channel"`
}

// Trace mirrors the canonical execution trace.
type Trace struct {
	Steps []StepRecord `json:"steps"`
}

// Hash returns the canonical trace hash.
func (t Trace) Hash() (string, error) {
	b, err := cdl.CanonicalBytes(t)
	if err != nil {
		return "", err
	}
	return cdl.Fingerprint(b), nil
}

type vm struct {
	doc      *cdl.EIRDoc
	fx       Fixture
	caps     map[string]bool
	states   map[string]any
	values   map[string]any
	tables   map[string][]map[string]any
	gates    map[string]string
	bindings map[string]map[string]any
	trace    []StepRecord
}

// Execute runs all sections in EIR order.
func Execute(doc *cdl.EIRDoc, fx Fixture, caps map[string]bool) (Trace, error) {
	v := &vm{doc: doc, fx: fx, caps: caps, bindings: map[string]map[string]any{}}
	v.states = map[string]any{}
	v.values = map[string]any{}
	v.tables = map[string][]map[string]any{}
	v.gates = map[string]string{}
	for _, s := range doc.Declarations.States {
		if s.Initial != "" {
			v.states[s.ID] = literal(s.Initial)
		}
	}
	for _, t := range doc.Declarations.Tables {
		v.tables[t.ID] = []map[string]any{}
	}
	for _, g := range doc.Declarations.Gates {
		v.gates[g] = "open"
	}
	for _, section := range doc.Sections {
		for _, step := range section.Steps {
			if err := v.step(step); err != nil {
				return Trace{}, err
			}
		}
	}
	return Trace{Steps: v.trace}, nil
}

func literal(s string) any {
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	return s
}

func (v *vm) snapshot() *vm {
	out := &vm{doc: v.doc, fx: v.fx, caps: v.caps, bindings: map[string]map[string]any{}}
	out.states = map[string]any{}
	out.values = map[string]any{}
	out.tables = map[string][]map[string]any{}
	out.gates = map[string]string{}
	for k, val := range v.states {
		out.states[k] = val
	}
	for k, val := range v.values {
		out.values[k] = val
	}
	for k, val := range v.gates {
		out.gates[k] = val
	}
	for k, rows := range v.tables {
		copied := make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			cp := map[string]any{}
			for kk, vv := range row {
				cp[kk] = vv
			}
			copied = append(copied, cp)
		}
		out.tables[k] = copied
	}
	return out
}

func (v *vm) adopt(snap *vm) {
	v.states = snap.states
	v.values = snap.values
	v.tables = snap.tables
	v.gates = snap.gates
}

func (v *vm) step(step cdl.EIRStep) error {
	if step.Owner == "AGENT" {
		return v.agentStep(step)
	}
	return v.machineStep(step)
}

func (v *vm) machineStep(step cdl.EIRStep) error {
	rec := StepRecord{Step: step.ID, Owner: step.Owner, EIRHash: v.doc.Envelope.EIRHash, Channel: "native", Effects: []Effect{}, Diagnostics: []Diagnostic{}}
	missing := []string{}
	for _, cap := range step.Requires {
		if !v.caps[cap] && !condHandles(step.Ops, cap) {
			missing = append(missing, cap)
		}
	}
	if len(missing) > 0 {
		rec.Outcome = "unsupported"
		rec.Authority = "none"
		for _, cap := range missing {
			rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CAPABILITY_UNAVAILABLE", Subject: cap})
		}
		v.trace = append(v.trace, rec)
		return nil
	}
	staged := v.snapshot()
	effects := []Effect{}
	outcome, diag, err := v.ops(staged, step.Ops, &effects)
	if err != nil {
		rec.Outcome = "failed"
		rec.Authority = "none"
		rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "EIR_EVAL_ERROR", Subject: err.Error()})
		v.trace = append(v.trace, rec)
		return nil
	}
	if outcome == "blocked" {
		blocked := v.snapshot()
		var ignored []Effect
		if _, _, err := v.ops(blocked, otherwiseOps(step.Ops), &ignored); err != nil {
			return err
		}
		v.adopt(blocked)
		effects = ignored
		rec.Outcome = "blocked"
		rec.Authority = "authoritative"
		rec.Effects = effects
		rec.Diagnostics = append(rec.Diagnostics, *diag)
		v.trace = append(v.trace, rec)
		return nil
	}
	v.adopt(staged)
	rec.Outcome = "committed"
	rec.Authority = "authoritative"
	rec.Effects = effects
	v.trace = append(v.trace, rec)
	return nil
}

// ops executes a step body. It returns outcome "", "blocked", or an error.
func (v *vm) ops(target *vm, ops []map[string]any, effects *[]Effect) (string, *Diagnostic, error) {
	saved := v.adoptFor(target)
	defer v.restore(saved)
	for _, op := range ops {
		switch {
		case op["requires"] != nil:
		case op["set"] != nil:
			set := anyMap(op["set"])
			name := fmt.Sprint(set["target"])
			value, err := v.eval(target, set["value"])
			if err != nil {
				return "", nil, err
			}
			if current, ok := target.lookup(name); ok && same(current, value) {
				continue
			}
			if _, ok := target.states[name]; ok {
				target.states[name] = value
			} else {
				target.values[name] = value
			}
			*effects = append(*effects, Effect{Target: name, Operation: "set", Value: fmt.Sprint(value)})
		case op["append"] != nil:
			appendOp := anyMap(op["append"])
			tableID := fmt.Sprint(appendOp["to"])
			unit := v.bindings[fmt.Sprint(appendOp["for"])]
			if unit == nil {
				return "", nil, fmt.Errorf("append with unbound variable")
			}
			row := map[string]any{}
			for _, name := range v.rowFields(tableID) {
				if val, ok := unit[name]; ok {
					row[name] = val
				}
			}
			target.tables[tableID] = append(target.tables[tableID], row)
			*effects = append(*effects, Effect{Target: tableID, Operation: "append", Value: fmt.Sprint(row[v.tableKey(tableID)])})
		case op["foreach"] != nil:
			loop := anyMap(op["foreach"])
			name := fmt.Sprint(loop["var"])
			records, ok := v.fx.Registries[fmt.Sprint(loop["in"])]
			if !ok {
				return "", nil, fmt.Errorf("registry has no fixture data")
			}
			for _, rec := range records {
				v.bindings[name] = rec
				if _, _, err := v.ops(target, anyMaps(loop["do"]), effects); err != nil {
					return "", nil, err
				}
			}
			delete(v.bindings, name)
		case op["guard"] != nil:
			guard := anyMap(op["guard"])
			val, err := v.eval(target, guard["expr"])
			if err != nil {
				return "", nil, err
			}
			if val != true {
				diag := &Diagnostic{Code: "GUARD_BLOCKED", Subject: subject(guard["expr"])}
				return "blocked", diag, nil
			}
		case op["if"] != nil:
			branch := anyMap(op["if"])
			val, err := v.eval(target, branch["cond"])
			if err != nil {
				return "", nil, err
			}
			next := anyMaps(branch["else"])
			if val == true {
				next = anyMaps(branch["then"])
			}
			if _, _, err := v.ops(target, next, effects); err != nil {
				return "", nil, err
			}
		case op["block"] != nil:
			for _, gate := range anyStrings(anyMap(op["block"])["gates"]) {
				target.gates[gate] = "blocked"
				*effects = append(*effects, Effect{Target: gate, Operation: "block", Value: "blocked"})
			}
		default:
			return "", nil, fmt.Errorf("unsupported operation")
		}
	}
	return "", nil, nil
}

// adoptFor points the receiver at target's live maps for the duration of a body.
func (v *vm) adoptFor(target *vm) vm {
	saved := *v
	v.states = target.states
	v.values = target.values
	v.tables = target.tables
	v.gates = target.gates
	return saved
}

func (v *vm) restore(saved vm) {
	v.states = saved.states
	v.values = saved.values
	v.tables = saved.tables
	v.gates = saved.gates
	v.bindings = saved.bindings
}

func (v *vm) lookup(name string) (any, bool) {
	if val, ok := v.states[name]; ok {
		return val, true
	}
	val, ok := v.values[name]
	return val, ok
}

func (v *vm) rowFields(tableID string) []string {
	table := v.tableDecl(tableID)
	if table == nil {
		return nil
	}
	for _, f := range v.doc.Declarations.Fields {
		if f.ID == table.RowType {
			out := make([]string, 0, len(f.Fields))
			for _, spec := range f.Fields {
				out = append(out, spec.Name)
			}
			return out
		}
	}
	return nil
}

func (v *vm) tableKey(tableID string) string {
	t := v.tableDecl(tableID)
	if t == nil {
		return ""
	}
	return t.Key
}

func (v *vm) tableDecl(tableID string) *cdl.EIRTable {
	for i := range v.doc.Declarations.Tables {
		if v.doc.Declarations.Tables[i].ID == tableID {
			return &v.doc.Declarations.Tables[i]
		}
	}
	return nil
}

func (v *vm) agentStep(step cdl.EIRStep) error {
	rec := StepRecord{Step: step.ID, Owner: step.Owner, EIRHash: v.doc.Envelope.EIRHash, Channel: "native", Effects: []Effect{}, Diagnostics: []Diagnostic{}}
	staged := v.snapshot()
	data := v.fx.Agent[step.ID]
	var effects []Effect
	for _, produced := range step.Produces {
		parts := strings.SplitN(produced, ".", 2)
		if len(parts) != 2 {
			rec.Outcome = "failed"
			rec.Authority = "none"
			rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CONTRACT_INVALID", Subject: "malformed produces entry"})
			v.trace = append(v.trace, rec)
			return nil
		}
		tableID, field := parts[0], parts[1]
		keyField := v.tableKey(tableID)
		rows := staged.tables[tableID]
		enum, isEnum := v.enumFor(fieldType(v.doc, tableID, field))
		for _, row := range rows {
			key := fmt.Sprint(row[keyField])
			entry, ok := data[key]
			if !ok {
				rec.Outcome = "failed"
				rec.Authority = "none"
				rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " missing result for " + key})
				v.trace = append(v.trace, rec)
				return nil
			}
			value, ok := entry[field]
			if !ok {
				rec.Outcome = "failed"
				rec.Authority = "none"
				rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " missing field for " + key})
				v.trace = append(v.trace, rec)
				return nil
			}
			if isEnum && !has(enum, fmt.Sprint(value)) {
				rec.Outcome = "failed"
				rec.Authority = "none"
				rec.Diagnostics = append(rec.Diagnostics, Diagnostic{Code: "CONTRACT_INVALID", Subject: produced + " value outside RESULT domain for " + key})
				v.trace = append(v.trace, rec)
				return nil
			}
			row[field] = value
			effects = append(effects, Effect{Target: key + "." + field, Operation: "draft-set", Value: fmt.Sprint(value)})
		}
	}
	v.adopt(staged)
	rec.Outcome = "committed"
	rec.Authority = "draft"
	rec.Effects = effects
	v.trace = append(v.trace, rec)
	return nil
}

func (v *vm) enumFor(id string) ([]string, bool) {
	for _, e := range v.doc.Declarations.Enums {
		if e.ID == id {
			return e.Values, true
		}
	}
	return nil, false
}

func fieldType(doc *cdl.EIRDoc, tableID, field string) string {
	var rowType string
	for _, t := range doc.Declarations.Tables {
		if t.ID == tableID {
			rowType = t.RowType
		}
	}
	for _, f := range doc.Declarations.Fields {
		if f.ID == rowType {
			for _, spec := range f.Fields {
				if spec.Name == field {
					return spec.Type
				}
			}
		}
	}
	return ""
}

func (v *vm) eval(target *vm, e any) (any, error) {
	saved := v.adoptFor(target)
	defer v.restore(saved)
	switch val := e.(type) {
	case bool, string:
		return val, nil
	case float64:
		if val == math.Trunc(val) {
			return int(val), nil
		}
		return val, nil
	case map[string]any:
		if val["ref"] != nil {
			return v.ref(fmt.Sprint(val["ref"]))
		}
		if val["op"] != nil {
			return v.op(val)
		}
		if val["count"] != nil {
			return v.count(anyMap(val["count"]))
		}
		if val["hash"] != nil {
			tableID := fmt.Sprint(anyMap(val["hash"])["table"])
			rows := append([]map[string]any(nil), v.tables[tableID]...)
			key := v.tableKey(tableID)
			sort.SliceStable(rows, func(i, j int) bool {
				return fmt.Sprint(rows[i][key]) < fmt.Sprint(rows[j][key])
			})
			blob, err := cdl.CanonicalBytes(rows)
			if err != nil {
				return nil, err
			}
			return cdl.Fingerprint(blob), nil
		}
		if val["available"] != nil {
			return v.caps[fmt.Sprint(val["available"])], nil
		}
		if val["unavailable"] != nil {
			return !v.caps[fmt.Sprint(val["unavailable"])], nil
		}
	}
	return nil, fmt.Errorf("unsupported expression")
}

func (v *vm) ref(name string) (any, error) {
	if id, pred, ok := cut(name); ok {
		for _, rule := range v.doc.Declarations.Rules {
			if rule.ID == id && rule.Predicate == pred {
				return v.eval(v, rule.Expr)
			}
		}
		return nil, fmt.Errorf("unknown rule predicate %s", name)
	}
	if val, ok := v.states[name]; ok {
		return val, nil
	}
	if val, ok := v.values[name]; ok {
		return val, nil
	}
	for _, e := range v.doc.Declarations.Enums {
		for _, value := range e.Values {
			if value == name {
				return value, nil
			}
		}
	}
	for _, s := range v.doc.Declarations.States {
		for _, value := range s.Values {
			if value == name {
				return literal(value), nil
			}
		}
	}
	return nil, fmt.Errorf("unresolved reference %s", name)
}

func (v *vm) op(expr map[string]any) (any, error) {
	operator := fmt.Sprint(expr["op"])
	if operator == "NOT" {
		val, err := v.eval(v, expr["e"])
		if err != nil {
			return nil, err
		}
		return !(val == true), nil
	}
	l, err := v.eval(v, expr["l"])
	if err != nil {
		return nil, err
	}
	r, err := v.eval(v, expr["r"])
	if err != nil {
		return nil, err
	}
	switch operator {
	case "==":
		return same(l, r), nil
	case "!=":
		return !same(l, r), nil
	case "AND":
		return l == true && r == true, nil
	case "OR":
		return l == true || r == true, nil
	case "+":
		return num(l) + num(r), nil
	case "-":
		return num(l) - num(r), nil
	case "<":
		return num(l) < num(r), nil
	case "<=":
		return num(l) <= num(r), nil
	case ">":
		return num(l) > num(r), nil
	case ">=":
		return num(l) >= num(r), nil
	}
	return nil, fmt.Errorf("unknown operator")
}

func (v *vm) count(count map[string]any) (any, error) {
	rows := v.tables[fmt.Sprint(count["table"])]
	where, ok := count["where"]
	if !ok {
		return len(rows), nil
	}
	selector := anyMap(where)
	want, err := v.eval(v, selector["equals"])
	if err != nil {
		return nil, err
	}
	n := 0
	for _, row := range rows {
		if same(row[fmt.Sprint(selector["field"])], want) {
			n++
		}
	}
	return n, nil
}

func condHandles(ops []map[string]any, capability string) bool {
	for _, op := range ops {
		if op["if"] == nil {
			continue
		}
		branch := anyMap(op["if"])
		cond := anyMap(branch["cond"])
		if fmt.Sprint(cond["available"]) == capability {
			return true
		}
		if condHandles(anyMaps(branch["then"]), capability) || condHandles(anyMaps(branch["else"]), capability) {
			return true
		}
	}
	return false
}

func otherwiseOps(ops []map[string]any) []map[string]any {
	for _, op := range ops {
		if op["guard"] != nil {
			return anyMaps(anyMap(op["guard"])["otherwise"])
		}
	}
	return nil
}

func anyMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func anyMaps(v any) []map[string]any {
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

func anyStrings(v any) []string {
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

func subject(e any) string {
	if m, ok := e.(map[string]any); ok && m["ref"] != nil {
		return fmt.Sprint(m["ref"])
	}
	return "condition"
}

func cut(raw string) (string, string, bool) {
	idx := strings.LastIndex(raw, ".")
	if idx <= 0 || idx == len(raw)-1 {
		return "", "", false
	}
	return raw[:idx], raw[idx+1:], true
}

func same(a, b any) bool {
	ai, aok := intOf(a)
	bi, bok := intOf(b)
	if aok && bok {
		return ai == bi
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func num(v any) int {
	n, _ := intOf(v)
	return n
}

func intOf(v any) (int, bool) {
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

func has(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}
