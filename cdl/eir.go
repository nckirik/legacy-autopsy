package cdl

import "encoding/json"

// Frozen pilot identities. They are recorded in every EIR envelope and every
// generated projection provenance header.
const (
	LanguageVersion  = "cdl/0.1"
	StdlibVersion    = "cdl-stdlib/0.1"
	GeneratorVersion = "cdl/0.1.0"
	EIRFormat        = 1
	ProtocolVersion  = "canonical-deconstruction/4.1.2"
)

// EIRDoc is the typed execution IR contract.
type EIRDoc struct {
	Envelope     EIREnvelope     `json:"envelope"`
	Declarations EIRDeclarations `json:"declarations"`
	Sections     []EIRSection    `json:"sections"`
	Projections  []EIRProjection `json:"projections"`
}

// EIREnvelope binds the document to its versions and hashes.
type EIREnvelope struct {
	EIRFormat         int    `json:"eir-format"`
	Language          string `json:"language"`
	Stdlib            string `json:"stdlib"`
	Protocol          string `json:"protocol"`
	SourceFingerprint string `json:"source-fingerprint"`
	Generator         string `json:"generator"`
	EIRHash           string `json:"eir-sha256,omitempty"`
}

// EIRDeclarations aggregates document and section declarations.
type EIRDeclarations struct {
	Capabilities    []EIRCapability `json:"capabilities"`
	Registries      []string        `json:"registries"`
	Artifacts       []string        `json:"artifacts"`
	Gates           []string        `json:"gates"`
	WorkflowTargets []string        `json:"workflow-targets,omitempty"`
	Types           []EIRType       `json:"types,omitempty"`
	States          []EIRState      `json:"states,omitempty"`
	Values          []EIRValue      `json:"values,omitempty"`
	Enums           []EIREnum       `json:"enums,omitempty"`
	Fields          []EIRField      `json:"fields,omitempty"`
	Tables          []EIRTable      `json:"tables,omitempty"`
	Rules           []EIRRule       `json:"rules,omitempty"`
}

// EIRCapability is a supplied capability with its kind.
type EIRCapability struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Version int    `json:"version"`
}

// EIRType is a named value type.
type EIRType struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Section string `json:"section"`
}

// EIRState is a declared state.
type EIRState struct {
	ID      string      `json:"id"`
	Values  []string    `json:"values"`
	Initial string      `json:"initial"`
	Allows  [][2]string `json:"allows"`
	Section string      `json:"section"`
}

// EIRValue is a declared durable value.
type EIRValue struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Lifetime string `json:"lifetime"`
	Section  string `json:"section"`
}

// EIREnum is a closed value domain.
type EIREnum struct {
	ID      string   `json:"id"`
	Values  []string `json:"values"`
	Section string   `json:"section"`
}

// EIRField is a record shape.
type EIRField struct {
	ID      string         `json:"id"`
	Fields  []EIRFieldSpec `json:"fields"`
	Section string         `json:"section"`
}

// EIRFieldSpec is one record field.
type EIRFieldSpec struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// EIRTable is a table shape over a registry.
type EIRTable struct {
	ID      string `json:"id"`
	Rows    string `json:"rows"`
	RowType string `json:"row-type"`
	Key     string `json:"key"`
	Section string `json:"section"`
}

// EIRRule is a named predicate.
type EIRRule struct {
	ID        string `json:"id"`
	Goal      string `json:"goal,omitempty"`
	Predicate string `json:"predicate,omitempty"`
	Expr      any    `json:"expr,omitempty"`
	Section   string `json:"section,omitempty"`
}

// EIRSection is one compiled section.
type EIRSection struct {
	ID       string    `json:"id"`
	Number   string    `json:"number"`
	Title    string    `json:"title"`
	Artifact string    `json:"artifact"`
	Goal     string    `json:"goal"`
	Uses     []EIRUse  `json:"uses,omitempty"`
	Requires []string  `json:"requires"`
	Steps    []EIRStep `json:"steps"`
}

// EIRUse is a semantic dependency.
type EIRUse struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// EIRStep is one executable step.
type EIRStep struct {
	ID        string           `json:"id"`
	Owner     string           `json:"owner"`
	UsesRules []string         `json:"uses-rules,omitempty"`
	Requires  []string         `json:"requires,omitempty"`
	Ops       []map[string]any `json:"ops,omitempty"`
	Evidence  []string         `json:"evidence,omitempty"`
	Produces  []string         `json:"produces,omitempty"`
	Result    string           `json:"result,omitempty"`
	Require   string           `json:"require,omitempty"`
}

// EIRProjection is a generated consumer view.
type EIRProjection struct {
	ID        string            `json:"id"`
	Channels  []string          `json:"channels"`
	Covers    string            `json:"covers,omitempty"`
	StepTexts map[string]string `json:"step-texts,omitempty"`
	RuleTexts map[string]string `json:"rule-texts,omitempty"`
	Omissions []EIROmission     `json:"omissions,omitempty"`
}

// EIROmission drops a step from a channel with a declared supply.
type EIROmission struct {
	Step       string   `json:"step"`
	Channels   []string `json:"channels,omitempty"`
	Reason     string   `json:"reason"`
	SuppliedBy string   `json:"supplied-by"`
}

// BuildEIR converts a resolved program into EIR. The returned bytes are the
// canonical document including the computed eir-sha256; the hash covers the
// document with an empty hash field.
func BuildEIR(prog *Program, versions Versions, sourceFingerprint string) (*EIRDoc, []byte, string, error) {
	if versions.Language == "" {
		versions = Versions{Language: LanguageVersion, Stdlib: StdlibVersion, Generator: GeneratorVersion}
	}
	doc := &EIRDoc{
		Envelope: EIREnvelope{
			EIRFormat:         EIRFormat,
			Language:          versions.Language,
			Stdlib:            versions.Stdlib,
			Protocol:          ProtocolVersion,
			SourceFingerprint: sourceFingerprint,
			Generator:         versions.Generator,
		},
	}
	for _, c := range prog.Globals.Capabilities {
		doc.Declarations.Capabilities = append(doc.Declarations.Capabilities, EIRCapability{ID: c.ID, Kind: c.Kind, Version: 1})
	}
	doc.Declarations.Registries = append(doc.Declarations.Registries, prog.Globals.Registries...)
	doc.Declarations.Artifacts = append(doc.Declarations.Artifacts, prog.Globals.Artifacts...)
	doc.Declarations.Gates = append(doc.Declarations.Gates, prog.Globals.Gates...)
	doc.Declarations.WorkflowTargets = append(doc.Declarations.WorkflowTargets, prog.Globals.WorkflowTargets...)
	for _, id := range prog.Globals.Rules {
		doc.Declarations.Rules = append(doc.Declarations.Rules, EIRRule{ID: id})
	}
	for _, sec := range prog.Sections {
		doc.Declarations.Types = append(doc.Declarations.Types, sectionTypes(sec)...)
		doc.Declarations.States = append(doc.Declarations.States, sectionStates(sec)...)
		doc.Declarations.Values = append(doc.Declarations.Values, sectionValues(sec)...)
		doc.Declarations.Enums = append(doc.Declarations.Enums, sectionEnums(sec)...)
		doc.Declarations.Fields = append(doc.Declarations.Fields, sectionFields(sec)...)
		doc.Declarations.Tables = append(doc.Declarations.Tables, sectionTables(sec)...)
		for _, r := range sec.Rules {
			doc.Declarations.Rules = append(doc.Declarations.Rules, EIRRule{
				ID: r.ID, Goal: r.Goal, Predicate: r.Predicate, Expr: exprToEIR(r.Expr), Section: sec.ID,
			})
		}
		doc.Sections = append(doc.Sections, sectionEIR(sec))
	}
	for _, p := range prog.Projections {
		doc.Projections = append(doc.Projections, EIRProjection{
			ID:        p.ID,
			Channels:  p.Channels,
			Covers:    p.Covers,
			StepTexts: p.StepTexts,
			RuleTexts: p.RuleTexts,
			Omissions: omissionsEIR(p.Omissions),
		})
	}
	hashInput := *doc
	hashInput.Envelope.EIRHash = ""
	hashBytes, err := CanonicalBytes(hashInput)
	if err != nil {
		return nil, nil, "", err
	}
	hash := Fingerprint(hashBytes)
	doc.Envelope.EIRHash = hash
	finalBytes, err := CanonicalBytes(doc)
	if err != nil {
		return nil, nil, "", err
	}
	return doc, finalBytes, hash, nil
}

// DecodeEIR decodes canonical EIR bytes.
func DecodeEIR(b []byte) (*EIRDoc, error) {
	var doc EIRDoc
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func sectionTypes(sec Section) []EIRType {
	var out []EIRType
	for _, t := range sec.Types {
		out = append(out, EIRType{ID: t.ID, Type: t.Type, Section: sec.ID})
	}
	return out
}

func sectionStates(sec Section) []EIRState {
	var out []EIRState
	for _, s := range sec.States {
		out = append(out, EIRState{ID: s.ID, Values: s.Values, Initial: s.Initial, Allows: s.Allows, Section: sec.ID})
	}
	return out
}

func sectionValues(sec Section) []EIRValue {
	var out []EIRValue
	for _, v := range sec.Values {
		out = append(out, EIRValue{ID: v.ID, Type: v.Type, Lifetime: v.Lifetime, Section: sec.ID})
	}
	return out
}

func sectionEnums(sec Section) []EIREnum {
	var out []EIREnum
	for _, e := range sec.Enums {
		out = append(out, EIREnum{ID: e.ID, Values: e.Values, Section: sec.ID})
	}
	return out
}

func sectionFields(sec Section) []EIRField {
	var out []EIRField
	for _, f := range sec.Fields {
		field := EIRField{ID: f.ID, Section: sec.ID}
		for _, spec := range f.Fields {
			field.Fields = append(field.Fields, EIRFieldSpec{Name: spec.Name, Type: spec.Type, Required: spec.Required})
		}
		out = append(out, field)
	}
	return out
}

func sectionTables(sec Section) []EIRTable {
	var out []EIRTable
	for _, t := range sec.Tables {
		out = append(out, EIRTable{ID: t.ID, Rows: t.Rows, RowType: t.RowType, Key: t.Key, Section: sec.ID})
	}
	return out
}

func sectionEIR(sec Section) EIRSection {
	out := EIRSection{
		ID: sec.ID, Number: sec.Number, Title: sec.Title, Artifact: sec.Artifact,
		Goal: sec.Goal, Requires: sec.Requires,
	}
	for _, u := range sec.Uses {
		out.Uses = append(out.Uses, EIRUse{Kind: u.Kind, ID: u.ID})
	}
	for _, s := range sec.Steps {
		step := EIRStep{
			ID: s.ID, Owner: s.Owner, UsesRules: s.UsesRules, Requires: s.Requires,
			Evidence: s.Evidence, Produces: s.Produces, Result: s.Result, Require: s.Require,
		}
		for _, op := range s.Ops {
			step.Ops = append(step.Ops, opToEIR(op))
		}
		out.Steps = append(out.Steps, step)
	}
	return out
}

func omissionsEIR(in []Omission) []EIROmission {
	var out []EIROmission
	for _, o := range in {
		out = append(out, EIROmission{Step: o.Step, Channels: o.Channels, Reason: o.Reason, SuppliedBy: o.SuppliedBy})
	}
	return out
}

func opToEIR(op Op) map[string]any {
	switch v := op.(type) {
	case RequiresOp:
		return map[string]any{"requires": map[string]any{"capabilities": v.Caps}}
	case SetOp:
		return map[string]any{"set": map[string]any{"target": v.Target, "value": exprToEIR(v.Value)}}
	case LetOp:
		return map[string]any{"let": map[string]any{"target": v.Target, "value": exprToEIR(v.Value)}}
	case AppendOp:
		return map[string]any{"append": map[string]any{"row-type": v.RowType, "for": v.For, "to": v.To}}
	case ForEachOp:
		body := make([]any, 0, len(v.Body))
		for _, b := range v.Body {
			body = append(body, opToEIR(b))
		}
		return map[string]any{"foreach": map[string]any{"var": v.Var, "in": v.In, "do": body}}
	case GuardOp:
		otherwise := make([]any, 0, len(v.Otherwise))
		for _, b := range v.Otherwise {
			otherwise = append(otherwise, opToEIR(b))
		}
		return map[string]any{"guard": map[string]any{"expr": exprToEIR(v.Expr), "otherwise": otherwise}}
	case IfOp:
		then := make([]any, 0, len(v.Then))
		for _, b := range v.Then {
			then = append(then, opToEIR(b))
		}
		els := make([]any, 0, len(v.Else))
		for _, b := range v.Else {
			els = append(els, opToEIR(b))
		}
		return map[string]any{"if": map[string]any{"cond": exprToEIR(v.Cond), "then": then, "else": els}}
	case BlockOp:
		return map[string]any{"block": map[string]any{"gates": v.Gates}}
	default:
		return map[string]any{"unsupported": true}
	}
}

func exprToEIR(e Expr) any {
	switch v := e.(type) {
	case nil:
		return nil
	case Ident:
		return map[string]any{"ref": v.Raw}
	case Literal:
		return v.Value
	case Binary:
		return map[string]any{"op": v.Op, "l": exprToEIR(v.L), "r": exprToEIR(v.R)}
	case Unary:
		return map[string]any{"op": v.Op, "e": exprToEIR(v.E)}
	case Call:
		switch v.Name {
		case "COUNT":
			if len(v.Args) == 1 {
				if table, ok := v.Args[0].(Ident); ok {
					sel := map[string]any{"table": table.Raw}
					if v.Where != nil {
						sel["where"] = map[string]any{"field": v.Where.Field, "equals": exprToEIR(v.Where.Value)}
					}
					return map[string]any{"count": sel}
				}
			}
		case "HASH":
			if len(v.Args) == 1 {
				if table, ok := v.Args[0].(Ident); ok {
					return map[string]any{"hash": map[string]any{"table": table.Raw}}
				}
			}
		case "AVAILABLE":
			if len(v.Args) == 1 {
				if cap, ok := v.Args[0].(Ident); ok {
					return map[string]any{"available": cap.Raw}
				}
			}
		case "UNAVAILABLE":
			if len(v.Args) == 1 {
				if cap, ok := v.Args[0].(Ident); ok {
					return map[string]any{"unavailable": cap.Raw}
				}
			}
		}
		args := make([]any, 0, len(v.Args))
		for _, a := range v.Args {
			args = append(args, exprToEIR(a))
		}
		return map[string]any{"call": v.Name, "args": args}
	default:
		return nil
	}
}
