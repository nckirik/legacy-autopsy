package store

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func TestAuditEnumRegistriesComplete(t *testing.T) {
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	diagnostics, err := AuditEnumRegistries(res.EIR)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("enum registry diagnostics:\n%s", strings.Join(diagnostics, "\n"))
	}
}

func TestAuditEnumRegistriesFindsDuplicateAndEmptyValues(t *testing.T) {
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	doc := res.EIR
	if len(doc.Declarations.Enums) == 0 {
		t.Fatal("compiled document declares no enums")
	}
	original := append([]string(nil), doc.Declarations.Enums[0].Values...)
	defer func() { doc.Declarations.Enums[0].Values = original }()

	doc.Declarations.Enums[0].Values = append(append([]string(nil), original...), original[0])
	diagnostics, err := AuditEnumRegistries(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !containsDiagnostic(diagnostics, "repeats value") {
		t.Fatalf("duplicate value not reported: %v", diagnostics)
	}

	doc.Declarations.Enums[0].Values = nil
	diagnostics, err = AuditEnumRegistries(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !containsDiagnostic(diagnostics, "declares no values") {
		t.Fatalf("empty enum not reported: %v", diagnostics)
	}
}

func TestEnumCheckValuesParsesGeneratedClauses(t *testing.T) {
	values, ok := enumCheckValues(`  "status" TEXT CHECK ("status" IN ('Active', 'O''Brien', 'Retired')),`, "status")
	if !ok {
		t.Fatal("CHECK clause not found")
	}
	if strings.Join(values, "|") != "Active|O'Brien|Retired" {
		t.Fatalf("values = %v", values)
	}
	if _, ok := enumCheckValues(`  "other" TEXT CHECK ("other" IN ('x')),`, "status"); ok {
		t.Fatal("matched the wrong column")
	}
}

func containsDiagnostic(diagnostics []string, substring string) bool {
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic, substring) {
			return true
		}
	}
	return false
}
