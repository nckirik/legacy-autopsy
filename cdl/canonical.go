package cdl

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CanonicalBytes serializes a value per the EIR canonical form: sorted object
// keys, arrays in semantic order, no insignificant whitespace, no HTML escaping,
// and no trailing newline.
func CanonicalBytes(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var normalized any
	if err := dec.Decode(&normalized); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(normalized); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// Fingerprint returns the sha256:<hex> identity of a byte sequence.
func Fingerprint(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}
