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
