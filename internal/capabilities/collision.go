package capabilities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// CollisionExtensionLengths are the declared extension lengths (§4.1).
var collisionExtensionLengths = []int{16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64}

var uppercaseHex = regexp.MustCompile(`^[A-F0-9]{12,64}$`)

// ExtendCollisionPrefixes extends two colliding uppercase-hex hashes to the first
// declared length where they differ. Existing shorter IDs receive aliases and
// controlled supersession; this function never mutates them.
func ExtendCollisionPrefixes(hashA, hashB string) (string, string, error) {
	if !uppercaseHex.MatchString(hashA) || !uppercaseHex.MatchString(hashB) {
		return "", "", fmt.Errorf("collision: hashes must be uppercase hex of at least 12 characters")
	}
	if hashA[:12] != hashB[:12] {
		return "", "", fmt.Errorf("collision: hashes do not collide at 12 characters")
	}
	if hashA == hashB {
		return "", "", fmt.Errorf("collision: identical hashes are not a collision")
	}
	for _, length := range collisionExtensionLengths {
		if len(hashA) < length || len(hashB) < length {
			break
		}
		if hashA[:length] != hashB[:length] {
			return hashA[:length], hashB[:length], nil
		}
	}
	return "", "", fmt.Errorf("collision: hashes do not differ up to 64 characters")
}

// ExtendCollision computes the full hashes of two distinct canonical keys and
// returns the first declared-length extension pair.
func ExtendCollision(keyA, keyB string) (string, string, error) {
	if keyA == keyB {
		return "", "", fmt.Errorf("collision: identical canonical keys")
	}
	sumA := sha256.Sum256([]byte(keyA))
	sumB := sha256.Sum256([]byte(keyB))
	hashA := strings.ToUpper(hex.EncodeToString(sumA[:]))
	hashB := strings.ToUpper(hex.EncodeToString(sumB[:]))
	if hashA[:12] != hashB[:12] {
		return "", "", fmt.Errorf("collision: keys do not collide at 12 characters")
	}
	return ExtendCollisionPrefixes(hashA, hashB)
}
