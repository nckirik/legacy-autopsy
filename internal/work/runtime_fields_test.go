package work

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/store"
)

// TestRuntimeInvocationFieldsInjectInvocationID proves that a declared
// invocation field receives the current invocation identity, never the record
// ID: lifecycle eligibility rows must record the classifying invocation.
func TestRuntimeInvocationFieldsInjectInvocationID(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{
		ID:                      "U-LIFECYCLE-TEST",
		Mode:                    "Discovery",
		IDFromFields:            []string{"scope-id"},
		RuntimeInvocationFields: []string{"classified-by-invocation"},
		ExpectedRecords:         1,
		ResultType:              "LIFECYCLE-ELIGIBILITY",
	}
	payload := "- **scope-id:** SRC-FILE-AAAAAAAAAAAA\n" +
		"- **classification:** Unknown\n" +
		"- **reason:** bounded static evidence\n" +
		"- **clm-evidence-bindings:** CLM-AAAAAAAAAAAA\n" +
		"- **classified-by-invocation:** PENDING\n" +
		"- **snapshot:** dev\n"
	executor := fakeExecutorWork(func(inv Invocation) (Proposal, error) {
		return Proposal{FieldID: "LIFECYCLE-ELIGIBILITY", Records: []Record{{Payload: []byte(payload)}}}, nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err != nil {
		t.Fatal(err)
	}
	stored, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", "SRC-FILE-AAAAAAAAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stored), "- **classified-by-invocation:** "+report.Units[0].InvocationID) {
		t.Fatalf("stored payload lacks the invocation ID:\n%s", stored)
	}
}

type fakeExecutorWork func(Invocation) (Proposal, error)

func (f fakeExecutorWork) Execute(inv Invocation) (Proposal, error) { return f(inv) }
