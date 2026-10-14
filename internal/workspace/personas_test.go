package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func personaTree(t *testing.T, directories ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, directory := range append([]string{"_shared"}, directories...) {
		path := filepath.Join(root, "personas", directory)
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestValidatePersonaOwnership(t *testing.T) {
	root := personaTree(t, "WEB-web-kullanici", "AD-admin")
	err := ValidatePersonaOwnership(root, []PersonaRegistration{
		{Prefix: "WEB", Slug: "web-kullanici", Directory: "WEB-web-kullanici"},
		{Prefix: "AD", Slug: "admin", Directory: "AD-admin"},
	})
	if err != nil {
		t.Fatalf("valid layout rejected: %v", err)
	}
}

func TestValidatePersonaOwnershipFailures(t *testing.T) {
	root := personaTree(t, "WEB-web-kullanici", "XX-unknown")
	err := ValidatePersonaOwnership(root, []PersonaRegistration{
		{Prefix: "WEB", Slug: "web-kullanici", Directory: "WEB-web-kullanici"},
		{Prefix: "AD", Slug: "admin", Directory: "AD-admin"},
	})
	if err == nil {
		t.Fatal("missing registered directory and unregistered directory accepted")
	}
	joined := err.Error()
	for _, want := range []string{"no workspace directory", "not registered in 0A"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("diagnostics missing %q:\n%s", want, joined)
		}
	}

	root = personaTree(t, "WEB-web-kullanici")
	err = ValidatePersonaOwnership(root, []PersonaRegistration{
		{Prefix: "WEB", Slug: "web-kullanici", Directory: "WEB-web-kullanici"},
		{Prefix: "WEB", Slug: "web-kullanici", Directory: "WEB-web-kullanici"},
	})
	joined = err.Error()
	if !strings.Contains(joined, "registered twice") {
		t.Fatalf("duplicate prefix diagnostics missing:\n%s", joined)
	}

	root = personaTree(t, "AD-admin")
	_ = root
	if err := ValidatePersonaOwnership(root, []PersonaRegistration{{Prefix: "_shared", Slug: "x", Directory: "_shared-x"}}); err == nil || !strings.Contains(err.Error(), "must not be registered") {
		t.Fatalf("_shared registration accepted: %v", err)
	}
}
