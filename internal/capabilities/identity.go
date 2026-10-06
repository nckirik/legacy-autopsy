package capabilities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var typedIDTypePattern = regexp.MustCompile(`^[A-Z][A-Z0-9-]*$`)

// TypedID computes the §4.1 canonical typed-record ID from the canonical key
// preimage parts. Collision extension beyond the base 12-character hash is a
// registry concern and is not performed here.
func TypedID(idType, namespace, owner, discriminator string) (string, error) {
	switch {
	case !typedIDTypePattern.MatchString(idType):
		return "", fmt.Errorf("typed-id: invalid protocol id type %q", idType)
	case namespace == "":
		return "", fmt.Errorf("typed-id: empty system namespace")
	case owner == "":
		return "", fmt.Errorf("typed-id: empty normalized owner coordinate")
	case discriminator == "":
		return "", fmt.Errorf("typed-id: empty semantic discriminator")
	}
	key := strings.Join([]string{idType, namespace, owner, discriminator}, "|")
	sum := sha256.Sum256([]byte(key))
	hash := strings.ToUpper(hex.EncodeToString(sum[:]))[:12]
	return idType + "-" + hash, nil
}

// SourceID computes the deterministic SRC identity for one source-inventory
// unit. The canonical key is `SRC|namespace|normalized-owner-coordinate|KIND`
// and the id is `SRC-KIND-<first 12 uppercase hex of SHA-256(key)>`, matching
// the §4.1 SRC-KIND-HASH form with kind included in the preimage.
func SourceID(namespace, kind, coordinate string) (string, error) {
	if namespace == "" {
		return "", fmt.Errorf("typed-id: empty system namespace")
	}
	if kind == "" {
		return "", fmt.Errorf("typed-id: empty SRC kind")
	}
	if coordinate == "" {
		return "", fmt.Errorf("typed-id: empty normalized owner coordinate")
	}
	key := strings.Join([]string{"SRC", namespace, coordinate, strings.ToUpper(kind)}, "|")
	sum := sha256.Sum256([]byte(key))
	hash := strings.ToUpper(hex.EncodeToString(sum[:]))[:12]
	return "SRC-" + strings.ToUpper(kind) + "-" + hash, nil
}
