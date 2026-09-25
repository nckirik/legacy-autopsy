// Package cdl implements the Canonical Deconstruction Language toolchain for the
// frozen pilot subset: parsing, resolution, checks, EIR emission, and prompt
// rendering. It must not import Legacy Autopsy runtime packages.
package cdl

// Versions pins the language, stdlib, and generator identities recorded in EIR.
type Versions struct {
	Language  string
	Stdlib    string
	Generator string
}

// Program is one parsed source file.
type Program struct {
	Globals     Globals
	Sections    []Section
	Projections []Projection
}

// Globals holds document-level declarations.
type Globals struct {
	Capabilities    []Capability
	Registries      []string
	Artifacts       []string
	Rules           []string
	Gates           []string
	WorkflowTargets []string
}

// Capability is a declared supplied capability.
type Capability struct{ ID, Kind string }

// Use is a section-level semantic dependency.
type Use struct{ Kind, ID string }

// Section is one normative section.
type Section struct {
	ID       string
	Number   string
	Title    string
	Artifact string
	Goal     string
	Uses     []Use
	Requires []string
	Types    []TypeDecl
	States   []StateDecl
	Values   []ValueDecl
	Enums    []EnumDecl
	Fields   []FieldDecl
	Tables   []TableDecl
	Rules    []RuleDecl
	Steps    []Step
}

// TypeDecl names a value type.
type TypeDecl struct{ ID, Type string }

// StateDecl declares values, initialization, and allowed transitions.
type StateDecl struct {
	ID      string
	Values  []string
	Initial string
	Allows  [][2]string
}

// ValueDecl is a durable computed value.
type ValueDecl struct{ ID, Type, Lifetime string }

// EnumDecl is a closed value domain.
type EnumDecl struct {
	ID     string
	Values []string
}

// FieldSpec is one record field.
type FieldSpec struct {
	Name     string
	Type     string
	Required bool
}

// FieldDecl is a record shape.
type FieldDecl struct {
	ID     string
	Fields []FieldSpec
}

// TableDecl is a table shape over a registry.
type TableDecl struct{ ID, Rows, RowType, Key string }

// RuleDecl is a named predicate or referenced rule.
type RuleDecl struct {
	ID        string
	Goal      string
	Predicate string
	Expr      Expr
}

// Step is one executable step with exactly one owner.
type Step struct {
	ID        string
	Owner     string
	UsesRules []string
	Requires  []string
	Ops       []Op
	Evidence  []string
	Produces  []string
	Result    string
	Require   string
}

// Projection is a generated consumer view.
type Projection struct {
	ID        string
	Channels  []string
	Covers    string
	StepTexts map[string]string
	RuleTexts map[string]string
	Omissions []Omission
}

// Omission drops a step from a channel with a declared supply.
type Omission struct {
	Step       string
	Channels   []string
	Reason     string
	SuppliedBy string
}

// Op is a MACHINE operation.
type Op interface{ op() }

// RequiresOp declares an operation-level capability requirement.
type RequiresOp struct{ Caps []string }

// SetOp assigns a state or value.
type SetOp struct {
	Target string
	Value  Expr
}

// LetOp binds a step-local name.
type LetOp struct {
	Target string
	Value  Expr
}

// AppendOp appends a row for a bound variable; To is resolved by the compiler.
type AppendOp struct{ RowType, For, To string }

// ForEachOp iterates a registry.
type ForEachOp struct {
	Var  string
	In   string
	Body []Op
}

// GuardOp blocks when the condition is not true.
type GuardOp struct {
	Expr      Expr
	Otherwise []Op
}

// IfOp is a two-branch conditional.
type IfOp struct {
	Cond Expr
	Then []Op
	Else []Op
}

// BlockOp blocks named gates.
type BlockOp struct{ Gates []string }

func (RequiresOp) op() {}
func (SetOp) op()      {}
func (LetOp) op()      {}
func (AppendOp) op()   {}
func (ForEachOp) op()  {}
func (GuardOp) op()    {}
func (IfOp) op()       {}
func (BlockOp) op()    {}

// Expr is an expression.
type Expr interface{ expr() }

// Ident references a declaration, rule predicate, or enum value.
type Ident struct{ Raw string }

// Literal is an integer, string, or boolean literal.
type Literal struct{ Value any }

// Binary is an infix operation.
type Binary struct {
	Op   string
	L, R Expr
}

// Unary is a prefix operation.
type Unary struct {
	Op string
	E  Expr
}

// Where is a COUNT selector.
type Where struct {
	Field string
	Value Expr
}

// Call is a function call; COUNT carries a Where.
type Call struct {
	Name  string
	Args  []Expr
	Where *Where
}

func (Ident) expr()   {}
func (Literal) expr() {}
func (Binary) expr()  {}
func (Unary) expr()   {}
func (Call) expr()    {}
