// Package capabilities provides the deterministic capability set for the pilot.
// Proposing providers and effect services are deliberately absent.
package capabilities

import "crypto/sha256"

// Set declares which capabilities are available in one execution.
type Set map[string]bool

// Deterministic returns the pilot's deterministic capabilities.
func Deterministic() Set {
	return Set{"hash.sha256": true, "table.serialize": true}
}

// Without returns a copy with the named capabilities removed.
func (s Set) Without(ids ...string) Set {
	out := Set{}
	for k, v := range s {
		out[k] = v
	}
	for _, id := range ids {
		delete(out, id)
	}
	return out
}

// Available reports capability availability.
func (s Set) Available(id string) bool { return s[id] }

// Hash returns the sha256 capability used by the pilot VM.
func Hash(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hexEncode(sum[:])
}

func hexEncode(b []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, 0, len(b)*2)
	for _, x := range b {
		out = append(out, digits[x>>4], digits[x&0x0f])
	}
	return string(out)
}
